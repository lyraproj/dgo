package internal

import (
	"math"

	"gopkg.in/yaml.v3"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

type (
	Float float64

	floatType int

	exactFloatType float64

	floatRangeType struct {
		min float64
		max float64
	}
)

const DefaultFloatType = floatType(0)

func FloatRangeType(min, max float64) dgo.FloatRangeType {
	if min == max {
		return exactFloatType(min)
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	if min == -math.MaxFloat64 && max == math.MaxFloat64 {
		return DefaultFloatType
	}
	return &floatRangeType{min: min, max: max}
}

func (t *floatRangeType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case exactFloatType:
		return t.IsInstance(float64(ot))
	case *floatRangeType:
		return t.min <= ot.min && ot.max <= t.max
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *floatRangeType) Equals(other interface{}) bool {
	if ot, ok := other.(*floatRangeType); ok {
		return *t == *ot
	}
	return false
}

func (t *floatRangeType) HashCode() int {
	h := int(dgo.TiFloatRange)
	if t.min > 0 {
		h = h*31 + int(t.min)
	}
	if t.max < math.MaxInt64 {
		h = h*31 + int(t.max)
	}
	return h
}

func (t *floatRangeType) Instance(value interface{}) bool {
	f, ok := ToFloat(value)
	return ok && t.IsInstance(f)
}

func (t *floatRangeType) IsInstance(value float64) bool {
	return t.min <= value && value <= t.max
}

func (t *floatRangeType) Max() float64 {
	return t.max
}

func (t *floatRangeType) Min() float64 {
	return t.min
}

func (t *floatRangeType) String() string {
	return TypeString(t)
}

func (t *floatRangeType) Type() dgo.Type {
	return &metaType{t}
}

func (t *floatRangeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloatRange
}

func (t exactFloatType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(exactFloatType); ok {
		return t == ot
	}
	return CheckAssignableTo(nil, other, t)
}

func (t exactFloatType) Equals(other interface{}) bool {
	return t == other
}

func (t exactFloatType) HashCode() int {
	return Float(t).HashCode() * 3
}

func (t exactFloatType) Instance(value interface{}) bool {
	f, ok := ToFloat(value)
	return ok && float64(t) == f
}

func (t exactFloatType) IsInstance(value float64) bool {
	return float64(t) == value
}

func (t exactFloatType) Max() float64 {
	return float64(t)
}

func (t exactFloatType) Min() float64 {
	return float64(t)
}

func (t exactFloatType) Type() dgo.Type {
	return &metaType{t}
}

func (t exactFloatType) String() string {
	return TypeString(t)
}

func (t exactFloatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloatExact
}

func (t exactFloatType) Value() dgo.Value {
	v := (Float)(t)
	return v
}

func (t floatType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case floatType, exactFloatType, *floatRangeType:
		return true
	}
	return false
}

func (t floatType) Equals(other interface{}) bool {
	_, ok := other.(floatType)
	return ok
}

func (t floatType) HashCode() int {
	return int(dgo.TiFloat)
}

func (t floatType) Instance(value interface{}) bool {
	_, ok := ToFloat(value)
	return ok
}

func (t floatType) IsInstance(value float64) bool {
	return true
}

func (t floatType) Max() float64 {
	return math.MaxFloat64
}

func (t floatType) Min() float64 {
	return -math.MaxFloat64
}

func (t floatType) String() string {
	return TypeString(t)
}

func (t floatType) Type() dgo.Type {
	return &metaType{t}
}

func (t floatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloat
}

func (v Float) Type() dgo.Type {
	return exactFloatType(v)
}

func (v Float) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	if ov, isFloat := ToFloat(other); isFloat {
		fv := float64(v)
		switch {
		case fv > ov:
			r = 1
		case fv < ov:
			r = -1
		default:
			r = 0
		}
		return
	}
	if oi, isInt := ToInt(other); isInt {
		fv := float64(v)
		ov := float64(oi)
		switch {
		case fv > ov:
			r = 1
		case fv < ov:
			r = -1
		default:
			r = 0
		}
		return
	}
	if other == Nil || other == nil {
		r = 1
	} else {
		ok = false
	}
	return
}

func (v Float) Equals(other interface{}) bool {
	f, ok := ToFloat(other)
	return ok && float64(v) == f
}

func (v Float) HashCode() int {
	return int(v)
}

func (v Float) MarshalYAML() (interface{}, error) {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!!float`, Value: v.String()}, nil
}

func (v Float) String() string {
	return util.Ftoa(float64(v))
}

func (v Float) ToFloat() float64 {
	return float64(v)
}

func (v Float) ToInt() int64 {
	return int64(v)
}

func (v Float) GoFloat() float64 {
	return float64(v)
}

func ToFloat(value interface{}) (v float64, ok bool) {
	ok = true
	switch value := value.(type) {
	case Float:
		v = float64(value)
	case float64:
		v = value
	case float32:
		v = float64(value)
	default:
		ok = false
	}
	return
}
