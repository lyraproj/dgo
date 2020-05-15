package internal

import (
	"reflect"
	"sort"

	"github.com/lyraproj/dgo/dgo"
	"github.com/tada/catch"
)

type (
	structVal struct {
		rs         reflect.Value
		fieldCount uint32
		frozen     bool
	}
)

func (v *structVal) AppendTo(w dgo.Indenter) {
	appendMapTo(v, w)
}

func (v *structVal) All(predicate dgo.EntryPredicate) bool {
	rv := v.rs
	return findField(v.frozen, rv.Type(), rv, func(e dgo.MapEntry) bool { return !predicate(e) }) == nil
}

func findField(frozen bool, rt reflect.Type, rv reflect.Value, predicate dgo.EntryPredicate) dgo.MapEntry {
	for i, n := 0, rt.NumField(); i < n; i++ {
		f := rt.Field(i)
		v := rv.Field(i)
		if f.Anonymous {
			ft := f.Type
			if ft.Kind() == reflect.Struct {
				if me := findField(frozen, ft, v, predicate); me != nil {
					return me
				}
				continue
			}
		}
		me := &mapEntry{&hstring{string: f.Name}, ValueFromReflected(v, frozen)}
		if predicate(me) {
			return me
		}
	}
	return nil
}

func (v *structVal) AllKeys(predicate dgo.Predicate) bool {
	return allKeys(v.rs.Type(), predicate)
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
		if !predicate(&hstring{string: f.Name}) {
			return false
		}
	}
	return true
}

func (v *structVal) AllValues(predicate dgo.Predicate) bool {
	rv := v.rs
	return allValues(v.frozen, rv.Type(), rv, predicate)
}

func allValues(frozen bool, rt reflect.Type, rv reflect.Value, predicate dgo.Predicate) bool {
	for i, n := 0, rt.NumField(); i < n; i++ {
		f := rt.Field(i)
		v := rv.Field(i)
		if f.Anonymous {
			ft := f.Type
			if ft.Kind() == reflect.Struct {
				if !allValues(frozen, ft, v, predicate) {
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
	rv := v.rs
	return findField(v.frozen, rv.Type(), rv, predicate) != nil
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
	a := MutableValues([]interface{}{value})
	v.Put(key, a)
	return a
}

func (v *structVal) ContainsKey(key interface{}) bool {
	if s, ok := stringKey(key); ok {
		return v.rs.FieldByName(s).IsValid()
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
	rv := v.rs
	findField(v.frozen, rv.Type(), rv, func(entry dgo.MapEntry) bool { actor(entry); return false })
}

func (v *structVal) EachEntry(actor dgo.EntryActor) {
	rv := v.rs
	findField(v.frozen, rv.Type(), rv, func(entry dgo.MapEntry) bool { actor(entry); return false })
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
	return vm.rs.Addr().Interface()
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
	rv := v.rs
	findField(v.frozen, rv.Type(), rv, func(entry dgo.MapEntry) bool {
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
	rs := reflect.New(v.rs.Type()).Elem() // create and dereference pointer to a zero value
	rs.Set(v.rs)                          // copy v.rs to the zero value

	for i, n := 0, rs.NumField(); i < n; i++ {
		ef := rs.Field(i)
		ev := ValueFromReflected(ef, false)
		if f, ok := ev.(dgo.Mutability); ok && !f.Frozen() {
			ReflectTo(f.FrozenCopy(), ef)
		}
	}
	return &structVal{rs: rs, frozen: true}
}

func (v *structVal) ThawedCopy() dgo.Value {
	// Perform a by-value copy of the struct
	rs := reflect.New(v.rs.Type()).Elem() // create and dereference pointer to a zero value
	rs.Set(v.rs)                          // copy v.rs to the zero value

	for i, n := 0, rs.NumField(); i < n; i++ {
		ef := rs.Field(i)
		ev := ValueFromReflected(ef, false)
		if f, ok := ev.(dgo.Mutability); ok {
			ReflectTo(f.ThawedCopy(), ef)
		}
	}
	return &structVal{rs: rs, frozen: false}
}

func (v *structVal) Find(predicate dgo.EntryPredicate) dgo.MapEntry {
	rv := v.rs
	return findField(v.frozen, rv.Type(), rv, predicate)
}

func stringKey(key interface{}) (string, bool) {
	if hs, ok := key.(*hstring); ok {
		return hs.string, true
	}
	if s, ok := key.(string); ok {
		return s, true
	}
	return ``, false
}

func (v *structVal) Get(key interface{}) dgo.Value {
	if s, ok := stringKey(key); ok {
		rv := v.rs
		fv := rv.FieldByName(s)
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
	if v.fieldCount == 0 {
		v.fieldCount = fieldCount(v.rs.Type())
	}
	return int(v.fieldCount)
}

func fieldCount(rt reflect.Type) uint32 {
	t := uint32(0)
	for i, n := 0, rt.NumField(); i < n; i++ {
		f := rt.Field(i)
		if f.Anonymous {
			ft := f.Type
			if ft.Kind() == reflect.Struct {
				t += fieldCount(ft)
				continue
			}
		}
		t++
	}
	return t
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
	if s, ok := stringKey(key); ok {
		rv := v.rs
		fv := rv.FieldByName(s)
		if fv.IsValid() {
			old := ValueFromReflected(fv, false)
			ReflectTo(Value(value), fv)
			return old
		}
	}
	panic(catch.Error(`%s has no field named '%s'`, v.rs.Type(), key))
}

func (v *structVal) PutAll(associations dgo.Map) {
	associations.EachEntry(func(e dgo.MapEntry) { v.Put(e.Key(), e.Value()) })
}

func (v *structVal) ReflectTo(value reflect.Value) {
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
	nv := &structVal{rs: reflect.New(v.rs.Type()).Elem()}
	nv.PutAll(m)
	return nv
}

func (v *structVal) ReflectType() reflect.Type {
	return v.rs.Type()
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
