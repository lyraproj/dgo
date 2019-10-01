package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
	"gopkg.in/yaml.v3"
)

type (
	structMap struct {
		rs     reflect.Value
		frozen bool
	}
)

func (v *structMap) String() string {
	return ToStringERP(v)
}

func (v *structMap) Type() dgo.Type {
	return &exactMapType{v}
}

func (v *structMap) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *structMap) deepEqual(seen []dgo.Value, other deepEqual) bool {
	return mapEqual(seen, v, other)
}

func (v *structMap) HashCode() int {
	return deepHashCode(nil, v)
}

func (v *structMap) deepHashCode(seen []dgo.Value) int {
	hs := make([]int, v.Len())
	i := 0
	v.EachEntry(func(e dgo.MapEntry) {
		hs[i] = deepHashCode(seen, e)
		i++
	})
	sort.Ints(hs)
	h := 1
	for i = range hs {
		h = h*31 + hs[i]
	}
	return h
}

func (v *structMap) Freeze() {
	// Perform a shallow copy of the struct
	if !v.frozen {
		v.frozen = true
		v.rs = reflect.ValueOf(v.rs.Interface())
		v.EachValue(func(e dgo.Value) {
			if f, ok := e.(dgo.Freezable); ok {
				f.Freeze()
			}
		})
	}
}

func (v *structMap) Frozen() bool {
	return v.frozen
}

func (v *structMap) FrozenCopy() dgo.Value {
	if v.frozen {
		return v
	}

	// Perform a by-value copy of the struct
	rs := reflect.New(v.rs.Type()).Elem() // create and dereference pointer to a zero value
	rs.Set(v.rs)                          // copy v.rs to the zero value

	for i, n := 0, rs.NumField(); i < n; i++ {
		ef := rs.Field(i)
		ev := ValueFromReflected(ef)
		if f, ok := ev.(dgo.Freezable); ok && !f.Frozen() {
			ReflectTo(f.FrozenCopy(), ef)
		}
	}
	return &structMap{rs: rs, frozen: true}
}

func (v *structMap) AppendTo(w util.Indenter) {
	appendMapTo(v, w)
}

func (v *structMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.rs.Addr().Interface())
}

func (v *structMap) MarshalYAML() (interface{}, error) {
	// A bit wasteful but this is currently the only way to create a yaml.Node
	// from a struct
	n := &yaml.Node{}
	b, err := yaml.Marshal(v.rs.Addr().Interface())
	if err == nil {
		if err = yaml.Unmarshal(b, n); err == nil {
			// n is the document node at this point.
			n = n.Content[0]
		}
	}
	return n, err
}

func (v *structMap) UnmarshalJSON(data []byte) error {
	if v.frozen {
		panic(frozenMap(`UnmarshalJSON`))
	}
	return json.Unmarshal(data, v.rs.Addr().Interface())
}

func (v *structMap) UnmarshalYAML(value *yaml.Node) error {
	if v.frozen {
		panic(frozenMap(`UnmarshalYAML`))
	}
	return value.Decode(v.rs.Addr().Interface())
}

func (v *structMap) All(predicate dgo.EntryPredicate) bool {
	rv := v.rs
	rt := rv.Type()
	for i, n := 0, rt.NumField(); i < n; i++ {
		if !predicate(&mapEntry{&hstring{s: rt.Field(i).Name}, ValueFromReflected(rv.Field(i))}) {
			return false
		}
	}
	return true
}

func (v *structMap) AllKeys(predicate dgo.Predicate) bool {
	rv := v.rs
	rt := rv.Type()
	for i, n := 0, rt.NumField(); i < n; i++ {
		if !predicate(&hstring{s: rt.Field(i).Name}) {
			return false
		}
	}
	return true
}

func (v *structMap) AllValues(predicate dgo.Predicate) bool {
	rv := v.rs
	for i, n := 0, rv.NumField(); i < n; i++ {
		if !predicate(ValueFromReflected(rv.Field(i))) {
			return false
		}
	}
	return true
}

func (v *structMap) Any(predicate dgo.EntryPredicate) bool {
	return !v.All(func(entry dgo.MapEntry) bool { return !predicate(entry) })
}

func (v *structMap) AnyKey(predicate dgo.Predicate) bool {
	return !v.AllKeys(func(entry dgo.Value) bool { return !predicate(entry) })
}

func (v *structMap) AnyValue(predicate dgo.Predicate) bool {
	return !v.AllValues(func(entry dgo.Value) bool { return !predicate(entry) })
}

func (v *structMap) Copy(frozen bool) dgo.Map {
	if frozen && v.frozen {
		return v
	}
	if frozen {
		return v.FrozenCopy().(dgo.Map)
	}
	// Thaw: Perform a by-value copy of the frozen struct
	rs := reflect.New(v.rs.Type()).Elem() // create and dereference pointer to a zero value
	rs.Set(v.rs)                          // copy v.rs to the zero value
	return &structMap{rs: rs, frozen: false}
}

