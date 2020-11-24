package internal

import (
	"math/big"
	"reflect"
	"time"

	"github.com/lyraproj/dgo/dgo"
)

type (
	durationType int

	durationVal time.Duration
)

// DefaultDurationType is the unconstrained Duration type
const DefaultDurationType = durationType(0)

var reflectDurationType = reflect.TypeOf(time.Duration(0))

func (t durationType) Assignable(ot dgo.Type) bool {
	switch ot.(type) {
	case durationVal, durationType:
		return true
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t durationType) Equals(v interface{}) bool {
	return t == Value(v)
}

func (t durationType) HashCode() dgo.Hash {
	return dgo.Hash(dgo.TiDuration)
}

func (t durationType) Instance(v interface{}) bool {
	switch v.(type) {
	case durationVal, time.Duration:
		return true
	}
	return false
}

func (t durationType) New(arg dgo.Value) dgo.Value {
	return newDuration(t, arg)
}

func (t durationType) ReflectType() reflect.Type {
	return reflectDurationType
}

func (t durationType) String() string {
	return TypeString(t)
}

func (t durationType) Type() dgo.Type {
	return MetaType(t)
}

func (t durationType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiDuration
}

func newDuration(t dgo.Type, arg dgo.Value) dgo.Duration {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`duration`, 1, 1)
		arg = args.Get(0)
	}
	var tv dgo.Duration
	switch arg := arg.(type) {
	case dgo.Duration:
		tv = arg
	case dgo.Integer:
		tv = durationVal(arg.GoInt())
	case dgo.Float:
		tv = durationVal(arg.GoFloat() * 1_000_000_000.0)
	case dgo.String:
		tv = DurationFromString(arg.GoString())
	default:
		panic(illegalArgument(`duration`, `duration|int|float|string`, []interface{}{arg}, 0))
	}
	if !t.Instance(tv) {
		panic(IllegalAssignment(t, tv))
	}
	return tv
}

// Duration returns the given duration as a dgo.Duration
func Duration(ts time.Duration) dgo.Duration {
	return durationVal(ts)
}

// DurationFromString returns the given duration string as a dgo.Duration. The string must be parsable
// by the time.ParseDuration() function. The function will panic if the given string cannot be parsed.
func DurationFromString(s string) dgo.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return durationVal(d)
}

func (v durationVal) Equals(other interface{}) bool {
	switch ov := other.(type) {
	case durationVal:
		return v == ov
	case time.Duration:
		return time.Duration(v) == ov
	}
	return false
}

func (v durationVal) Float() dgo.Float {
	return floatVal(v.SecondsWithFraction())
}

func (v durationVal) GoDuration() time.Duration {
	return time.Duration(v)
}

func (v durationVal) HashCode() dgo.Hash {
	return dgo.Hash(v)
}

func (v durationVal) Integer() dgo.Integer {
	return intVal(v)
}

func (v durationVal) ReflectTo(value reflect.Value) {
	d := time.Duration(v)
	switch value.Kind() {
	case reflect.Interface:
		value.Set(reflect.ValueOf(d))
	case reflect.Ptr:
		value.Set(reflect.ValueOf(&d))
	default:
		value.SetInt(int64(d))
	}
}

func (v durationVal) SecondsWithFraction() float64 {
	return float64(v) / 1_000_000_000.0
}

func (v durationVal) String() string {
	return TypeString(v)
}

func (v durationVal) ToBigFloat() *big.Float {
	return big.NewFloat(v.SecondsWithFraction())
}

func (v durationVal) ToBigInt() *big.Int {
	return big.NewInt(int64(v))
}

func (v durationVal) ToFloat() (float64, bool) {
	return v.SecondsWithFraction(), true
}

func (v durationVal) ToInt() (int64, bool) {
	return int64(v), true
}

func (v durationVal) ToUint() (uint64, bool) {
	if v >= 0 {
		return uint64(v), true
	}
	return 0, false
}

func (v durationVal) Type() dgo.Type {
	return v
}

// Duration exact type implementation

func (v durationVal) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v durationVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v durationVal) Generic() dgo.Type {
	return DefaultDurationType
}

func (v durationVal) New(arg dgo.Value) dgo.Value {
	return newDuration(v, arg)
}

func (v durationVal) ReflectType() reflect.Type {
	return reflectDurationType
}

func (v durationVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiDurationExact
}
