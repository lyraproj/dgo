package internal

import (
	"reflect"
	"sort"

	"github.com/lyraproj/dgo/dgo"
	"github.com/tada/catch"
)

type (
	structVal struct {
		native
		frozen bool
	}
)

// Struct creates a dgo.Struct which implements the dgo.Map interface backed by the given go struct. It panics if
// rm's kind is not reflect.Struct.
func Struct(rs reflect.Value, frozen bool) dgo.Struct {
	return &structVal{native: native{_rv: rs}, frozen: frozen}
}

func (v *structVal) AppendTo(w dgo.Indenter) {
	appendMapTo(v, w)
}

func (v *structVal) All(predicate dgo.EntryPredicate) bool {
	rv := &v._rv
	return v.findField(rv.Type(), rv, func(e dgo.MapEntry) bool { return !predicate(e) }) == nil
}

func (v *structVal) findField(rt reflect.Type, rv *reflect.Value, predicate dgo.EntryPredicate) dgo.MapEntry {
	for i, n := 0, rt.NumField(); i < n; i++ {
		f := rt.Field(i)
		fv := rv.Field(i)
		if f.Anonymous {
			ft := f.Type
			if ft.Kind() == reflect.Struct {
				if me := v.findField(ft, &fv, predicate); me != nil {
					return me
				}
				continue
			}
		}
		me := &mapEntry{String(f.Name), ValueFromReflected(fv, v.frozen)}
		if predicate(me) {
			return me
		}
	}
	return nil
}

func (v *structVal) AllKeys(predicate dgo.Predicate) bool {
	return allKeys(v.ReflectType(), predicate)
}

func allKeys(rt reflect.Type, predicate dgo.Predicate) bool {
	for i, n := 0, rt.NumField(); i < n; i++ {
		f := rt.Field(i)
		if f.Anonymous {
			ft := f.Type
			if ft.Kind() == reflect.Struct {
				if !allKeys(ft, predicate) {
					return false
				}
				continue
			}
		}
		if !predicate(String(f.Name)) {
			return false
		}
	}
	return true
}

func (v *structVal) AllValues(predicate dgo.Predicate) bool {
	rv := &v._rv
	return allValues(v.frozen, rv.Type(), rv, predicate)
}

func allValues(frozen bool, rt reflect.Type, rv *reflect.Value, predicate dgo.Predicate) bool {
	for i, n := 0, rt.NumField(); i < n; i++ {
		f := rt.Field(i)
		v := rv.Field(i)
		if f.Anonymous {
			ft := f.Type
			if ft.Kind() == reflect.Struct {
				if !allValues(frozen, ft, &v, predicate) {
					return false
				}
				continue
			}
		}
		if !predicate(ValueFromReflected(rv.Field(i), frozen)) {
			return false
		}
	}
	return true
}

func (v *structVal) Any(predicate dgo.EntryPredicate) bool {
	rv := &v._rv
	return v.findField(rv.Type(), rv, predicate) != nil
}

func (v *structVal) AnyKey(predicate dgo.Predicate) bool {
	return !v.AllKeys(func(entry dgo.Value) bool { return !predicate(entry) })
}

func (v *structVal) AnyValue(predicate dgo.Predicate) bool {
	return !v.AllValues(func(entry dgo.Value) bool { return !predicate(entry) })
}

func (v *structVal) AppendAt(key, value interface{}) dgo.Array {
	if a, ok := v.Get(key).(dgo.Array); ok {
		a.Add(value)
		return a
	}
	a := MutableValues(value)
	v.Put(key, a)
	return a
}

func stringKey(key interface{}) string {
	s := ``
	switch key := key.(type) {
	case string:
		s = key
	case dgo.String:
		s = key.GoString()
	}
	return s
}

func (v *structVal) ContainsKey(key interface{}) bool {
	if s := stringKey(key); s != `` {
		return v.FieldByName(s).IsValid()
	}
	return false
}

func (v *structVal) Copy(frozen bool) dgo.Map {
	if frozen && v.frozen {
		return v
	}
	if frozen {
		return v.FrozenCopy().(dgo.Map)
	}
	return v.ThawedCopy().(dgo.Map)
}

