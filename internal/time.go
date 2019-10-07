package internal

import (
	"reflect"
	"time"

	"github.com/lyraproj/dgo/dgo"
	"gopkg.in/yaml.v3"
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

// Time returns the given timestamp as a dgo.Time
func Time(ts time.Time) dgo.Time {
	return (*timeVal)(&ts)
}

// TimeFromString returns the given time string as a dgo.Time. The string must conform to
// the time.RFC3339 or time.RFC3339Nano format. The function will panic if the given string
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

func (v *timeVal) GoTime() *time.Time {
	return (*time.Time)(v)
}

func (v *timeVal) HashCode() int {
	return int((*time.Time)(v).UnixNano())
}

func (v *timeVal) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   `!!timestamp`,
		Value: (*time.Time)(v).Format(time.RFC3339Nano),
		Style: yaml.TaggedStyle}, nil
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
