package dgo

type (
	// AliasProvider replaces aliases with their concrete type.
	//
	// The parser uses this interface to perform in-place replacement of aliases
	AliasProvider interface {
		Replace(Type) Type
	}

	// AliasContainer is implemented by types that can contain other types.
	//
	// The parser uses this interface to perform in-place replacement of aliases
	AliasContainer interface {
		Resolve(AliasProvider)
	}

	// An AliasMap maps names to types and vice versa.
	AliasMap interface {
		// GetName returns the name for the given type and returns it or nil if no
		// type has been added using that name has been added
		GetName(t Type) String

		// GetType returns the type with the given name and returns it or nil if no
		// such type has been added with that name.
		GetType(n String) Type

		// Add adds the type t with the given name to this map
		Add(t Type, name String)
	}

	// ExactType is implemented by types that match exactly one value
	ExactType interface {
		Type

		Value() Value
	}

	// Named is implemented by named types such as the StructMap
	Named interface {
		Name() string
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
		DeepInstance(guard RecursionGuard, value interface{}) bool
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
