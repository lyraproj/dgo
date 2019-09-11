# dgo (Dynamic Go)

[![Go Report Card](https://goreportcard.com/badge/github.com/lyraproj/dgo)](https://goreportcard.com/report/github.com/lyraproj/dgo)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/lyraproj/dgo/dgo)

Dgo's main objectives are: Type Constraints, Immutability, Collections, Serialization, Encapsulation, Extendability,
and Performance. It is designed to make working with dynamic values an effortless and type safe task.

## Install
Dgo is a go module and if the Go version is < 1.13, go modules must be enabled. This is done by setting the environment
variable GO111MODULE=on before an attempt is made to install:
```sh
export GO111MODULE=on
```
To install the module under $GOPATH/src:
```sh
go get github.com/lyraproj/dgo
```

## Example of usage:
Let's assume some kind of typed parameters in YAML that the user enters like this:
```yaml
host: example.com
port: 22
```
The task is to create a user friendly description of a parameter, also in YAML, which can be used to validate
the above parameters. Something like this:
```yaml
host:
  type: string[1]
  name: sample/service_host
  required: true
port:
  type: 1..999
  name: sample/service_port
```
The value of each `type` is a [dgo type](docs/types.md). They limit the host parameter to a non empty string
and the port parameter to an integer in the range 1-999. A special `required` entry is used to denote whether
or not a parameter value must be present.

Next, the parameter descriptions must converted into a `StructType`. This type is a special variant of the
`MapType` that will constrain its values to a specified set of associations. The conversion is made using
the function `newtype.StructFromMap()`. This function requires an argument that is an instance of the type
`map[string](dgo|type|{type:dgo|type,required?:bool,...})` and it just so happens that the above friendly
parameter description fits that bill. 

A Go function that reads the description can look like this:
```go
func loadDesc(yamlData []byte) (dgo.StructType, error) {
  data := vf.MutableMap(nil)
  if err := yaml.Unmarshal(yamlData, data); err != nil {
    return nil, err
  }
  return newtype.StructFromMap(false, data), nil
}
```

The `StructType` returned from the `loadDesc()` function can now be used to validate the parameters.
 
Here is a sample function loads the type, the parameters, and then uses the type to validate the parameters. All
errors are simply printed on stdout:
```go
func validate(yamlDesc, yamlParams []byte) {
  pt, err := loadDesc(yamlDesc)
  if err != nil {
    fmt.Println(err)
    return
  }
  params := vf.MutableMap(nil)
  if err = yaml.Unmarshal(yamlParams, params); err != nil {
    fmt.Println(err)
    return
  }
  for _, err := range pt.Validate(nil, params) {
    fmt.Println(err)
  }
}
```
For more examples of how to use the above validation, please take a look at
[parameter_test.go](examples_test/parameter_test.go).

## Type Constraints

### What's wrong with Go's type system?
Go is a typed language but the types are not very descriptive. It is for instance not possible to declare a type
that corresponds only to a specific range of integers, a string that must confirm to a specific pattern, or a slice
that can contain only integers or floats. The only generic type provided by Go is the empty interface `interface{}`
which corresponds to every possible value in the system. A slice containing ints and floats must be declared as
`[]interface{}` but doing so creates a list that can also contain any other type of value. What's needed is a type
such as `[](int|float)`. Other examples can be `map[string](string|int)` (string keyed map of string
and int values), or `[]0..15` (slice of integers ranging from 0 - 15).

In many situations it is desirable to declare more restrictive type constraints.
[TypeScript](https://www.typescriptlang.org/docs/handbook/basic-types.html) has good model for such constraints.
Others can be found in [Python Type Hints](https://www.python.org/dev/peps/pep-0484/) and
[Puppet Types](/puppetlabs/puppet-specifications/blob/master/language/types_values_variables.md).

### Language syntax
Dgo defines a [type language of its own](docs/types.md) which is designed to be close to Go itself. A parser
and a stringifier are provided for this syntax. New parsers and stringifiers can be added to support other syntaxes. 

### Type Assignability
As with go reflect, types can be compared for assignability. A type is assignable from another type if the other
type is equally or more restricitive, e.g. the type `int` is assignable from the range `0..10` (and all other
integer ranges). The type `string` is assignable from the pattern `/abc/`, the type `"a"|"b"|"c"` or
any other type that restricts a string. A length constrained `string[10,20]` is assignable from `string[12,17]`, etc.

### Type Instance check
A type can be used to validate if a value is an instance of that type. The integer `3` is an instance of the
range type `1..4`, the string `"abc"` is an instance of the pattern type `/b/`, etc.

### The type of a value
All values have a type that is backed by the value itself. The type will consider its value, and only that value,
to be an instance of itself. E.g. string "hello" is represented by the type `"hello"`. That type in turn is assignable
to `string`, `string[5]`, `string[0,10]`, `"hello"|"goodbye"`, but it is not assignable to `string[0,4]` or
`"hi"|"bye"`. In other words, the value type is assignable to another type if the value that it represents is an
instance of that other type.

#### Collection types
The default type for an Array or a Map can be overridden. This is particularly useful when working with
mutable collections. If an explicit type is given to the collection, it will instead ensure that any modifications
made to it will be in conformance with that type. The default type is backed by the collection and changes
dynamically when the collection is modified.

## Immutability

Non primitives in Go (array, slice, map, struct) are mutable and it's the programmers responsibility to ensure that
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

## Serialization
Support for JSON and YAML is built in to Dgo. Support for [gob](https://golang.org/pkg/encoding/gob/) is in the
pipeline.

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
than welcome. Ping us on the [Puppet Cloudnative Slack](https://app.slack.com/client/T0AQJ2CTU/CHPSJ9L4F), create
an [issue](../../issues), or file a [PR](../../pulls). 

## Pipeline
- dgo annotations for go struct members, e.g.
    ```go
    type Input struct {
      Variables map[string]interface{} `dgo:"map[string]dgo"`
      Parameters map[string]interface{} `dgo:"map[string]{name: string[1], type: dgo, required?: bool}"`
    }
    ```
- Go [gob](https://golang.org/pkg/encoding/gob/) support to enable full binary exchange of values and types.
- Distributed type aliases (aliases using URI reference)
- Type extension, i.e. how a type such as a map with explicit associations can be extended by another type.
