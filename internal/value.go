package internal

import (
	"fmt"
	"math"
	"reflect"
	"regexp"

	"github.com/lyraproj/dgo/dgo"
)

// Value returns the dgo.Value representation of its argument
func Value(v interface{}) dgo.Value {
	// This function is kept very small to enable inlining so this
	// if statement should not be baked in to the grand switch
	// in the value function
	if gv, ok := v.(dgo.Value); ok {
		return gv
	}
	return value(v)
}

func value(v interface{}) dgo.Value {
	if v == nil {
		return Nil
	}
	switch v := v.(type) {
	case bool:
		return Boolean(v)
	case string:
		return String(v)
	case int:
		return Integer(int64(v))
	case int8:
		return Integer(int64(v))
	case int16:
		return Integer(int64(v))
	case int32:
		return Integer(int64(v))
	case uint8:
		return Integer(int64(v))
	case uint16:
		return Integer(int64(v))
	case uint32:
		return Integer(int64(v))
	case int64:
		return Integer(v)
	case uint:
		if v == math.MaxUint64 {
			panic(fmt.Errorf(`value %d overflows int64`, v))
		}
		return Integer(int64(v))
	case uint64:
		if v == math.MaxUint64 {
			panic(fmt.Errorf(`value %d overflows int64`, v))
		}
		return Integer(int64(v))
	case float32:
		return Float(float64(v))
	case float64:
		return Float(v)
	case []byte:
		return Binary(v)
	case []string:
		return Strings(v)
	case []int:
		return Integers(v)
	case *regexp.Regexp:
		return (*Regexp)(v)
	case error:
		return &errw{v}
	case reflect.Value:
		return ValueFromReflected(v)
	default:
		return ValueFromReflected(reflect.ValueOf(v))
	}
}

var dgoValueType = reflect.TypeOf((*dgo.Value)(nil)).Elem()

func ValueFromReflected(vr reflect.Value) (pv dgo.Value) {
	// Invalid shouldn't happen, but needs a check
	if !vr.IsValid() {
		return Nil
	}

	switch vr.Kind() {
	case reflect.Slice:
		if vr.IsNil() {
			pv = Nil
		} else {
			pv = ArrayFromReflected(vr, true)
		}
	case reflect.Map:
		if vr.IsNil() {
			pv = Nil
		} else {
			pv = MapFromReflected(vr, true)
		}
	case reflect.String:
		pv = String(vr.String())
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8,
		reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		if i, ok := ToInt(vr.Interface()); ok {
			pv = Integer(i)
		}
	case reflect.Bool:
		pv = Boolean(vr.Bool())
	case reflect.Float64, reflect.Float32:
		if f, ok := ToFloat(vr.Interface()); ok {
			pv = Float(f)
		}
	case reflect.Interface:
		if vr.IsNil() {
			pv = Nil
			break
		}
		pv = Value(vr.Interface())
	case reflect.Ptr:
		if vr.IsNil() {
			pv = Nil
		}
	}
	if pv != nil {
		return
	}
	if vr.Type().AssignableTo(dgoValueType) {
		return vr.Interface().(dgo.Value)
	}
	if vf, ok := wellKnown[vr.Type()]; ok {
		pv = vf(vr)
	} else {
		// Value as unsafe. Immutability is not guaranteed
		pv = native(vr)
	}
	return
}

// Add well known types like regexp, time, etc. here
var wellKnown map[reflect.Type]func(reflect.Value) dgo.Value
var wellKnownTypes map[reflect.Type]dgo.Type

func init() {
	wellKnown = map[reflect.Type]func(reflect.Value) dgo.Value{
		reflect.TypeOf(&regexp.Regexp{}): func(v reflect.Value) dgo.Value { return (*Regexp)(v.Interface().(*regexp.Regexp)) },
	}
	wellKnownTypes = map[reflect.Type]dgo.Type{
		reflect.TypeOf(&regexp.Regexp{}): DefaultRegexpType,
	}
}
