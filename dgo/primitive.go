/*
Package dgo contains all interfaces for the dgo types and values.

Other packages

https://godoc.org/github.com/lyraproj/dgo/typ contains the default dgo type constants
https://godoc.org/github.com/lyraproj/dgo/newtype contains methods to create types
https://godoc.org/github.com/lyraproj/dgo/vf contains methods to create Dgo values fom go values
*/
package dgo

import (
	"fmt"
	"reflect"
	"regexp"
)

type (
	// A Value represents an immutable value of some Type
	//
	// Values should not be compared using == (depending on value type, it may result in a panic: runtime error:
	//  comparing uncomparable). Instead, the method Equal or the function Identical should be used.
	//
	// There is no Value representation of "null" other than comparing to nil. The "undefined" found in some languages
	// such as TypeScript doesn't exist either. Instead, places where it should be relevant, such as when examining
	// if a Map contains a Value (or nil) for a certain key, methods will return a Value together with a bool that
	// indicates if a mapping was present or not.
	Value interface {
		fmt.Stringer

		// Type returns the type of this value
		Type() Type

		// Equals returns true if this value if equal to the given value. For complex objects, this
		// comparison is deep.
		Equals(other interface{}) bool

		// HashCode returns the computed hash code of the value
		HashCode() int
	}

	// A Type describes an immutable Value. The Type is in itself also a Value
	Type interface {
		Value

		// Assignable returns true if a variable or parameter of this type can be hold a value of the other type
		Assignable(other Type) bool

		// Instance returns true if the value is an instance of this type
		Instance(value interface{}) bool

		// TypeIdentifier returns a unique identifier for this type. The TypeIdentifier is intended to be used by
		// decorators providing string representation of the type
		TypeIdentifier() TypeIdentifier
	}

	// Number is implemented by Float and Integer implementations
	Number interface {
		// ToInt returns this number as an int64
		ToInt() int64

		// ToFloat returns this number as an float64
		ToFloat() float64
	}

	// Integer value is an int64 that implements the Value interface
	Integer interface {
		Value
		Number
		Comparable

		// GoInt returns the Go native representation of this value
		GoInt() int64
	}

	// IntegerRangeType describes integers that are within an inclusive range
	IntegerRangeType interface {
		Type

		// IsInstance returns true if the given int64 is an instance of this type
		IsInstance(int64) bool

		// Max returns the maximum constraint
		Max() int64

		// Min returns the minimum constraint
		Min() int64
	}

	// Float value is a float64 that implements the Value interface
	Float interface {
		Value
		Number
		Comparable

		// GoFloat returns the Go native representation of this value
		GoFloat() float64
	}

	// FloatRangeType describes floating point numbers that are within an inclusive range
	FloatRangeType interface {
		Type

		// IsInstance returns true if the given float64 is an instance of this type
		IsInstance(float64) bool

		// Max returns the maximum constraint
		Max() float64

		// Min returns the minimum constraint
		Min() float64
	}

	// String value is a string that implements the Value interface
	String interface {
		Value
		Comparable

		// GoString returns the Go native representation of this value
		GoString() string
	}

	// Regexp value is a string that implements the Value interface
	Regexp interface {
		Value

		// GoRegexp returns the Go native representation of this value
		GoRegexp() *regexp.Regexp
	}

	// StringType is a SizedType.
	StringType interface {
		SizedType
	}

	// Boolean value
	Boolean interface {
		Value
		Comparable

		// GoBool returns the Go native representation of this value
		GoBool() bool
	}

	// BooleanType matches the true and false literals
	BooleanType interface {
		Type

		// IsInstance returns true if the Go native value is represented by this type
		IsInstance(value bool) bool
	}

	// NativeType is the type for all Native values
	NativeType interface {
		Type

		// GoType returns the reflect.Type
		GoType() reflect.Type
	}

	// Native is a wrapper of a runtime value such as a chan or a pointer for which there is no proper immutable Value
	// representation
	Native interface {
		Value
		Freezable

		// GoValue returns the Go native representation of this value
		GoValue() interface{}
	}

	// Comparable imposes natural ordering on its implementations. A Comparable is only comparable to other
	// values of its own type with the exception of Nil which is less than everything else and the special
	// case when Integer is compared to Float. Such a comparison will convert the Integer to a Float.
	Comparable interface {
		// CompareTo compares this value with the given value for order. Returns a negative integer, zero, or a positive
		// integer as this value is less than, equal to, or greater than the specified object and a bool that indicates
		// if the comparison was at all possible.
		CompareTo(other interface{}) (int, bool)
	}

	// RecursionGuard guards against endless recursion when checking if one deep type is assignable to another.
	// A Hit is detected once both a and b has been added more than once (can be on separate calls to Append). The
	// RecursionGuard is in itself immutable.
	RecursionGuard interface {
		// Append creates a new RecursionGuard guaranteed to contain both a and b. The new instance is returned.
		Append(a, b Value) RecursionGuard

		// Hit returns true if both a and b has been appended more than once.
		Hit() bool

		// Swap returns the guard with its two internal guards for a and b swapped.
		Swap() RecursionGuard
	}
)
