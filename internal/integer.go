package internal

import (
	"fmt"
	"math"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/lyraproj/dgo/dgo"
)

type (
	Integer int64

	integerType int

	exactIntegerType int64

	integerRangeType struct {
		min int64
		max int64
	}
)

const DefaultIntegerType = integerType(0)

func IntegerRangeType(min, max int64) dgo.IntegerRangeType {
	if min == max {
		return exactIntegerType(min)
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	if min == math.MinInt64 && max == math.MaxInt64 {
		return DefaultIntegerType
	}
	return &integerRangeType{min: min, max: max}
}

func (t *integerRangeType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case exactIntegerType:
		return t.IsInstance(int64(ot))
	case *integerRangeType:
		return t.min <= ot.min && ot.max <= t.max
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *integerRangeType) Equals(other interface{}) bool {
	if ot, ok := other.(*integerRangeType); ok {
		return *t == *ot
	}
	return false
}

func (t *integerRangeType) HashCode() int {
	h := int(dgo.IdIntegerRange)
	if t.min > 0 {
		h = h*31 + int(t.min)
	}
	if t.max < math.MaxInt64 {
		h = h*31 + int(t.max)
	}
	return h
}

func (t *integerRangeType) Instance(value interface{}) bool {
	if ov, ok := ToInt(value); ok {
		return t.IsInstance(ov)
	}
	return false
}

func (t *integerRangeType) IsInstance(value int64) bool {
	return t.min <= value && value <= t.max
}

func (t *integerRangeType) Max() int64 {
	return t.max
}

func (t *integerRangeType) Min() int64 {
	return t.min
}

func (t *integerRangeType) String() string {
	return TypeString(t)
}

func (t *integerRangeType) Type() dgo.Type {
	return &metaType{t}
}

func (t *integerRangeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.IdIntegerRange
}

func (t exactIntegerType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(exactIntegerType); ok {
		return t == ot
	}
	return CheckAssignableTo(nil, other, t)
}

func (t exactIntegerType) Equals(other interface{}) bool {
	return t == other
}

func (t exactIntegerType) HashCode() int {
	return Integer(t).HashCode() * 5
}

func (t exactIntegerType) Instance(value interface{}) bool {
	ov, ok := ToInt(value)
	return ok && int64(t) == ov
}

func (t exactIntegerType) IsInstance(value int64) bool {
	return int64(t) == value
}

func (t exactIntegerType) Max() int64 {
	return int64(t)
}

func (t exactIntegerType) Min() int64 {
	return int64(t)
}

func (t exactIntegerType) String() string {
	return TypeString(t)
}

func (t exactIntegerType) Type() dgo.Type {
	return &metaType{t}
}

func (t exactIntegerType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.IdIntegerExact
}

func (t exactIntegerType) Value() dgo.Value {
	v := (Integer)(t)
	return v
}

func (t integerType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case integerType, exactIntegerType, *integerRangeType:
		return true
	}
	return false
}

func (t integerType) Equals(other interface{}) bool {
	_, ok := other.(integerType)
	return ok
}

func (t integerType) HashCode() int {
	return int(dgo.IdInteger)
}

func (t integerType) Instance(value interface{}) bool {
	switch value.(type) {
	case dgo.Integer, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return true
	}
	return false
}

func (t integerType) IsInstance(value int64) bool {
	return true
}

func (t integerType) Max() int64 {
	return math.MaxInt64
}

func (t integerType) Min() int64 {
	return math.MinInt64
}

func (t integerType) String() string {
	return TypeString(t)
}

func (t integerType) Type() dgo.Type {
	return &metaType{t}
}

func (t integerType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.IdInteger
}

func (v Integer) Type() dgo.Type {
	return exactIntegerType(v)
}

func (v Integer) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	if oi, isInt := ToInt(other); isInt {
		mv := int64(v)
		if mv > oi {
			r = 1
		} else if mv < oi {
			r = -1
		} else {
			r = 0
		}
	} else if ov, isFloat := ToFloat(other); isFloat {
		fv := float64(v)
		if fv > ov {
			r = 1
		} else if fv < ov {
			r = -1
		} else {
			r = 0
		}
	} else if other == Nil || other == nil {
		r = 1
	} else {
		ok = false
	}
	return
}

func (v Integer) Equals(other interface{}) bool {
	i, ok := ToInt(other)
	return ok && int64(v) == i
}

func (v Integer) HashCode() int {
	return int(v ^ (v >> 32))
}

func (v Integer) MarshalYAML() (interface{}, error) {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!!int`, Value: v.String()}, nil
}

func (v Integer) String() string {
	return strconv.Itoa(int(v))
}

func (v Integer) ToFloat() float64 {
	return float64(v)
}

func (v Integer) ToInt() int64 {
	return int64(v)
}

func (v Integer) GoInt() int64 {
	return int64(v)
}

func ToInt(value interface{}) (v int64, ok bool) {
	ok = true
	switch value := value.(type) {
	case Integer:
		v = int64(value)
	case int:
		v = int64(value)
	case int64:
		v = value
	case int32:
		v = int64(value)
	case int16:
		v = int64(value)
	case int8:
		v = int64(value)
	case uint:
		if value == math.MaxUint64 {
			panic(fmt.Errorf(`value %d overflows int64`, value))
		}
		v = int64(value)
	case uint64:
		if value == math.MaxUint64 {
			panic(fmt.Errorf(`value %d overflows int64`, value))
		}
		v = int64(value)
	case uint32:
		v = int64(value)
	case uint16:
		v = int64(value)
	case uint8:
		v = int64(value)
	default:
		ok = false
	}
	return
}
