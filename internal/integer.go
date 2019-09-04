package internal

import (
	"fmt"
	"math"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// intVal is an int64 that implements the dgo.Value interface
	intVal int64

	integerType int

	exactIntegerType int64

	integerRangeType struct {
		min int64
		max int64
	}
)

// DefaultIntegerType is the unconstrained Integer type
const DefaultIntegerType = integerType(0)

// IntegerRangeType returns a dgo.IntegerRangeType that is limited to the inclusive range given by min and max
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
	h := int(dgo.TiIntegerRange)
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
	return dgo.TiIntegerRange
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
	return intVal(t).HashCode() * 5
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
	return dgo.TiIntegerExact
}

func (t exactIntegerType) Value() dgo.Value {
	v := (intVal)(t)
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
	return int(dgo.TiInteger)
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
	return dgo.TiInteger
}

// Integer returns the dgo.Integer for the given int64
func Integer(v int64) dgo.Integer {
	return intVal(v)
}

func (v intVal) Type() dgo.Type {
	return exactIntegerType(v)
}

func (v intVal) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	if oi, isInt := ToInt(other); isInt {
		mv := int64(v)
		switch {
		case mv > oi:
			r = 1
		case mv < oi:
			r = -1
		default:
			r = 0
		}
		return
	}
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
	if other == Nil || other == nil {
		r = 1
	} else {
		ok = false
	}
	return
}

func (v intVal) Equals(other interface{}) bool {
	i, ok := ToInt(other)
	return ok && int64(v) == i
}

func (v intVal) HashCode() int {
	return int(v ^ (v >> 32))
}

func (v intVal) MarshalYAML() (interface{}, error) {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!!int`, Value: v.String()}, nil
}

func (v intVal) String() string {
	return strconv.Itoa(int(v))
}

func (v intVal) ToFloat() float64 {
	return float64(v)
}

func (v intVal) ToInt() int64 {
	return int64(v)
}

func (v intVal) GoInt() int64 {
	return int64(v)
}

// ToInt returns the given value as a int64 if, and only if, the value type is one of the go int types. An
// additional boolean is returned to indicate if that was the case or not.
func ToInt(value interface{}) (v int64, ok bool) {
	ok = true
	switch value := value.(type) {
	case intVal:
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
