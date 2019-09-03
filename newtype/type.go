package newtype

import (
	"reflect"

	"github.com/lyraproj/got/dgo"
	"github.com/lyraproj/got/internal"
)

// FromReflected returns teh got.Type that represents the given reflected type
func FromReflected(vt reflect.Type) dgo.Type {
	return internal.TypeFromReflected(vt)
}

// ParseFile parses the given content into a got.Type.
func Parse(content string) dgo.Type {
	return internal.Parse(content)
}

// ParseFile parses the given content into a got.Type. The filename is used in error messages.
func ParseFile(fileName, content string) dgo.Type {
	return internal.ParseFile(fileName, content)
}
