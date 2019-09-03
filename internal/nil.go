package internal

import (
	"encoding/json"

	"gopkg.in/yaml.v3"

	"github.com/lyraproj/got/util"

	"github.com/lyraproj/got/dgo"
)

type (
	// nilType represents an nil value without constraints (-1), constrained to false (0) or constrained to true(1)
	nilType int

	nilValue int
)

func (nilValue) AppendTo(w *util.Indenter) {
	w.Append(`null`)
}

func (nilValue) CompareTo(other interface{}) (int, bool) {
	if Nil == other || nil == other {
		return 0, true
	}
	return -1, true
}

func (nilValue) HashCode() int {
	return 131
}

func (nilValue) Equals(other interface{}) bool {
	return Nil == other || nil == other
}

func (nilValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}

func (nilValue) MarshalYAML() (interface{}, error) {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!!null`, Value: `null`}, nil
}

func (nilValue) String() string {
	return `null`
}

func (nilValue) Type() dgo.Type {
	return DefaultNilType
}

const DefaultNilType = nilType(0)

const Nil = nilValue(0)

func (t nilType) Assignable(ot dgo.Type) bool {
	_, ok := ot.(nilType)
	return ok || CheckAssignableTo(nil, ot, t)
}

func (t nilType) Equals(v interface{}) bool {
	return t == v
}

func (t nilType) HashCode() int {
	return int(1 + dgo.IdNil)
}

func (t nilType) Instance(v interface{}) bool {
	return Nil == v || nil == v
}

func (t nilType) Type() dgo.Type {
	return &metaType{t}
}

func (t nilType) String() string {
	return `nil`
}

func (t nilType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.IdNil
}
