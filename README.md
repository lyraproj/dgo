# dgo (Dynamic Go)

Dgo's main objectives are: Type Constraints, Immutability, Collections, Serialization, Encapsulation, Extendability,
and Performance. It is designed to make working with dynamic values an effortless and type safe task.

## Type Constraints

Go is a typed language but the types are not very descriptive. It is for instance not possible to declare a type
that corresponds only to a specific range of integers, a string that must confirm to a specific pattern, or a slice
that can contain only integers or floats. The only generic type provided by Go is the empty interface `interface{}`
which corresponds to every possible value in the system. A list containing ints and floats must be declared as
`[]interface{}` but doing so creates a list that can also contain any other type of value. What if I want to declare
a type like `[](int|float)` (slice of int and float values) or `map[string](string|int)` (string keyed map of string
and int values), or `[]0..15` (slice of integers ranging from 0 - 15)?

In many situations it is desirable to declare more restrictive type constraints.
[TypeScript](https://www.typescriptlang.org/docs/handbook/basic-types.html) has good model for such constraints.
Others can be found in [Python Type Hints](https://www.python.org/dev/peps/pep-0484/) and
[Puppet Types](/puppetlabs/puppet-specifications/blob/master/language/types_values_variables.md).

Dgo defines a [type language of its own](docs/types.md) which is designed to be close to Go itself. A parser
and a stringifier are provided for this syntax. New parsers and stringifiers can be added to support other syntaxes. 

## Immutability

Non primitives in Go (array, slice, map, struct) are mutable and it's the programmers responsibility to ensure that
access to such values are synchronized when they are accessed from multiple go routines.

Dgo guarantees that all values can be 100% immutable by exposing all values through interfaces and hiding the
implementation, thus enabling concurrency safe coding without the need of synchronization.

#### Collections
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

## Disclaimer
Dgo is work in progress. It is far from production ready at this point. If you stumbled on this project and
managed to read this far, please don't hesitate to reach out. Opinions, ideas, and even contributions are more than
welcome. Send an email to [Thomas Hallgren](mailto:thomas.hallgren@puppet.com), ping me on the
[Puppet Cloudnative Slack](https://app.slack.com/client/T0AQJ2CTU/CHPSJ9L4F), create an [issue](../../issues), or
file a [PR](../../pulls). 
