package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
	"gopkg.in/yaml.v3"
)

type (
	array struct {
		slice  []dgo.Value
		typ    dgo.ArrayType
		frozen bool
	}

	// defaultArrayType is the unconstrained array type
	defaultArrayType int

	// sizedArrayType represents array with element type constraint and a size constraint
	sizedArrayType struct {
		elementType dgo.Type
		min         int
		max         int
	}

	// tupleType represents an array with an exact number of ordered element types.
	tupleType array

	// exactArrayType only matches the array that it represents
	exactArrayType array
)

// DefaultArrayType is the unconstrained Array type
const DefaultArrayType = defaultArrayType(0)

func arrayTypeOne(args []interface{}) dgo.ArrayType {
	switch a0 := Value(args[0]).(type) {
	case dgo.Type:
		return newArrayType(a0, 0, math.MaxInt64)
	case dgo.Integer:
		return newArrayType(nil, int(a0.GoInt()), math.MaxInt64)
	default:
		panic(illegalArgument(`Array`, `Type or Integer`, args, 0))
	}
}

func arrayTypeTwo(args []interface{}) dgo.ArrayType {
	a1, ok := Value(args[1]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Array`, `Integer`, args, 1))
	}
	switch a0 := Value(args[0]).(type) {
	case dgo.Type:
		return newArrayType(a0, int(a1.GoInt()), math.MaxInt64)
	case dgo.Integer:
		return newArrayType(nil, int(a0.GoInt()), int(a1.GoInt()))
	default:
		panic(illegalArgument(`Array`, `Type or Integer`, args, 0))
	}
}

func arrayTypeThree(args []interface{}) dgo.ArrayType {
	a0, ok := Value(args[0]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Array`, `Type`, args, 0))
	}
	a1, ok := Value(args[1]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Array`, `Integer`, args, 1))
	}
	a2, ok := Value(args[2]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`ArrayType`, `Integer`, args, 2))
	}
	return newArrayType(a0, int(a1.GoInt()), int(a2.GoInt()))
}

// ArrayType returns a type that represents an Array value
func ArrayType(args ...interface{}) dgo.ArrayType {
	switch len(args) {
	case 0:
		return DefaultArrayType
	case 1:
		return arrayTypeOne(args)
	case 2:
		return arrayTypeTwo(args)
	case 3:
		return arrayTypeThree(args)
	default:
		panic(fmt.Errorf(`illegal number of arguments for Array. Expected 0 - 3, got %d`, len(args)))
	}
}

func newArrayType(elementType dgo.Type, min, max int) dgo.ArrayType {
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	if elementType == nil {
		elementType = DefaultAnyType
	}
	if min == 0 && max == math.MaxInt64 && elementType == DefaultAnyType {
		// Unbounded
		return DefaultArrayType
	}
	return &sizedArrayType{elementType: elementType, min: min, max: max}
}

func (t defaultArrayType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case defaultArrayType, *tupleType, *exactArrayType, *sizedArrayType:
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t defaultArrayType) ElementType() dgo.Type {
	return DefaultAnyType
}

func (t defaultArrayType) Equals(other interface{}) bool {
	return t == other
}

func (t defaultArrayType) HashCode() int {
	return int(dgo.TiArray)
}

func (t defaultArrayType) Instance(value interface{}) bool {
	_, ok := value.(dgo.Array)
	return ok
}

func (t defaultArrayType) Max() int {
	return math.MaxInt64
}

func (t defaultArrayType) Min() int {
	return 0
}

func (t defaultArrayType) String() string {
	return TypeString(t)
}

func (t defaultArrayType) Type() dgo.Type {
	return &metaType{t}
}

func (t defaultArrayType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArray
}

func (t defaultArrayType) Unbounded() bool {
	return true
}

func (t *sizedArrayType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *sizedArrayType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	switch ot := other.(type) {
	case defaultArrayType:
		return false // lacks size
	case *sizedArrayType:
		return t.min <= ot.min && ot.max <= t.max && t.elementType.Assignable(ot.elementType)
	case *tupleType:
		l := len(ot.slice)
		return t.min <= l && l <= t.max && allAssignable(guard, t.elementType, ot.slice)
	case *exactArrayType:
		l := len(ot.slice)
		return t.min <= l && l <= t.max && t.elementType.Assignable(ot.ElementType())
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *sizedArrayType) ElementType() dgo.Type {
	return t.elementType
}

func (t *sizedArrayType) Equals(other interface{}) bool {
	if ot, ok := other.(*sizedArrayType); ok {
		return t.min == ot.min && t.max == ot.max && t.elementType.Equals(ot.elementType)
	}
	return false
}

func (t *sizedArrayType) HashCode() int {
	h := int(dgo.TiArray)
	if t.min > 0 {
		h = h*31 + t.min
	}
	if t.max < math.MaxInt64 {
		h = h*31 + t.max
	}
	if DefaultAnyType != t.elementType {
		h = h*31 + t.elementType.HashCode()
	}
	return h
}

func (t *sizedArrayType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *sizedArrayType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	if ov, ok := value.(*array); ok {
		l := len(ov.slice)
		return t.min <= l && l <= t.max && allInstance(guard, t.elementType, ov.slice)
	}
	return false
}

func (t *sizedArrayType) Max() int {
	return t.max
}

func (t *sizedArrayType) Min() int {
	return t.min
}

func (t *sizedArrayType) Resolve(ap dgo.AliasProvider) {
	t.elementType = ap.Replace(t.elementType)
}

func (t *sizedArrayType) String() string {
	return TypeString(t)
}

func (t *sizedArrayType) Type() dgo.Type {
	return &metaType{t}
}

func (t *sizedArrayType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArrayElementSized
}

func (t *sizedArrayType) Unbounded() bool {
	return t.min == 0 && t.max == math.MaxInt64
}

func (t *exactArrayType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *exactArrayType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	es := t.slice
	switch ot := other.(type) {
	case defaultArrayType:
		return false // lacks size
	case *sizedArrayType:
		l := len(es)
		return ot.min == l && ot.max == l && assignableToAll(guard, ot.elementType, es)
	case *tupleType:
		os := ot.slice
		l := len(es)
		if l != len(os) {
			return false
		}
		for i := range es {
			if !Assignable(guard, es[i].Type(), os[i].(dgo.Type)) {
				return false
			}
		}
		return true
	case *exactArrayType:
		return sliceEquals(nil, es, ot.slice)
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *exactArrayType) ElementType() dgo.Type {
	switch len(t.slice) {
	case 0:
		return DefaultAnyType
	case 1:
		return t.slice[0].Type()
	}
	return (*allOfValueType)(t)
}

func (t *exactArrayType) ElementTypes() dgo.Array {
	es := t.slice
	ts := make([]dgo.Value, len(es))
	for i := range es {
		ts[i] = es[i].Type()
	}
	return &array{slice: ts, frozen: true}
}

func (t *exactArrayType) Equals(other interface{}) bool {
	if ot, ok := other.(*exactArrayType); ok {
		return (*array)(t).Equals((*array)(ot))
	}
	return false
}

func (t *exactArrayType) HashCode() int {
	return (*array)(t).HashCode()*7 + int(dgo.TiArrayExact)
}

func (t *exactArrayType) Instance(value interface{}) bool {
	if ot, ok := value.(*array); ok {
		return (*array)(t).Equals(ot)
	}
	return false
}

func (t *exactArrayType) Max() int {
	return len(t.slice)
}

func (t *exactArrayType) Min() int {
	return len(t.slice)
}

func (t *exactArrayType) String() string {
	return TypeString(t)
}

func (t *exactArrayType) Value() dgo.Value {
	a := (*array)(t)
	return a
}

func (t *exactArrayType) Type() dgo.Type {
	return &metaType{t}
}

func (t *exactArrayType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArrayExact
}

func (t *exactArrayType) Unbounded() bool {
	return false
}

// DefaultTupleType is the unconstrained Tuple type
var DefaultTupleType = &tupleType{}

// TupleType creates a new TupleTupe based on the given types
func TupleType(types []dgo.Type) dgo.TupleType {
	l := len(types)
	if l == 0 {
		return DefaultTupleType
	}
	es := make([]dgo.Value, l)
	for i := range types {
		es[i] = types[i]
	}
	return &tupleType{slice: es, frozen: true}
}

func (t *tupleType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *tupleType) deepAssignable1(guard dgo.RecursionGuard, ot *tupleType) bool {
	es := t.slice
	os := ot.slice
	if len(os) != len(es) {
		return false
	}
	for i := range es {
		if !Assignable(guard, es[i].(dgo.Type), os[i].(dgo.Type)) {
			return false
		}
	}
	return true
}

func (t *tupleType) deepAssignable2(guard dgo.RecursionGuard, ot *exactArrayType) bool {
	es := t.slice
	os := ot.slice
	if len(os) != len(es) {
		return false
	}
	for i := range es {
		if !Instance(guard, es[i].(dgo.Type), os[i]) {
			return false
		}
	}
	return true

}

func (t *tupleType) deepAssignable3(guard dgo.RecursionGuard, ot *sizedArrayType) bool {
	es := t.slice
	if ot.Min() == len(es) && ot.Max() == len(es) {
		et := ot.ElementType()
		for i := range es {
			if !Assignable(guard, es[i].(dgo.Type), et) {
				return false
			}
		}
	}
	return true
}

func (t *tupleType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	dflt := len(t.slice) == 0
	switch ot := other.(type) {
	case defaultArrayType:
		return dflt
	case *tupleType:
		if dflt {
			return true
		}
		return t.deepAssignable1(guard, ot)
	case *exactArrayType:
		if dflt {
			return true
		}
		return t.deepAssignable2(guard, ot)
	case *sizedArrayType:
		if dflt {
			return true
		}
		return t.deepAssignable3(guard, ot)
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *tupleType) ElementType() dgo.Type {
	switch len(t.slice) {
	case 0:
		return DefaultAnyType
	case 1:
		return t.slice[0].(dgo.Type)
	default:
		return (*allOfType)(t)
	}
}

func (t *tupleType) ElementTypes() dgo.Array {
	return (*array)(t)
}

func (t *tupleType) Equals(other interface{}) bool {
	if ot, ok := other.(*tupleType); ok {
		return (*array)(t).Equals((*array)(ot))
	}
	return false
}

func (t *tupleType) HashCode() int {
	return (*array)(t).HashCode()*7 + int(dgo.TiTuple)
}

func (t *tupleType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *tupleType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	if ov, ok := value.(*array); ok {
		es := t.slice
		if len(es) == 0 {
			return true
		}
		s := ov.slice
		if len(s) == len(es) {
			for i := range es {
				if !Instance(guard, es[i].(dgo.Type), s[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (t *tupleType) Max() int {
	if l := len(t.slice); l > 0 {
		return l
	}
	return math.MaxInt64
}

func (t *tupleType) Min() int {
	return len(t.slice)
}

func (t *tupleType) Resolve(ap dgo.AliasProvider) {
	resolveSlice(t.slice, ap)
}

func (t *tupleType) String() string {
	return TypeString(t)
}

func (t *tupleType) Type() dgo.Type {
	return &metaType{t}
}

func (t *tupleType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiTuple
}

func (t *tupleType) Unbounded() bool {
	return len(t.slice) == 0
}

// Array creates a new frozen array that contains a copy of the given slice
func Array(values []dgo.Value) dgo.Array {
	arr := make([]dgo.Value, len(values))
	for i := range values {
		e := values[i]
		if f, ok := e.(dgo.Freezable); ok {
			e = f.FrozenCopy()
		} else if e == nil {
			e = Nil
		}
		arr[i] = e
	}
	return &array{slice: arr, frozen: true}
}

// ArrayFromReflected creates a new array that contains a copy of the given reflected slice
func ArrayFromReflected(vr reflect.Value, frozen bool) dgo.Value {
	if vr.IsNil() {
		return Nil
	}

	ix := vr.Interface()
	if bs, ok := ix.([]byte); ok {
		return Binary(bs, frozen)
	}

	top := vr.Len()
	var arr []dgo.Value
	if vs, ok := ix.([]dgo.Value); ok {
		arr = vs
		if frozen {
			arr = sliceCopy(arr)
		}
	} else {
		arr = make([]dgo.Value, top)
		for i := 0; i < top; i++ {
			arr[i] = ValueFromReflected(vr.Index(i))
		}
	}

	if frozen {
		for i := range arr {
			if f, ok := arr[i].(dgo.Freezable); ok {
				arr[i] = f.FrozenCopy()
			}
		}
	}
	return &array{slice: arr, frozen: frozen}
}

// MutableArray creates a new mutable array that wraps the given slice. Unset entries in the
// slice will be replaced by Nil.
func MutableArray(t dgo.ArrayType, values []dgo.Value) dgo.Array {
	if t != nil {
		l := len(values)
		if l < t.Min() {
			panic(IllegalSize(t, l))
		}
		if l > t.Max() {
			panic(IllegalSize(t, l))
		}
		if tt, ok := t.(*tupleType); ok {
			es := tt.slice
			for i := range values {
				e := values[i]
				et := es[i].(dgo.Type)
				if !et.Instance(e) {
					panic(IllegalAssignment(et, e))
				}
			}
		} else {
			et := t.ElementType()
			if DefaultAnyType != et {
				for i := range values {
					e := values[i]
					if !et.Instance(e) {
						panic(IllegalAssignment(et, e))
					}
				}
			}
		}
	}
	return &array{slice: ReplaceNil(values), typ: t, frozen: false}
}

// MutableValues returns a frozen dgo.Array that represents the given values
func MutableValues(t dgo.ArrayType, values []interface{}) dgo.Array {
	s := make([]dgo.Value, len(values))
	for i := range values {
		s[i] = Value(values[i])
	}
	return MutableArray(t, s)
}

func valueSlice(values []interface{}, frozen bool) []dgo.Value {
	cp := make([]dgo.Value, len(values))
	if frozen {
		for i := range values {
			v := Value(values[i])
			if f, ok := v.(dgo.Freezable); ok {
				v = f.FrozenCopy()
			}
			cp[i] = v
		}
	} else {
		for i := range values {
			cp[i] = Value(values[i])
		}
	}
	return cp
}

// Integers returns a dgo.Array that represents the given ints
func Integers(values []int) dgo.Array {
	cp := make([]dgo.Value, len(values))
	for i := range values {
		cp[i] = intVal(values[i])
	}
	return &array{slice: cp, frozen: true}
}

// Strings returns a dgo.Array that represents the given strings
func Strings(values []string) dgo.Array {
	cp := make([]dgo.Value, len(values))
	for i := range values {
		cp[i] = makeHString(values[i])
	}
	return &array{slice: cp, frozen: true}
}

// Values returns a frozen dgo.Array that represents the given values
func Values(values []interface{}) dgo.Array {
	return &array{slice: valueSlice(values, true), frozen: true}
}

func (v *array) assertType(e dgo.Value, pos int) {
	if t := v.typ; t != nil {
		sz := len(v.slice)
		if pos >= sz {
			sz++
			if sz > t.Max() {
				panic(IllegalSize(t, sz))
			}
		}
		var et dgo.Type
		if tp, ok := t.(*tupleType); ok {
			et = tp.slice[pos].(dgo.Type)
		} else {
			et = t.ElementType()
		}
		if !et.Instance(e) {
			panic(IllegalAssignment(et, e))
		}
	}
}

func (v *array) assertTypes(values dgo.Array) {
	if t := v.typ; t != nil {
		addedSize := values.Len()
		if addedSize == 0 {
			return
		}
		sz := len(v.slice)
		if sz+addedSize > t.Max() {
			panic(IllegalSize(t, sz+addedSize))
		}
		et := t.ElementType()
		for i := 0; i < addedSize; i++ {
			e := values.Get(i)
			if !et.Instance(e) {
				panic(IllegalAssignment(et, e))
			}
		}
	}
}

func (v *array) Add(vi interface{}) {
	if v.frozen {
		panic(frozenArray(`Add`))
	}
	val := Value(vi)
	v.assertType(val, len(v.slice))
	v.slice = append(v.slice, val)
}

func (v *array) AddAll(values dgo.Array) {
	if v.frozen {
		panic(frozenArray(`AddAll`))
	}
	v.assertTypes(values)
	v.slice = values.AppendToSlice(v.slice)
}

func (v *array) AddValues(values ...interface{}) {
	if v.frozen {
		panic(frozenArray(`AddValues`))
	}
	va := valueSlice(values, false)
	v.assertTypes(&array{slice: va})
	v.slice = append(v.slice, va...)
}

func (v *array) All(predicate dgo.Predicate) bool {
	a := v.slice
	for i := range a {
		if !predicate(a[i]) {
			return false
		}
	}
	return true
}

func (v *array) Any(predicate dgo.Predicate) bool {
	a := v.slice
	for i := range a {
		if predicate(a[i]) {
			return true
		}
	}
	return false
}

func (v *array) AppendTo(w util.Indenter) {
	w.AppendRune('[')
	ew := w.Indent()
	a := v.slice
	for i := range a {
		if i > 0 {
			w.AppendRune(',')
		}
		ew.NewLine()
		ew.AppendValue(v.slice[i])
	}
	w.NewLine()
	w.AppendRune(']')
}

func (v *array) AppendToSlice(slice []dgo.Value) []dgo.Value {
	return append(slice, v.slice...)
}

func (v *array) CompareTo(other interface{}) (int, bool) {
	return compare(nil, v, Value(other))
}

func (v *array) deepCompare(seen []dgo.Value, other deepCompare) (int, bool) {
	ov, ok := other.(*array)
	if !ok {
		return 0, false
	}
	a := v.slice
	b := ov.slice
	top := len(a)
	max := len(b)
	r := 0
	if top < max {
		r = -1
		max = top
	} else if top > max {
		r = 1
	}

	for i := 0; i < max; i++ {
		if _, ok = a[i].(dgo.Comparable); !ok {
			r = 0
			break
		}
		var c int
		if c, ok = compare(seen, a[i], b[i]); !ok {
			r = 0
			break
		}
		if c != 0 {
			r = c
			break
		}
	}
	return r, ok
}

func (v *array) Copy(frozen bool) dgo.Array {
	if frozen && v.frozen {
		return v
	}
	cp := sliceCopy(v.slice)
	if frozen {
		for i := range cp {
			if f, ok := cp[i].(dgo.Freezable); ok {
				cp[i] = f.FrozenCopy()
			}
		}
	}
	return &array{slice: cp, typ: v.typ, frozen: frozen}
}

func (v *array) Each(doer dgo.Doer) {
	a := v.slice
	for i := range a {
		doer(a[i])
	}
}

func (v *array) EachWithIndex(doer dgo.DoWithIndex) {
	a := v.slice
	for i := range a {
		doer(a[i], i)
	}
}

func (v *array) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *array) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ov, ok := other.(*array); ok {
		return sliceEquals(seen, v.slice, ov.slice)
	}
	return false
}

func (v *array) Freeze() {
	if v.frozen {
		return
	}
	v.frozen = true
	a := v.slice
	for i := range a {
		if f, ok := a[i].(dgo.Freezable); ok {
			f.Freeze()
		}
	}
}

func (v *array) Frozen() bool {
	return v.frozen
}

func (v *array) FrozenCopy() dgo.Value {
	return v.Copy(true)
}

func (v *array) GoSlice() []dgo.Value {
	if v.frozen {
		return sliceCopy(v.slice)
	}
	return v.slice
}

func (v *array) HashCode() int {
	return v.deepHashCode(nil)
}

func (v *array) deepHashCode(seen []dgo.Value) int {
	h := 1
	s := v.slice
	for i := range s {
		h = h*31 + deepHashCode(seen, s[i])
	}
	return h
}

func (v *array) Get(index int) dgo.Value {
	return v.slice[index]
}

func (v *array) IndexOf(vi interface{}) int {
	val := Value(vi)
	a := v.slice
	for i := range a {
		if val.Equals(a[i]) {
			return i
		}
	}
	return -1
}

func (v *array) Insert(pos int, vi interface{}) {
	if v.frozen {
		panic(frozenArray(`Insert`))
	}
	val := Value(vi)
	v.assertType(val, pos)
	v.slice = append(v.slice[:pos], append([]dgo.Value{val}, v.slice[pos:]...)...)
}

func (v *array) Len() int {
	return len(v.slice)
}

func (v *array) MapTo(t dgo.ArrayType, mapper dgo.Mapper) dgo.Array {
	if t == nil {
		return v.Map(mapper)
	}
	a := v.slice
	l := len(a)
	if l < t.Min() {
		panic(IllegalSize(t, l))
	}
	if l > t.Max() {
		panic(IllegalSize(t, l))
	}
	et := t.ElementType()
	vs := make([]dgo.Value, len(a))

	for i := range a {
		mv := Value(mapper(a[i]))
		if !et.Instance(mv) {
			panic(IllegalAssignment(et, mv))
		}
		vs[i] = mv
	}
	return &array{slice: vs, typ: t, frozen: v.frozen}
}

func (v *array) Map(mapper dgo.Mapper) dgo.Array {
	a := v.slice
	vs := make([]dgo.Value, len(a))
	for i := range a {
		vs[i] = Value(mapper(a[i]))
	}
	return &array{slice: vs, frozen: v.frozen}
}

func (v *array) One(predicate dgo.Predicate) bool {
	a := v.slice
	f := false
	for i := range a {
		if predicate(a[i]) {
			if f {
				return false
			}
			f = true
		}
	}
	return f
}

func (v *array) Reduce(mi interface{}, reductor func(memo dgo.Value, elem dgo.Value) interface{}) dgo.Value {
	memo := Value(mi)
	a := v.slice
	for i := range a {
		memo = Value(reductor(memo, a[i]))
	}
	return memo
}

func (v *array) removePos(pos int) dgo.Value {
	a := v.slice
	if pos >= 0 && pos < len(a) {
		newLen := len(a) - 1
		if v.typ != nil {
			if v.typ.Min() > newLen {
				panic(IllegalSize(v.typ, newLen))
			}
		}
		val := a[pos]
		copy(a[pos:], a[pos+1:])
		a[newLen] = nil // release to GC
		v.slice = a[:newLen]
		return val
	}
	return nil
}

func (v *array) Remove(pos int) dgo.Value {
	if v.frozen {
		panic(frozenArray(`Remove`))
	}
	return v.removePos(pos)
}

func (v *array) RemoveValue(value interface{}) bool {
	if v.frozen {
		panic(frozenArray(`RemoveValue`))
	}
	return v.removePos(v.IndexOf(value)) != nil
}

func (v *array) Reject(predicate dgo.Predicate) dgo.Array {
	vs := make([]dgo.Value, 0)
	a := v.slice
	for i := range a {
		e := a[i]
		if !predicate(e) {
			vs = append(vs, e)
		}
	}
	return &array{slice: vs, typ: v.typ, frozen: v.frozen}
}

func (v *array) SameValues(other dgo.Array) bool {
	return len(v.slice) == other.Len() && v.ContainsAll(other)
}

func (v *array) ContainsAll(other dgo.Array) bool {
	oa := other.(*array)
	a := v.slice
	b := oa.slice
	l := len(a)
	if l < len(b) {
		return false
	}
	if l == 0 {
		return true
	}

	// Keep track of elements that have been found equal using a copy
	// where such elements are set to nil. This avoids excessive calls
	// to Equals
	vs := sliceCopy(b)
	for i := range b {
		ea := a[i]
		f := false
		for j := range vs {
			if be := vs[j]; be != nil {
				if be.Equals(ea) {
					vs[j] = nil
					f = true
					break
				}
			}
		}
		if !f {
			return false
		}
	}
	return true
}

func (v *array) Select(predicate dgo.Predicate) dgo.Array {
	vs := make([]dgo.Value, 0)
	a := v.slice
	for i := range a {
		e := a[i]
		if predicate(e) {
			vs = append(vs, e)
		}
	}
	return &array{slice: vs, typ: v.typ, frozen: v.frozen}
}

func (v *array) Set(pos int, vi interface{}) dgo.Value {
	if v.frozen {
		panic(frozenArray(`Set`))
	}
	val := Value(vi)
	v.assertType(val, pos)
	old := v.slice[pos]
	v.slice[pos] = val
	return old
}

func (v *array) SetType(t dgo.ArrayType) {
	if v.frozen {
		panic(frozenArray(`SetType`))
	}
	if t.Instance(v) {
		v.typ = t
		return
	}
	panic(IllegalAssignment(t, v))
}

func (v *array) Sort() dgo.Array {
	sa := v.slice
	if len(sa) < 2 {
		return v
	}
	sorted := sliceCopy(sa)
	sort.SliceStable(sorted, func(i, j int) bool {
		a := sorted[i]
		b := sorted[j]
		if ac, ok := a.(dgo.Comparable); ok {
			var c int
			if c, ok = ac.CompareTo(b); ok {
				return c < 0
			}
		}
		return a.Type().TypeIdentifier() < b.Type().TypeIdentifier()
	})
	return &array{slice: sorted, typ: v.typ, frozen: v.frozen}
}

func (v *array) String() string {
	return util.ToStringERP(v)
}

func (v *array) ToMap() dgo.Map {
	ms := v.slice
	top := len(ms)

	ts := top / 2
	if top%2 != 0 {
		ts++
	}
	tbl := make([]*hashNode, tableSizeFor(ts))
	hl := len(tbl) - 1
	m := &hashMap{table: tbl, len: ts, frozen: v.frozen}

	for i := 0; i < top; {
		mk := ms[i]
		i++
		var mv dgo.Value = Nil
		if i < top {
			mv = ms[i]
			i++
		}
		hk := hl & hash(mk.HashCode())
		nd := &hashNode{mapEntry: mapEntry{key: mk, value: mv}, hashNext: tbl[hk], prev: m.last}
		if m.first == nil {
			m.first = nd
		} else {
			m.last.next = nd
		}
		m.last = nd
		tbl[hk] = nd
	}
	return m
}

func (v *array) ToMapFromEntries() (dgo.Map, bool) {
	ms := v.slice
	top := len(ms)
	tbl := make([]*hashNode, tableSizeFor(top))
	hl := len(tbl) - 1
	m := &hashMap{table: tbl, len: top, frozen: v.frozen}

	for i := range ms {
		nd, ok := ms[i].(*hashNode)
		if !ok {
			var ea *array
			if ea, ok = ms[i].(*array); ok && len(ea.slice) == 2 {
				nd = &hashNode{mapEntry: mapEntry{key: ea.slice[0], value: ea.slice[1]}}
			} else {
				return nil, false
			}
		} else if nd.hashNext != nil {
			// Copy node, it belongs to another map
			c := *nd
			c.next = nil // this one might not get assigned below
			nd = &c
		}

		hk := hl & hash(nd.key.HashCode())
		nd.hashNext = tbl[hk]
		nd.prev = m.last
		if m.first == nil {
			m.first = nd
		} else {
			m.last.next = nd
		}
		m.last = nd
		tbl[hk] = nd
	}
	return m, true
}

func (v *array) Type() dgo.Type {
	if v.typ == nil {
		return (*exactArrayType)(v)
	}
	return v.typ
}

func (v *array) Unique() dgo.Array {
	a := v.slice
	top := len(a)
	if top < 2 {
		return v
	}
	tbl := make([]*hashNode, tableSizeFor(int(float64(top)/loadFactor)))
	hl := len(tbl) - 1
	u := make([]dgo.Value, top)
	ui := 0

nextVal:
	for i := range a {
		k := a[i]
		hk := hl & hash(k.HashCode())
		for e := tbl[hk]; e != nil; e = e.hashNext {
			if k.Equals(e.key) {
				continue nextVal
			}
		}
		tbl[hk] = &hashNode{mapEntry: mapEntry{key: k}, hashNext: tbl[hk]}
		u[ui] = k
		ui++
	}
	if ui == top {
		return v
	}
	return &array{slice: u[:ui], typ: v.typ, frozen: v.frozen}
}

func (v *array) MarshalJSON() ([]byte, error) {
	return []byte(util.ToStringERP(v)), nil
}

func (v *array) MarshalYAML() (interface{}, error) {
	a := v.slice
	s := make([]*yaml.Node, len(a))
	var err error
	for i := range a {
		s[i], err = yamlEncodeValue(a[i])
		if err != nil {
			return nil, err
		}
	}
	return &yaml.Node{Kind: yaml.SequenceNode, Tag: `!!seq`, Content: s}, nil
}

func (v *array) Pop() (dgo.Value, bool) {
	if v.frozen {
		panic(frozenArray(`Pop`))
	}
	p := len(v.slice) - 1
	if p >= 0 {
		return v.removePos(p), true
	}
	return nil, false
}

func (v *array) UnmarshalJSON(b []byte) error {
	if v.frozen {
		panic(frozenArray(`UnmarshalJSON`))
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	t, err := dec.Token()
	if err == nil {
		if delim, ok := t.(json.Delim); !ok || delim != '[' {
			return errors.New("expecting data to be an array")
		}
		var a *array
		a, err = jsonDecodeArray(dec)
		if err == nil {
			*v = *a
		}
	}
	return err
}

func (v *array) UnmarshalYAML(n *yaml.Node) error {
	if v.frozen {
		panic(frozenArray(`UnmarshalYAML`))
	}
	if n.Kind != yaml.SequenceNode {
		return errors.New("expecting data to be an array")
	}
	a, err := yamlDecodeArray(n)
	if err == nil {
		*v = *a
	}
	return err
}

func (v *array) With(vi interface{}) dgo.Array {
	val := Value(vi)
	v.assertType(val, len(v.slice))
	return &array{slice: append(v.slice, val), typ: v.typ, frozen: v.frozen}
}

func (v *array) WithAll(values dgo.Array) dgo.Array {
	if values.Len() == 0 {
		return v
	}
	v.assertTypes(values)
	return &array{slice: values.AppendToSlice(v.slice), typ: v.typ, frozen: v.frozen}
}

func (v *array) WithValues(values ...interface{}) dgo.Array {
	if len(values) == 0 {
		return v
	}
	va := valueSlice(values, v.frozen)
	v.assertTypes(&array{slice: va})
	return &array{slice: append(v.slice, va...), typ: v.typ, frozen: v.frozen}
}

// ReplaceNil performs an in-place replacement of nil interfaces with the NilValue
func ReplaceNil(vs []dgo.Value) []dgo.Value {
	for i := range vs {
		if vs[i] == nil {
			vs[i] = Nil
		}
	}
	return vs
}

// allInstance returns true when all elements of slice vs are assignable to the given type t
func allInstance(guard dgo.RecursionGuard, t dgo.Type, vs []dgo.Value) bool {
	if t == DefaultAnyType {
		return true
	}
	for i := range vs {
		if !Instance(guard, t, vs[i]) {
			return false
		}
	}
	return true
}

// allAssignable returns true when all types in the given slice s are assignable to the given type t
func allAssignable(guard dgo.RecursionGuard, t dgo.Type, s []dgo.Value) bool {
	for i := range s {
		if !Assignable(guard, t, s[i].(dgo.Type)) {
			return false
		}
	}
	return true
}

// assignableToAll returns true when the given type t is assignable the type of all elements of slice vs
func assignableToAll(guard dgo.RecursionGuard, t dgo.Type, vs []dgo.Value) bool {
	for i := range vs {
		if !Assignable(guard, vs[i].Type(), t) {
			return false
		}
	}
	return true
}

func frozenArray(f string) error {
	return fmt.Errorf(`%s called on a frozen Array`, f)
}

func sliceCopy(s []dgo.Value) []dgo.Value {
	c := make([]dgo.Value, len(s))
	copy(c, s)
	return c
}

func resolveSlice(ts []dgo.Value, ap dgo.AliasProvider) {
	for i := range ts {
		ts[i] = ap.Replace(ts[i].(dgo.Type))
	}
}
