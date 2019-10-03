package typ

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// AllOf is the default AllOf type. Since it contains no types, everything is
// assignable to it.
var AllOf = internal.DefaultAllOfType

// AnyOf is the default AnyOf type. Since it contains no types, nothing is
// assignable to it except all AnyOf types.
var AnyOf = internal.DefaultAnyOfType

// OneOf is the default OneOf type. Since it contains no types, nothing is
// assignable to it except all OneOf types
var OneOf = internal.DefaultOneOfType

// Array represents all array values
const Array = internal.DefaultArrayType

// Tuple represents all Array values since it has no elements with type constraints
var Tuple = internal.DefaultTupleType

// Not represents a negated type. The default Not is negated Any so no other type
// is assignable to it.
var Not = internal.DefaultNotType

// Map is the unconstrained type. It represents all Map values
const Map = internal.DefaultMapType

// Any is a type that represents all values
const Any = internal.DefaultAnyType

// Nil is a type that represents the nil Value
const Nil = internal.DefaultNilType

// Boolean is a type that represents both true and false
const Boolean = internal.DefaultBooleanType

// False is a type that only represents the value false
const False = internal.FalseType

// True is a type that only represents the value true
const True = internal.TrueType

// Float is a type that represents all floats
const Float = internal.DefaultFloatType

// Integer is a type that represents all integers
const Integer = internal.DefaultIntegerType

// Regexp is a type that represents all regexps
const Regexp = internal.DefaultRegexpType

// Binary is a type that represents all Binary values
var Binary = internal.DefaultBinaryType

// String is a type that represents all strings
const String = internal.DefaultStringType

// DgoString is a type that represents all strings with Dgo syntax
const DgoString = internal.DefaultDgoStringType

// Error is a type that represents all implementation of error
const Error = internal.DefaultErrorType

// Native is a type that represents all Native values
var Native = internal.DefaultNativeType

// Sensitive is a type that represents Sensitive values
var Sensitive = internal.DefaultSensitiveType

// Generic returns the generic form of the given type. All non exact types are considered generic
// and will be returned directly. Exact types will loose information about what instance they represent
// and also range and size information. Nested types will return a generic version of the contained
// types as well.
func Generic(t dgo.Type) dgo.Type {
	return internal.Generic(t)
}
