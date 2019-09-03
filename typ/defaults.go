package typ

import "github.com/lyraproj/got/internal"

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
const Binary = internal.DefaultBinaryType

// String is a type that represents all strings
const String = internal.DefaultStringType

// Error is a type that represents all implementation of error
const Error = internal.DefaultErrorType

// Native is a type that represents all Native values
var Native = internal.DefaultNativeType
