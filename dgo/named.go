package dgo

type (
	// Constructor is a function that creates a Value base on an initializer Map
	Constructor func(Value) Value

	// InitArgExtractor is a function that can extract initializer arguments from an
	// instance.
	InitArgExtractor func(Value) Value

	// NamedType is implemented by types that are named and made available using an AliasMap
	NamedType interface {
		Type

		// ExtractInitArg extracts the initializer argument from an instance.
		ExtractInitArg(Value) Value

		// Name returns the name of this type
		Name() string

		// New creates instances of this type.
		New(Value) Value

		// ValueString returns the given value as a string. The Value must be an instance of this type.
		ValueString(value Value) string
	}
)
