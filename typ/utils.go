package typ

import (
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
)

// Generic returns the generic form of the given type. All non exact types are considered generic
// and will be returned directly. Exact types will loose information about what instance they represent
// and also range and size information. Nested types will return a generic version of the contained
// types as well.
func Generic(t dgo.Type) dgo.Type {
	return internal.Generic(t)
}
