# dgo (Dynamic Go)

[![](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![](https://goreportcard.com/badge/github.com/lyraproj/dgo)](https://goreportcard.com/report/github.com/lyraproj/dgo)
[![](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/lyraproj/dgo)
[![](https://github.com/lyraproj/dgo/workflows/Dgo%20Test/badge.svg)](https://github.com/lyraproj/dgo/actions)
[![](https://coveralls.io/repos/github/lyraproj/dgo/badge.svg)](https://coveralls.io/github/lyraproj/dgo)

Dgo's main objectives are: Type Constraints, Immutability, Collections, Serialization, Encapsulation, Extendability,
and Performance. It is designed to make working with dynamic values an effortless and type safe task.

## Install
To use dgo, first install the latest version of the library:
```sh
go get github.com/lyraproj/dgo
```
Next, include needed packages in your application. Pick what you need from the following packages:
```go
import (
  "github.com/lyraproj/dgo/dgo"
  "github.com/lyraproj/dgo/typ"
  "github.com/lyraproj/dgo/tf"
  "github.com/lyraproj/dgo/vf"
)
```

## Type Constraints

### Dgo types versus Go's native types
Go is a typed language, but the types are not very descriptive. It is for instance not possible to declare a type
that corresponds only to a specific range of integers, or a string that must confirm to a specific pattern. All such
constraints must be expressed as code wherever a value of the type is assigned. In dgo a type describing a range
of integers can be declared as `0..15` and a pattern constrained string can be declared as `/^[a-z]+$/`.

It is also not possible to declare type combinations such as a slice that can contain only integers or floats. If
a value can be of more than one of go's native types, then it must be declared as an `interface{}` which corresponds
to every possible value in the system. This means that a slice containing ints and floats must be declared as
`[]interface{}` wich is a declaration of a slice that may contain any type of value. In dgo, such a type can be
declared as `[](int|float)`. Other examples:
 
- `map[string](string|int)` (string keyed map of string
- `[]0..15` (slice of integers ranging from 0 - 15).
- `"red"|"green"|"blue"` (enumeration of strings)
- `2|8|10|16` (enumeration of integers)

Dgo is influenced by restrictive type constraint languages such as:
- [CUE](https://cue.googlesource.com/cue/+/HEAD/doc/ref/spec.md)
- [TypeScript](https://www.typescriptlang.org/docs/handbook/basic-types.html)
- [Python Type Hints](https://www.python.org/dev/peps/pep-0484/)
- [Puppet Types](/puppetlabs/puppet-specifications/blob/master/language/types_values_variables.md)

### Language syntax
Dgo defines a [type language of its own](docs/types.md) which is close to Go itself. A parser and a stringifier are
provided for this syntax. New parsers and stringifiers can be added to support other syntaxes. 

### Type Assignability
As with go reflect, types can be compared for assignability. A type is assignable from another type if the other
type is equally or more restricitive, e.g. the type `int` is assignable from the range `0..10` (and all other
integer ranges). The type `string` is assignable from the pattern `/abc/`, the type `"a"|"b"|"c"` or
any other type that restricts a string. A length constrained `string[10,20]` is assignable from `string[12,17]`, etc.

### Type Instance check
A type can be used to validate if a value is an instance of that type. The integer `3` is an instance of the
range type `1..8`, the string `"abc"` is an instance of the pattern type `/b/`, etc.

### The type of a value
All values have a type that is backed by the value itself. The type will consider its value, and only that value,
to be an instance of itself. E.g. string "hello" is represented by the type `"hello"`. That type in turn is assignable
to `string`, `string[5]`, `string[0,10]`, `"hello"|"goodbye"`, but it is not assignable to `string[0,4]` or
`"hi"|"bye"`. In other words, the value type is assignable to another type if the value that it represents is an
instance of that other type.

## Immutability

Non primitives in Go (array, slice, map, struct) are mutable, and it's the programmer's responsibility to ensure that
access to such values are synchronized when they are accessed from multiple go routines.

Dgo guarantees that all values can be 100% immutable by exposing all values through interfaces and hiding the
implementation, thus enabling concurrency safe coding without the need to synchronized use of shared resources
using mutexes.

An`Array` or a `Map` can be created as a mutable collection but can be made immutable by calling the method `Freeze()`
or `Copy(true)` (argument `true` requests a frozen copy). Both calls are recursive and ensures that the collection
and all its contained values are frozen. `Freeze` performs an in-place recursive freeze of all values while `Copy`
will copy unfrozen objects before freezing them to ensure that the original and all its contained values
are not frozen as a consequence of the call.

A frozen object can never be unfrozen. The only way to resume mutability is to do `Copy(false)` which returns a
mutable copy.

## Extend the type of struct fields using tags
A Go struct can be accessed like a `dgo.Map` and the dgo type of the fields of a struct can be set using
go field tags, e.g.

```go
type Example struct {
  A map[string]string `dgo:"map[string]/^[a-z].*/"` // map of strings where the first character is a-z
  B int `dgo:"0..999"` // map of integers in the inclusive range 0 to 999
}
```

The type declared in the dgo tag must be assignable to the actual field type.

## Serialization
Support for JSON is built in to the Dgo module. Support for [gob](https://golang.org/pkg/encoding/gob/) is in
the pipeline.

Transformations between dgo and [cty](https://github.com/zclconf/go-cty) is provided by the
[dgocty](https://github.com/lyraproj/dgocty) module 

Transformations between dgo and [pcore](https://github.com/lyraproj/pcore) is provided by the
[pcore](https://github.com/lyraproj/dgopcore) module 

## Encapsulation

It's often desirable to encapsulate common behavior of values in a way that relieves the programmer from trivial
concerns. For instance:

- In Go, you cannot define a generic behavior for equality comparison. It's either `==` or `reflect.DeepEqual()`
 and both have limitations. `==` cannot compare slices or structs containing slices. The DeepEqual method compares
 _all_ fields of a struct, exported and unexported. Dgo solves this by letting all values implement the `Equals()`
 method. This method is then used throughout the Dgo framework.
- Go has no concept of a hash code. Keys in hashes may only be values that are considered comparable on Go. Dgo
solves this by letting all values implement the `HashCode()` method. The hash code can be used for several purposes,
including map keys and computing unique sets of values.
- Since Go has generic value (besides the `interface{}`, it provides no way to specify generic natural ordering. Dgo
 provides a `Comparable` interface.

## Extendability

The functionality of the Dgo basic types is exposed through interfaces to enable the type/value system to be expanded.

## Performance

Dgo is designed with performance in mind. The Map implementation is, although it's generic, actually faster
than a `map[string]interface{}`. It is a lot faster than a `map[interface{}]interface{}` required if you want to
use dynamic keyes with Go. Nevertheless, the Dgo Map can use any value as a key, even arrays, maps, and types.

The dgo.Value implementation for primitives like Bool, Integer, and Float are just redefined Go types, and as such,
they consume the same amount of memory and reuses the same logic i.e.: 
```go
type integer int64 // implements dgo.Integer
type float float64 // implements dgo.Float
type boolean bool  // implements dgo.Boolean
```

The dgo.String is different because it caches the hash code once it has been computed. A string is stored as:
```go
type hstring struct {
  s string
  h int
}
```

## How to get involved
We value your opinion very much. Please don't hesitate to reach out. Opinions, ideas, and contributions are more
than welcome. Create an [issue](../../issues), or file a [PR](../../pulls). 

## Pipeline
- Go [gob](https://golang.org/pkg/encoding/gob/) support to enable full binary exchange of values and types.
- Distributed type aliases (aliases using URI reference)
- Type extension, i.e. how a type such as a map with explicit associations can be extended by another type.
