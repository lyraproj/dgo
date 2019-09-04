package vf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// UnmarshalYAML decodes the YAML representation of the given bytes into a dgo.Value
func UnmarshalYAML(b []byte) (dgo.Value, error) {
	return internal.UnmarshalYAML(b)
}