func (v *structMap) Each(actor dgo.Actor) {
	v.All(func(entry dgo.MapEntry) bool { actor(entry); return true })
}

func (v *structMap) EachEntry(actor dgo.EntryActor) {
	v.All(func(entry dgo.MapEntry) bool { actor(entry); return true })
}

func (v *structMap) EachKey(actor dgo.Actor) {
	v.AllKeys(func(entry dgo.Value) bool { actor(entry); return true })
}

func (v *structMap) EachValue(actor dgo.Actor) {
	v.AllValues(func(entry dgo.Value) bool { actor(entry); return true })
}

func stringKey(key interface{}) (string, bool) {
	if hs, ok := key.(*hstring); ok {
		return hs.s, true
	}
	if s, ok := key.(string); ok {
		return s, true
	}
	return ``, false
}

func (v *structMap) Find(predicate dgo.EntryPredicate) dgo.MapEntry {
	rv := v.rs
	rt := rv.Type()
	for i, n := 0, rt.NumField(); i < n; i++ {
		e := &mapEntry{&hstring{s: rt.Field(i).Name}, ValueFromReflected(rv.Field(i))}
		if predicate(e) {
			return e
		}
	}
	return nil
}

func (v *structMap) Get(key interface{}) dgo.Value {
	if s, ok := stringKey(key); ok {
		rv := v.rs
		fv := rv.FieldByName(s)
		if fv.IsValid() {
			return ValueFromReflected(fv)
		}
	}
	return nil
}

func (v *structMap) Keys() dgo.Array {
	return arrayFromIterator(v.Len(), v.EachKey)
}

func (v *structMap) Len() int {
	return v.rs.NumField()
}

func (v *structMap) Map(mapper dgo.EntryMapper) dgo.Map {
	c := v.toHashMap()
	for e := c.first; e != nil; e = e.next {
		e.value = Value(mapper(e))
	}
	c.frozen = v.frozen
	return c
}

func (v *structMap) Merge(associations dgo.Map) dgo.Map {
	if associations.Len() == 0 || v == associations {
		return v
	}
	c := v.toHashMap()
	c.PutAll(associations)
	c.frozen = v.frozen && associations.Frozen()
	return c
}

func (v *structMap) Put(key, value interface{}) dgo.Value {
	if v.frozen {
		panic(frozenMap(`Put`))
	}
	if s, ok := stringKey(key); ok {
		rv := v.rs
		fv := rv.FieldByName(s)
		if fv.IsValid() {
			old := ValueFromReflected(fv)
			ReflectTo(Value(value), fv)
			return old
		}
	}
	panic(fmt.Errorf(`%s has no field named '%s'`, v.rs.Type(), key))
}

func (v *structMap) PutAll(associations dgo.Map) {
	associations.EachEntry(func(e dgo.MapEntry) { v.Put(e.Key(), e.Value()) })
}

func (v *structMap) ReflectTo(value reflect.Value) {
	if value.Kind() == reflect.Ptr {
		if v.frozen {
			// Don't expose pointer to frozen struct
			rs := reflect.New(v.rs.Type()) // create pointer to a zero value
			rs.Elem().Set(v.rs)            // copy v.rs to the zero value
			value.Set(rs)
		} else {
			value.Set(v.rs.Addr())
		}
	} else {
		// copy struct by value
		value.Set(v.rs)
	}
}

func (v *structMap) Remove(key interface{}) dgo.Value {
	panic(errors.New(`struct fields cannot be removed`))
}

func (v *structMap) RemoveAll(keys dgo.Array) {
	panic(errors.New(`struct fields cannot be removed`))
}

func (v *structMap) SetType(t interface{}) {
	panic(errors.New(`struct type is read only`))
}

func (v *structMap) StringKeys() bool {
	return true
}

func (v *structMap) Values() dgo.Array {
	return arrayFromIterator(v.Len(), v.EachValue)
}

func (v *structMap) With(key, value interface{}) dgo.Map {
	c := v.toHashMap()
	c.Put(key, value)
	c.frozen = v.frozen
	return c
}

func (v *structMap) Without(key interface{}) dgo.Map {
	if v.Get(key) == nil {
		return v
	}
	c := v.toHashMap()
	c.Remove(key)
	c.frozen = v.frozen
	return c
}

func (v *structMap) WithoutAll(keys dgo.Array) dgo.Map {
	c := v.toHashMap()
	c.RemoveAll(keys)
	c.frozen = v.frozen
	return c
}

func (v *structMap) toHashMap() *hashMap {
	c := MapWithCapacity(int(float64(v.Len())/loadFactor), nil)
	v.EachEntry(func(entry dgo.MapEntry) {
		c.Put(entry.Key(), entry.Value())
	})
	return c.(*hashMap)
}
