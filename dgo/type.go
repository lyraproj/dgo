package dgo

type (
	// TypeIdentifier is a unique identifier for each type known to the system. The order of the TypeIdentifier
	// determines the sort order for elements that are not comparable
	TypeIdentifier int

	// ExactType is implemented by types that match exactly one value
	ExactType interface {
		Type

		Value() Value
	}

	// SizedType is implemented by types that may have a size constraint
	// such as String, Array, or Map
	SizedType interface {
		Type

		// Max returns the maximum size for instances of this type
		Max() int

		// Min returns the minimum size for instances of this type
		Min() int

		// Unbounded returns true when the type has no size constraint
		Unbounded() bool
	}

	// DeepAssignable is implemented by values that need deep Assignable comparisons.
	DeepAssignable interface {
		DeepAssignable(guard RecursionGuard, other Type) bool
	}

	// DeepInstance is implemented by values that need deep Intance comparisons.
	DeepInstance interface {
		DeepInstance(guard RecursionGuard, other Value) bool
	}

	// ReverseAssignable indicates that the check for assignable must continue by delegating to the
	// type passed as an argument to the Assignable method. The reason is that types like AllOf, AnyOf
	// OneOf or types representing exact slices or maps, might need to check if individual types are
	// assignable.
	//
	// All implementations of Assignable must take into account the argument may implement this interface
	// do a reverse by calling the CheckAssignableTo function
	ReverseAssignable interface {
		// AssignableTo returns true if a variable or parameter of the other type can be hold a value of this type.
		// All implementations of Assignable must take into account that the given type might implement this method
		// do a reverse check before returning false.
		//
		// The guard is part of the internal endless recursion mechanism and should be passed as nil unless provided
		// by a DeepAssignable caller.
		AssignableTo(guard RecursionGuard, other Type) bool
	}
)

const (
	// TiNil is the type identifier for the Nil type
	TiNil = TypeIdentifier(iota)
	// TiAny is the type identifier for the Any type
	TiAny
	// TiMeta is the type identifier for the Meta type
	TiMeta
	// TiBoolean is the type identifier for the Boolean type
	TiBoolean
	// TiFalse is the type identifier for the False type
	TiFalse
	// TiTrue is the type identifier for the True type
	TiTrue
	// TiInteger is the type identifier for the Integer type
	TiInteger
	// TiIntegerExact is the type identifier for the exact Integer type
	TiIntegerExact
	// TiIntegerRange is the type identifier for the Integer range type
	TiIntegerRange
	// TiFloat is the type identifier for the Float type
	TiFloat
	// TiFloatExact is the type identifier for the exact Float type
	TiFloatExact
	// TiFloatRange is the type identifier for the Float range type
	TiFloatRange
	// TiBinary is the type identifier for the Binary type
	TiBinary
	// TiString is the type identifier for the String type
	TiString
	// TiStringExact is the type identifier for the exact String type
	TiStringExact
	// TiStringPattern is the type identifier for the String pattern type
	TiStringPattern
	// TiStringSized is the type identifier for the size constrained String type
	TiStringSized

	// TiRegexp is the type identifier for the Regexp type
	TiRegexp
	// TiRegexpExact is the type identifier for the exact Regexp type
	TiRegexpExact
	// TiNative is the type identifier for the Native type
	TiNative

	// TiArray is the type identifier for the Array type
	TiArray
	// TiArrayExact is the type identifier for the exact Array type
	TiArrayExact
	// TiArrayElementSized is the type identifier for the element and size constrained Array type
	TiArrayElementSized
	// TiElementsExact is the type identifier for the element type of an exact Array type
	TiElementsExact

	// TiTuple is the type identifier for the Tuple type
	TiTuple

	// TiMap is the type identifier for the Map type
	TiMap
	// TiMapExact is the type identifier for exact Map type
	TiMapExact
	// TiMapValuesExact is the type identifier the value type of the exact Map type
	TiMapValuesExact
	// TiMapKeysExact is the type identifier for the key type of the exact Map type
	TiMapKeysExact
	// TiMapSized is the type identifier for the key, value, and size constrained Map type
	TiMapSized
	// TiMapEntry is the type identifier for the map entry type of a Struct type
	TiMapEntry
	// TiMapEntryExact is the type identifier the map entry type of the exact Map type
	TiMapEntryExact
	// TiStruct is the type identifier for the Struct type
	TiStruct

	// TiNot is the type identifier for the Not type
	TiNot
	// TiAllOf is the type identifier for the AllOf type
	TiAllOf
	// TiAnyOf is the type identifier for the AnyOf type
	TiAnyOf
	// TiOneOf is the type identifier for the OneOf type
	TiOneOf

	// TiError is the type identifier for for the Error type
	TiError

	// TiDgoString is the type identifier for for the DgoString type
	TiDgoString
)
