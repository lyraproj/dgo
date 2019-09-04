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
	TiNil = TypeIdentifier(iota)
	TiAny
	TiMeta
	TiBoolean
	TiFalse
	TiTrue
	TiInteger
	TiIntegerExact
	TiIntegerRange
	TiFloat
	TiFloatExact
	TiFloatRange
	TiBinary
	TiBinaryExact
	TiString
	TiStringExact
	TiStringPattern
	TiStringSized

	TiRegexp
	TiRegexpExact
	TiNative

	TiArray
	TiArrayExact
	TiArrayElementSized
	TiElementsExact

	TiTuple

	TiMap
	TiMapExact
	TiMapValuesExact
	TiMapKeysExact
	TiMapSized
	TiMapEntry
	TiMapEntryExact
	TiStruct

	TiNot
	TiAllOf
	TiAnyOf
	TiOneOf

	TiError
)
