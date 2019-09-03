package internal

import (
	"github.com/lyraproj/got/dgo"
	"gopkg.in/yaml.v3"
)

type (
	// booleanType represents an boolean without constraints (-1), constrained to false (0) or constrained to true(1)
	booleanType int

	Boolean bool
)

const DefaultBooleanType = booleanType(-1)
const FalseType = booleanType(0)
const TrueType = booleanType(1)

const True = Boolean(true)

const False = Boolean(false)

func (t booleanType) Assignable(ot dgo.Type) bool {
	if ob, ok := ot.(booleanType); ok {
		return t < 0 || t == ob
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t booleanType) Equals(v interface{}) bool {
	return t == v
}

func (t booleanType) HashCode() int {
	return int(t.TypeIdentifier())
}

func (t booleanType) Instance(v interface{}) bool {
	if bv, ok := v.(Boolean); ok {
		return t.IsInstance(bool(bv))
	}
	if bv, ok := v.(bool); ok {
		return t.IsInstance(bv)
	}
	return false
}

func (t booleanType) IsInstance(v bool) bool {
	return t < 0 || (t == 1) == v
}

func (t booleanType) String() string {
	return TypeString(t)
}

func (t booleanType) Type() dgo.Type {
	return &metaType{t}
}

func (t booleanType) TypeIdentifier() dgo.TypeIdentifier {
	switch t {
	case FalseType:
		return dgo.IdFalse
	case TrueType:
		return dgo.IdTrue
	default:
		return dgo.IdBoolean
	}
}

func (v Boolean) GoBool() bool {
	return bool(v)
}

func (v Boolean) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	switch ov := other.(type) {
	case Boolean:
		r = 0
		if v {
			if !ov {
				r = 1
			}
		} else if ov {
			r = -1
		}
	case nilValue:
		r = 1
	default:
		ok = false
	}
	return
}

func (v Boolean) Equals(other interface{}) bool {
	if ov, ok := other.(Boolean); ok {
		return v == ov
	}
	if ov, ok := other.(bool); ok {
		return bool(v) == ov
	}
	return false
}

func (v Boolean) HashCode() int {
	if v {
		return 1231
	}
	return 1237
}

func (v Boolean) MarshalYAML() (interface{}, error) {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!!bool`, Value: v.String()}, nil
}

func (v Boolean) String() string {
	if v {
		return `true`
	}
	return `false`
}

func (v Boolean) Type() dgo.Type {
	if v {
		return booleanType(1)
	}
	return booleanType(0)
}
