package internal

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"

	"github.com/lyraproj/dgo/dgo"
	"github.com/tada/catch"
)

type (
	// uintVal is an uint64 that implements the dgo.Value interface
	uintVal uint64
)

// DefaultUint64Type is the unconstrained Uint64 type
var DefaultUint64Type = IntegerType(Int64(0), Uint64(math.MaxUint64), true)

var reflectUint64Type = reflect.TypeOf(uint64(0))

// Uint64 returns the dgo.Uint64 for the given uint64
func Uint64(v uint64) dgo.Uint64 {
	return uintVal(v)
}

func (v uintVal) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v uintVal) CompareTo(other interface{}) (int, bool) {
	r := 0
	ok := true

	mv := uint64(v)
	compare64 := func(fv uint64) {
		switch {
		case mv > fv:
			r = 1
		case mv < fv:
			r = -1
		}
	}

	compareBig := func(ov *big.Int) {
		if ov.IsUint64() {
			compare64(ov.Uint64())
		} else {
			r = -ov.Sign()
		}
	}

	switch ov := other.(type) {
	case nil, nilValue:
		r = 1
	case uintVal:
		compare64(uint64(ov))
	case uint:
		compare64(uint64(ov))
	case uint64:
		compare64(ov)
	case *bigIntVal:
		compareBig(ov.Int)
	case *big.Int:
		compareBig(ov)
	case dgo.Float:
		r, ok = v.Float().CompareTo(ov)
	default: // all other int types
		var iv int64
		iv, ok = ToInt(other)
		if ok {
			if iv < 0 {
				r = 1
			} else {
				compare64(uint64(iv))
			}
		}
	}
	return r, ok
}

func (v uintVal) Dec() dgo.Integer {
	return v - 1
}

func (v uintVal) Equals(other interface{}) bool {
	i, ok := ToUint(other)
	return ok && uint64(v) == i
}

func (v uintVal) Float() dgo.Float {
	return floatVal(v)
}

func (v uintVal) Format(s fmt.State, format rune) {
	doFormat(uint64(v), s, format)
}

func (v uintVal) Generic() dgo.Type {
	return DefaultUint64Type
}

func (v uintVal) GoInt() int64 {
	if v <= math.MaxInt64 {
		return int64(v)
	}
	panic(catch.Error(`UInt64.ToInt(): value %d cannot fit into an int64`, v))
}

func (v uintVal) GoUint() uint64 {
	return uint64(v)
}

func (v uintVal) HashCode() dgo.Hash {
	return dgo.Hash(v ^ (v >> 32))
}

func (v uintVal) Inc() dgo.Integer {
	return v + 1
}

func (v uintVal) Inclusive() bool {
	return true
}

func (v uintVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v uintVal) Integer() dgo.Integer {
	return v
}

func (v uintVal) Max() dgo.Integer {
	return v
}

func (v uintVal) Min() dgo.Integer {
	return v
}

func (v uintVal) New(arg dgo.Value) dgo.Value {
	return newInt(v, arg)
}

func (v uintVal) ReflectTo(value reflect.Value) {
	switch {
	case strings.HasPrefix(value.Type().Name(), "uint"):
		value.SetUint(uint64(v))
	case strings.HasPrefix(value.Type().Name(), "int"):
		value.SetInt(v.GoInt())
	default:
		valueTypeReflectTo(v, uint64(v), value)
	}
}

func (v uintVal) ReflectType() reflect.Type {
	return reflectUint64Type
}

func (v uintVal) String() string {
	return TypeString(v)
}

func (v uintVal) ToBigInt() *big.Int {
	return new(big.Int).SetUint64(uint64(v))
}

func (v uintVal) ToBigFloat() *big.Float {
	return new(big.Float).SetUint64(uint64(v))
}

func (v uintVal) ToFloat() (float64, bool) {
	return float64(v), true
}

func (v uintVal) ToInt() (int64, bool) {
	if v <= math.MaxInt64 {
		return int64(v), true
	}
	return 0, false
}

func (v uintVal) ToUint() (uint64, bool) {
	return uint64(v), true
}

func (v uintVal) Type() dgo.Type {
	return v
}

func (v uintVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiUint64Exact
}

// ToUint returns the (value as an uint64, true) if it fits into that data type, (0, false) if not
func ToUint(value interface{}) (uint64, bool) {
	ok := true
	v := uint64(0)
	switch value := value.(type) {
	case uintVal:
		v = uint64(value)
	case uint64:
		v = value
	case uint:
		v = uint64(value)
	case dgo.BigInt:
		bi := value.GoBigInt()
		if ok = bi.IsUint64(); ok {
			v = bi.Uint64()
		}
	case *big.Int:
		if ok = value.IsUint64(); ok {
			v = value.Uint64()
		}
	default:
		var iv int64
		if iv, ok = ToInt(value); iv >= 0 {
			v = uint64(iv)
		} else {
			v = 0
			ok = false
		}
	}
	return v, ok
}

func unsignedToInteger(v uint64) dgo.Integer {
	if v <= math.MaxInt64 {
		return intVal(int64(v))
	}
	return uintVal(v)
}