func (v *structVal) Each(actor dgo.Consumer) {
	rv := &v._rv
	v.findField(rv.Type(), rv, func(entry dgo.MapEntry) bool { actor(entry); return false })
}

func (v *structVal) EachEntry(actor dgo.EntryActor) {
	rv := &v._rv
	v.findField(rv.Type(), rv, func(entry dgo.MapEntry) bool { actor(entry); return false })
}

func (v *structVal) EachKey(actor dgo.Consumer) {
	v.AllKeys(func(entry dgo.Value) bool { actor(entry); return true })
}

func (v *structVal) EachValue(actor dgo.Consumer) {
	v.AllValues(func(entry dgo.Value) bool { actor(entry); return true })
}

func (v *structVal) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *structVal) GoStruct() interface{} {
	vm := v
	if vm.frozen {
		vm = vm.ThawedCopy().(*structVal)
	}
	return vm.Addr().Interface()
}

func (v *structVal) deepEqual(seen []dgo.Value, other deepEqual) bool {
	return mapEqual(seen, v, other)
}

func (v *structVal) HashCode() dgo.Hash {
	return deepHashCode(nil, v)
}

func (v *structVal) deepHashCode(seen []dgo.Value) dgo.Hash {
	hs := make([]int, v.Len())
	i := 0
	rv := &v._rv
	v.findField(rv.Type(), rv, func(entry dgo.MapEntry) bool {
		hs[i] = int(deepHashCode(seen, entry))
		i++
		return false
	})
	sort.Ints(hs)
	h := dgo.Hash(1)
	for i = range hs {
		h = h*31 + dgo.Hash(hs[i])
	}
	return h
}

func (v *structVal) Frozen() bool {
	return v.frozen
}

func (v *structVal) FrozenCopy() dgo.Value {
	if v.frozen {
		return v
	}

	// Perform a by-value copy of the struct
	rs := reflect.New(v.ReflectType()).Elem() // create and dereference pointer to a zero value
	rs.Set(v._rv)                             // copy v.rs to the zero value

	for i, n := 0, rs.NumField(); i < n; i++ {
		ef := rs.Field(i)
		ev := ValueFromReflected(ef, false)
		if f, ok := ev.(dgo.Mutability); ok && !f.Frozen() {
			ReflectTo(f.FrozenCopy(), ef)
		}
	}
	return Struct(rs, true)
}

func (v *structVal) ThawedCopy() dgo.Value {
	// Perform a by-value copy of the struct
	rs := reflect.New(v.ReflectType()).Elem() // create and dereference pointer to a zero value
	rs.Set(v._rv)                             // copy v.rs to the zero value

	for i, n := 0, rs.NumField(); i < n; i++ {
		ef := rs.Field(i)
		ev := ValueFromReflected(ef, false)
		if f, ok := ev.(dgo.Mutability); ok {
			ReflectTo(f.ThawedCopy(), ef)
		}
	}
	return Struct(rs, false)
}

func (v *structVal) Find(predicate dgo.EntryPredicate) dgo.MapEntry {
	rv := &v._rv
	return v.findField(rv.Type(), rv, predicate)
}

func (v *structVal) Get(key interface{}) dgo.Value {
	if s := stringKey(key); s != `` {
		fv := v.FieldByName(s)
		if fv.IsValid() {
			return ValueFromReflected(fv, v.frozen)
		}
	}
	return nil
}

func (v *structVal) Keys() dgo.Array {
	return arrayFromIterator(v.Len(), v.EachKey)
}

func (v *structVal) Len() int {
	return getTags(v.ReflectType()).fieldCount()
}

func (v *structVal) Map(mapper dgo.EntryMapper) dgo.Map {
	c := v.toHashMap()
	for e := c.first; e != nil; e = e.next {
		e.value = Value(mapper(e))
	}
	c.frozen = v.frozen
	return c
}

func (v *structVal) Merge(associations dgo.Map) dgo.Map {
	if associations.Len() == 0 || v == associations {
		return v
	}
	c := v.toHashMap()
	c.PutAll(associations)
	c.frozen = v.frozen && associations.Frozen()
	return c
}

