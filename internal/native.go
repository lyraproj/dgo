package internal

import (
	"fmt"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	native reflect.Value

	nativeType struct {
		rt reflect.Type
	}
)

var DefaultNativeType = &nativeType{}

func Native(rv reflect.Value) dgo.Native {
	return native(rv)
}

func (t *nativeType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*nativeType); ok {
		if t.rt == nil {
			return true
		}
		if ot.rt == nil {
			return false
		}
		return ot.rt.AssignableTo(t.rt)
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *nativeType) Equals(other interface{}) bool {
	if ot, ok := other.(*nativeType); ok {
		return t.rt == ot.rt
	}
	return false
}

func (t *nativeType) HashCode() int {
	return stringHash(t.rt.Name())*31 + int(dgo.IdNative)
}

func (t *nativeType) Type() dgo.Type {
	return &metaType{t}
}

func (t *nativeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.IdNative
}

func (t *nativeType) Instance(value interface{}) bool {
	if ov, ok := toReflected(value); ok {
		if t.rt == nil {
			return true
		}
		return ov.Type().AssignableTo(t.rt)
	}
	return false
}

func (t *nativeType) String() string {
	return TypeString(t)
}

func (t *nativeType) GoType() reflect.Type {
	return t.rt
}

func (v native) Equals(other interface{}) bool {
	if b, ok := toReflected(other); ok {
		a := reflect.Value(v)
		k := a.Kind()
		if k != b.Kind() {
			return false
		}
		if k == reflect.Func {
			return a.Pointer() == b.Pointer()
		}
		return reflect.DeepEqual(a.Interface(), b.Interface())
	}
	return false
}

func (v native) Freeze() {
	if !v.Frozen() {
		panic(fmt.Errorf(`native value cannot be frozen`))
	}
}

func (v native) Frozen() bool {
	return reflect.Value(v).Kind() == reflect.Func
}

func (v native) FrozenCopy() dgo.Value {
	if v.Frozen() {
		return v
	}
	panic(fmt.Errorf(`native value cannot be frozen`))
}

func (v native) HashCode() int {
	rv := reflect.Value(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Uintptr:
		p := rv.Pointer()
		return int(p ^ (p >> 32))
	case reflect.Struct:
		n := rv.NumField()
		h := 1
		for i := 0; i < n; i++ {
			h = h*31 + ValueFromReflected(rv.Field(i)).HashCode()
		}
		return h
	}
	return 1234
}

func (v native) String() string {
	rv := reflect.Value(v)
	if rv.CanInterface() {
		if s, ok := rv.Interface().(fmt.Stringer); ok {
			return s.String()
		}
	}
	return rv.String()
}

func (v native) Type() dgo.Type {
	return &nativeType{reflect.Value(v).Type()}
}

func (v native) GoValue() interface{} {
	return reflect.Value(v).Interface()
}

func toReflected(value interface{}) (reflect.Value, bool) {
	if ov, ok := value.(native); ok {
		return reflect.Value(ov), true
	}
	if _, ok := value.(dgo.Value); ok {
		return reflect.Value{}, false
	}
	return reflect.ValueOf(value), true
}
