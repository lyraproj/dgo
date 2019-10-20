package internal

import (
	"fmt"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	function reflect.Value
)

func (f *function) Equals(other interface{}) bool {
	if b, ok := toReflected(other); ok {
		return reflect.Func == b.Kind() && (*reflect.Value)(f).Pointer() == b.Pointer()
	}
	return false
}

func (f *function) Type() dgo.Type {
	return &nativeType{(*reflect.Value)(f).Type()}
}

func (f *function) HashCode() int {
	return int((*reflect.Value)(f).Pointer())
}

func (f *function) Call(args ...interface{}) []dgo.Value {
	convertReturn := func(rr []reflect.Value) []dgo.Value {
		vr := make([]dgo.Value, len(rr))
		for i := range rr {
			vr[i] = ValueFromReflected(rr[i])
		}
		return vr
	}

	mx := len(args)
	m := (*reflect.Value)(f)
	t := m.Type()
	if t.IsVariadic() {
		nv := t.NumIn() - 1 // number of non variadic
		if mx < nv {
			panic(fmt.Errorf(`illegal number of arguments. Expected at least %d, got %d`, nv, mx))
		}
		rr := make([]reflect.Value, nv+1)
		for i := 0; i < nv; i++ {
			rr[i] = reflect.New(t.In(i)).Elem()
			ReflectTo(Value(args[i]), rr[i])
		}

		// Create the variadic slice
		vt := t.In(nv)
		vz := mx - nv
		vs := reflect.MakeSlice(vt, vz, vz)
		rr[nv] = vs

		for i := 0; i < vz; i++ {
			ReflectTo(Value(args[i+nv]), vs.Index(i))
		}
		return convertReturn(m.CallSlice(rr))
	}

	if mx != t.NumIn() {
		panic(fmt.Errorf(`illegal number of arguments. Expected %d, got %d`, t.NumIn(), mx))
	}

	rr := make([]reflect.Value, mx)
	for i := range args {
		rr[i] = reflect.New(t.In(i)).Elem()
		ReflectTo(Value(args[i]), rr[i])
	}
	return convertReturn(m.Call(rr))
}

func (f *function) String() string {
	return (*reflect.Value)(f).String()
}
