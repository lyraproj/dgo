package vf

import (
	"encoding/json"

	"github.com/lyraproj/got/dgo"
	"github.com/lyraproj/got/internal"
	"github.com/lyraproj/got/util"
)

// MarshalJSON returns the JSON encoding for the given got.Value
func MarshalJSON(v interface{}) ([]byte, error) {
	// Default Indentable output is JSON
	if i, ok := v.(util.Indentable); ok {
		return []byte(util.ToString(i)), nil
	}
	return json.Marshal(v)
}

// UnmarshalJSON decodes the JSON representation of the given bytes into a got.Value
func UnmarshalJSON(b []byte) (dgo.Value, error) {
	return internal.UnmarshalJSON(b)
}
