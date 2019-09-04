package newtype

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// FromReflected returns teh dgo.Type that represents the given reflected type
func FromReflected(vt reflect.Type) dgo.Type {
	return internal.TypeFromReflected(vt)
}

// Parse parses the given content into a dgo.Type.
func Parse(content string) dgo.Type {
	return internal.Parse(content)
}

// ParseFile parses the given content into a dgo.Type. The filename is used in error messages.
func ParseFile(fileName, content string) dgo.Type {
	return internal.ParseFile(fileName, content)
}
