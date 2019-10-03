package newtype

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Sensitive returns a Sensitive dgo.Type that wraps the given dgo.Type
func Sensitive(wrappedType dgo.Type) dgo.Type {
	return internal.SensitiveType(wrappedType)
}
