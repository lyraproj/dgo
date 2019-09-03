package dgo

type (
	// BinaryType is the type that represents a Binary value
	BinaryType interface {
		Type

		IsInstance([]byte) bool
	}

	// Binary represents a sequence of bytes. Its string representation is a strictly base64 encoded string
	Binary interface {
		Value

		// GoBytes returns a copy of the internal array to ensure immutability
		GoBytes() []byte
	}
)
