package internal

import (
	"math"
	"reflect"
	"time"

	"github.com/lyraproj/dgo/dgo"
)

type (
	timeType int

	exactTimeType time.Time

	timeVal time.Time
)

// DefaultTimeType is the unconstrainted Time type
const DefaultTimeType = timeType(0)

var reflectTimeType = reflect.TypeOf(time.Time{})

func (t timeType) Assignable(ot dgo.Type) bool {
	switch ot.(type) {
	case timeType, *exactTimeType:
		return true
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t timeType) Equals(v interface{}) bool {
	return t == v
}

func (t timeType) HashCode() int {
	return int(dgo.TiTime)
}

func (t timeType) Instance(v interface{}) bool {
	switch v.(type) {
	case *timeVal, *time.Time, time.Time:
		return true
	}
	return false
}

func (t timeType) New(arg dgo.Value) dgo.Value {
	return newTime(t, arg)
}

func (t timeType) ReflectType() reflect.Type {
	return reflectTimeType
}

func (t timeType) String() string {
	return TypeString(t)
}

func (t timeType) Type() dgo.Type {
	return &metaType{t}
}

func (t timeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiTime
}

func (t *exactTimeType) Assignable(other dgo.Type) bool {
	return t.Equals(other)
}

func (t *exactTimeType) Equals(other interface{}) bool {
	if ot, ok := other.(*exactTimeType); ok {
		return (*time.Time)(t).Equal(*(*time.Time)(ot))
	}
	return false
}

func (t *exactTimeType) Generic() dgo.Type {
	return DefaultTimeType
}

func (t *exactTimeType) HashCode() int {
	return (*timeVal)(t).HashCode()*31 + int(dgo.TiTimeExact)
}

func (t *exactTimeType) Instance(value interface{}) bool {
	switch ov := value.(type) {
	case *timeVal:
		return t.IsInstance(*(*time.Time)(ov))
	case time.Time:
		return t.IsInstance(ov)
	}
	return false
}

func (t *exactTimeType) IsInstance(tv time.Time) bool {
	return (*time.Time)(t).Equal(tv)
}

func (t *exactTimeType) New(arg dgo.Value) dgo.Value {
	return newTime(t, arg)
}

func (t *exactTimeType) ReflectType() reflect.Type {
	return reflectTimeType
}

func (t *exactTimeType) String() string {
	return TypeString(t)
}

func (t *exactTimeType) Type() dgo.Type {
	return &metaType{t}
}

func (t *exactTimeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiTimeExact
}

func (t *exactTimeType) Value() dgo.Value {
	return (*timeVal)(t)
}

func newTime(t dgo.Type, arg dgo.Value) dgo.Time {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`time`, 1, 1)
		arg = args.Get(0)
	}
	var tv dgo.Time
	switch arg := arg.(type) {
	case dgo.Time:
		tv = arg
	case dgo.Integer:
		tv = Time(time.Unix(arg.GoInt(), 0))
	case dgo.Float:
		s, f := math.Modf(arg.GoFloat())
		tv = Time(time.Unix(int64(s), int64(f*1000000000.0)))
	case dgo.String:
		tv = TimeFromString(arg.GoString())
	default:
		panic(illegalArgument(`time`, `time|string`, []interface{}{arg}, 0))
	}
	if !t.Instance(tv) {
		panic(IllegalAssignment(t, tv))
	}
	return tv
}

// Time returns the given timestamp as a dgo.Time
func Time(ts time.Time) dgo.Time {
	return (*timeVal)(&ts)
}

// TimeFromString returns the given time string as a dgo.Time. The string must conform to
// the time.RFC3339 or time.RFC3339Nano format. The goFunc will panic if the given string
// cannot be parsed.
func TimeFromString(s string) dgo.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}
	return (*timeVal)(&ts)
}

func (v *timeVal) Equals(other interface{}) bool {
	switch ov := other.(type) {
	case *timeVal:
		return (*time.Time)(v).Equal(*(*time.Time)(ov))
	case time.Time:
		return (*time.Time)(v).Equal(ov)
	case *time.Time:
		return (*time.Time)(v).Equal(*ov)
	}
	return false
}

func (v *timeVal) SecondsWithFraction() float64 {
	t := (*time.Time)(v)
	y := t.Year()
	// Timestamps that represent a date before the year 1678 or after 2262 can
	// be represented as nanoseconds in an int64.
	if 1678 < y && y < 2262 {
		return float64(t.UnixNano()) / 1000000000.0
	}
	// Fall back to microsecond precision
	us := t.Unix()*1000000 + int64(t.Nanosecond())/1000
	return float64(us) / 1000000.0
}

func (v *timeVal) GoTime() time.Time {
	return *(*time.Time)(v)
}

func (v *timeVal) HashCode() int {
	return int((*time.Time)(v).UnixNano())
}

func (v *timeVal) ReflectTo(value reflect.Value) {
	rv := reflect.ValueOf((*time.Time)(v))
	k := value.Kind()
	if !(k == reflect.Ptr || k == reflect.Interface) {
		rv = rv.Elem()
	}
	value.Set(rv)
}

func (v *timeVal) String() string {
	return (*time.Time)(v).Format(time.RFC3339Nano)
}

func (v *timeVal) Type() dgo.Type {
	return (*exactTimeType)(v)
}