func (v *structVal) Put(key, value interface{}) dgo.Value {
	if v.frozen {
		panic(frozenMap(`Put`))
	}
	n, tp := getTags(v.ReflectType()).fieldType(key)
	if n != `` {
		fv := v.FieldByName(n)
		if fv.IsValid() {
			v := Value(value)
			if !(tp == nil || tp.Instance(v)) {
				panic(IllegalAssignment(tp, v))
			}
			old := ValueFromReflected(fv, false)
			ReflectTo(v, fv)
			return old
		}
	}
	panic(catch.Error(`%s has no field named '%v'`, v.ReflectType(), key))
}

func (v *structVal) PutAll(associations dgo.Map) {
	associations.EachEntry(func(e dgo.MapEntry) { v.Put(e.Key(), e.Value()) })
}

func (v *structVal) ReflectTo(value reflect.Value) {
	if value.Kind() == reflect.Ptr {
		if v.frozen {
			// Don't expose pointer to frozen struct
			rs := reflect.New(v.ReflectType()) // create pointer to a zero value
			rs.Elem().Set(v._rv)               // copy v.rs to the zero value
			value.Set(rs)
		} else {
			value.Set(v.Addr())
		}
	} else {
		// copy struct by value
		value.Set(v._rv)
	}
}

func (v *structVal) Remove(_ interface{}) dgo.Value {
	panic(catch.Error(`struct fields cannot be removed`))
}

func (v *structVal) RemoveAll(_ dgo.Array) {
	panic(catch.Error(`struct fields cannot be removed`))
}

func (v *structVal) String() string {
	return TypeString(v)
}

func (v *structVal) StringKeys() bool {
	return true
}

func (v *structVal) Type() dgo.Type {
	return v
}

func (v *structVal) Values() dgo.Array {
	return arrayFromIterator(v.Len(), v.EachValue)
}

func (v *structVal) With(key, value interface{}) dgo.Map {
	c := v.toHashMap()
	c.Put(key, value)
	c.frozen = v.frozen
	return c
}

func (v *structVal) Without(key interface{}) dgo.Map {
	if v.Get(key) == nil {
		return v
	}
	c := v.toHashMap()
	c.Remove(key)
	c.frozen = v.frozen
	return c
}

func (v *structVal) WithoutAll(keys dgo.Array) dgo.Map {
	c := v.toHashMap()
	c.RemoveAll(keys)
	c.frozen = v.frozen
	return c
}

func (v *structVal) toHashMap() *hashMap {
	c := MapWithCapacity(v.Len())
	v.EachEntry(func(entry dgo.MapEntry) {
		c.Put(entry.Key(), entry.Value())
	})
	return c.(*hashMap)
}

func (v *structVal) Additional() bool {
	return false
}

func (v *structVal) Assignable(other dgo.Type) bool {
	return v.Equals(other)
}

func (v *structVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v *structVal) EachEntryType(actor func(dgo.StructMapEntry)) {
	eachEntryType(v, actor)
}

func (v *structVal) Generic() dgo.Type {
	return genericMapType(v)
}

func (v *structVal) GetEntryType(key interface{}) dgo.StructMapEntry {
	return entryType(v, key)
}

func (v *structVal) KeyType() dgo.Type {
	return keyType(v)
}

func (v *structVal) Max() int {
	return v.Len()
}

func (v *structVal) Min() int {
	return v.Len()
}

func (v *structVal) New(arg dgo.Value) dgo.Value {
	m := newMap(v, arg)
	nv := Struct(reflect.New(v.ReflectType()).Elem(), false)
	nv.PutAll(m)
	return nv
}

func (v *structVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMapExact
}

func (v *structVal) Unbounded() bool {
	return false
}

func (v *structVal) ValueType() dgo.Type {
	return valueType(v)
}

func (v *structVal) Validate(keyLabel func(key dgo.Value) string, value interface{}) []error {
	return validate(v, keyLabel, value)
}

func (v *structVal) ValidateVerbose(value interface{}, out dgo.Indenter) bool {
	return validateVerbose(v, value, out)
}
