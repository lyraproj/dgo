package internal

import (
	"math/big"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/tada/catch"
)

type (
	bigIntVal struct {
		*big.Int
	}

	defaultBigIntType struct {
		defaultIntegerType
	}

	bigIntType struct {
		integerType
	}
)

// DefaultBigIntType is the unconstrained Int64 type
var DefaultBigIntType = &defaultBigIntType{defaultIntegerType(dgo.TiBigInt)}

var reflectBigIntType = reflect.TypeOf(&big.Int{})

func (t *defaultBigIntType) New(arg dgo.Value) dgo.Value {
	return newBigInt(t, arg)
}

func (t *defaultBigIntType) ReflectType() reflect.Type {
	return reflectBigIntType
}

func (t *defaultBigIntType) Equals(other interface{}) bool {
	_, ok := Value(other).(*defaultBigIntType)
	return ok
}

func (t *bigIntType) New(arg dgo.Value) dgo.Value {
	return newBigInt(t, arg)
}

func (t *bigIntType) ReflectType() reflect.Type {
	return reflectBigIntType
}

func (t *bigIntType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBigIntRange
}

// BigInt returns the dgo.BigInt for the given *big.Int
func BigInt(v *big.Int) dgo.BigInt {
	return &bigIntVal{v}
}

func (v *bigIntVal) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v *bigIntVal) compare64(ov int64) int {
	r := 0
	if v.IsInt64() {
		iv := v.Int64()
		switch {
		case iv > ov:
			r = 1
		case iv < ov:
			r = -1
		}
	} else {
		r = v.Sign()
	}
	return r
}

func (v *bigIntVal) compareU64(ov uint64) int {
	r := 0
	if v.IsUint64() {
		iv := v.Uint64()
		switch {
		case iv > ov:
			r = 1
		case iv < ov:
			r = -1
		}
	} else {
		r = v.Sign()
	}
	return r
}

func (v *bigIntVal) CompareTo(other interface{}) (int, bool) {
	r := 0
	ok := true

	switch ov := other.(type) {
	case nil, nilValue:
		r = 1
	case *bigIntVal:
		r = v.Cmp(ov.Int)
	case intVal:
		r = v.compare64(int64(ov))
	case uintVal:
		r = v.compareU64(uint64(ov))
	case *big.Int:
		r = v.Cmp(ov)
	case int:
		r = v.compare64(int64(ov))
	case uint:
		r = v.compareU64(uint64(ov))
	case uint64:
		r = v.compareU64(ov)
	case float64:
		r = v.compare64(int64(ov))
	case float32:
		r = v.compare64(int64(ov))
	case *big.Float:
		bi, a := ov.Int(nil)
		r = v.Cmp(bi)
		if r == 0 {
			r = int(a)
		}
	case dgo.Number:
		r, ok = v.CompareTo(ov.Integer())
	default:
		var iv int64
		if iv, ok = ToInt(ov); ok {
			r = v.compare64(iv)
		}
	}
	return r, ok
}

func (v *bigIntVal) Dec() dgo.Integer {
	return &bigIntVal{new(big.Int).Sub(v.Int, big.NewInt(1))}
}

func (v *bigIntVal) Equals(other interface{}) bool {
	switch ov := other.(type) {
	case *bigIntVal:
		return v.Cmp(ov.Int) == 0
	case *big.Int:
		return v.Cmp(ov) == 0
	case uint:
		return v.Cmp(new(big.Int).SetUint64(uint64(ov))) == 0
	case uint64:
		return v.Cmp(new(big.Int).SetUint64(ov)) == 0
	default:
		if v.IsInt64() {
			return intVal(v.Int64()).Equals(ov)
		}
	}
	return false
}

func (v *bigIntVal) Float() dgo.Float {
	return &bigFloatVal{v.ToBigFloat()}
}

func (v *bigIntVal) Generic() dgo.Type {
	return DefaultBigIntType
}

func (v *bigIntVal) GoBigInt() *big.Int {
	return v.Int
}

func (v *bigIntVal) GoInt() int64 {
	if i, ok := v.ToInt(); ok {
		return i
	}
	panic(catch.Error(`BigInt.ToInt(): value %d cannot fit into an int64`, v))
}

func (v *bigIntVal) HashCode() dgo.Hash {
	return bigIntHash(v.Int)
}

func (v *bigIntVal) Inc() dgo.Integer {
	return &bigIntVal{new(big.Int).Add(v.Int, big.NewInt(1))}
}

func (v *bigIntVal) Inclusive() bool {
	return true
}

func (v *bigIntVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v *bigIntVal) Integer() dgo.Integer {
	return v
}

func (v *bigIntVal) Max() dgo.Integer {
	return v
}

func (v *bigIntVal) Min() dgo.Integer {
	return v
}

func (v *bigIntVal) New(arg dgo.Value) dgo.Value {
	return newBigInt(v, arg)
}

func (v *bigIntVal) ReflectTo(value reflect.Value) {
	pointerTypeReflectTo(v, reflect.ValueOf(v.Int), value)
}

func (v *bigIntVal) ReflectType() reflect.Type {
	return reflectBigIntType
}

func (v *bigIntVal) ToBigFloat() *big.Float {
	return new(big.Float).SetInt(v.Int)
}

func (v *bigIntVal) ToBigInt() *big.Int {
	return v.Int
}

func (v *bigIntVal) ToInt() (int64, bool) {
	if v.IsInt64() {
		return v.Int64(), true
	}
	return 0, false
}

func (v *bigIntVal) ToUint() (uint64, bool) {
	if v.IsUint64() {
		return v.Uint64(), true
	}
	return 0, false
}

func (v *bigIntVal) ToFloat() (float64, bool) {
	return demoteToFloat64(v.ToBigFloat())
}

func (v *bigIntVal) Type() dgo.Type {
	return v
}

func (v *bigIntVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBigIntExact
}

func bigIntHash(v *big.Int) dgo.Hash {
	hc := dgo.Hash(0)
	for _, w := range v.Bits() {
		hc = hc*31 + dgo.Hash(w)
	}
	return hc
}

func newBigInt(t dgo.Type, arg dgo.Value) (i dgo.Integer) {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`int`, 1, 2)
		if args.Len() == 2 {
			i = bigIntFromConvertible(args.Get(0), int(args.Arg(`int`, 1, radixType).(dgo.Integer).GoInt()))
		} else {
			i = bigIntFromConvertible(args.Get(0), 0)
		}
	} else {
		i = bigIntFromConvertible(arg, 0)
	}
	if !t.Instance(i) {
		panic(IllegalAssignment(t, i))
	}
	return i
}

func bigIntFromConvertible(from dgo.Value, radix int) dgo.Integer {
	switch from := from.(type) {
	case dgo.Number:
		bi := from.Integer()
		if _, ok := bi.(dgo.BigInt); !ok {
			bi = BigInt(bi.ToBigInt())
		}
		return bi
	case dgo.Boolean:
		if from.GoBool() {
			return BigInt(big.NewInt(1))
		}
		return BigInt(big.NewInt(0))
	case dgo.String:
		if bi, ok := new(big.Int).SetString(from.GoString(), radix); ok {
			return BigInt(bi)
		}
	}
	panic(catch.Error(`the value '%v' cannot be converted to an int`, from))
}
