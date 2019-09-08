package vf

import (
	"encoding/json"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/util"
)

// MarshalJSON returns the JSON encoding for the given dgo.Value
func MarshalJSON(v interface{}) ([]byte, error) {
	// Default Indentable output is JSON
	if i, ok := v.(util.Indentable); ok {
		return []byte(util.ToStringERP(i)), nil
	}
	return json.Marshal(v)
}

// UnmarshalJSON decodes the JSON representation of the given bytes into a dgo.Value
func UnmarshalJSON(b []byte) (dgo.Value, error) {
	return internal.UnmarshalJSON(b)
}
