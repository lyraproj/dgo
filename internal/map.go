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

const initialCapacity = 1 << 4
const maximumCapacity = 1 << 30
const loadFactor = 0.75

type (
	// defaultMapType is the unconstrained map type
	defaultMapType int

	// exactMapType represents a map exactly
	exactMapType hashMap

	// sizedMapType represents a map with constraints on key type, value type, and size
	sizedMapType struct {
		keyType   dgo.Type
		valueType dgo.Type
		min       int
		max       int
	}

	mapEntry struct {
		key   dgo.Value
		value dgo.Value
	}

	exactEntryType mapEntry

	hashNode struct {
		mapEntry
		hashNext *hashNode
		next     *hashNode
		prev     *hashNode
	}

	// hashMap is an unsorted Map that uses a hash table
	hashMap struct {
		table  []*hashNode
		typ    dgo.MapType
		len    int
		first  *hashNode
		last   *hashNode
		frozen bool
	}
)

func (t *exactEntryType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*exactEntryType); ok {
		return (*mapEntry)(t).Equals((*mapEntry)(ot))
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *exactEntryType) Equals(other interface{}) bool {
	if ot, ok := other.(*exactEntryType); ok {
		return (*mapEntry)(t).Equals((*mapEntry)(ot))
	}
	return false
}

func (t *exactEntryType) HashCode() int {
	return (*mapEntry)(t).HashCode()
}

func (t *exactEntryType) Instance(value interface{}) bool {
	if ot, ok := value.(*hashNode); ok {
		return (*mapEntry)(t).Equals(ot)
	}
	return false
}

func (t *exactEntryType) KeyType() dgo.Type {
	return (*mapEntry)(t).key.Type()
}

func (t *exactEntryType) Required() bool {
	return true
}

func (t *exactEntryType) String() string {
	return TypeString(t)
}

func (t *exactEntryType) Type() dgo.Type {
	return &metaType{t}
}

func (t *exactEntryType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMapEntryExact
}

func (t *exactEntryType) Value() dgo.Value {
	v := (*mapEntry)(t)
	return v
}

func (t *exactEntryType) ValueType() dgo.Type {
	return (*mapEntry)(t).value.Type()
}

func (v *mapEntry) AppendTo(w util.Indenter) {
	w.AppendValue(v.key)
	w.Append(`:`)
	if w.Indenting() {
		w.Append(` `)
	}
	w.AppendValue(v.value)
}

func (v *mapEntry) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *mapEntry) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ov, ok := other.(*hashNode); ok {
		return equals(seen, v.key, ov.key) && equals(seen, v.value, ov.value)
	}
	return false
}

func (v *mapEntry) HashCode() int {
	return deepHashCode(nil, v)
}

func (v *mapEntry) deepHashCode(seen []dgo.Value) int {
	return deepHashCode(seen, v.key) ^ deepHashCode(seen, v.value)
}

func (v *mapEntry) Freeze() {
	if f, ok := v.value.(dgo.Freezable); ok {
		f.Freeze()
	}
}

func (v *mapEntry) Frozen() bool {
	if f, ok := v.value.(dgo.Freezable); ok && !f.Frozen() {
		return false
	}
	return true
}

func (v *mapEntry) FrozenCopy() dgo.Value {
	if v.Frozen() {
		return v
	}
	c := &mapEntry{key: v.key, value: v.value}
	c.copyFreeze()
	return c
}

func (v *hashNode) FrozenCopy() dgo.Value {
	if v.Frozen() {
		return v
	}
	c := &hashNode{mapEntry: mapEntry{key: v.key, value: v.value}}
	c.copyFreeze()
	return c
}

func (v *mapEntry) Key() dgo.Value {
	return v.key
}

func (v *mapEntry) String() string {
	return util.ToStringERP(v)
}

func (v *mapEntry) Type() dgo.Type {
	return (*exactEntryType)(v)
}

func (v *mapEntry) Value() dgo.Value {
	return v.value
}

func (v *mapEntry) copyFreeze() {
	if f, ok := v.value.(dgo.Freezable); ok {
		v.value = f.FrozenCopy()
	}
}

var emptyMap = &hashMap{frozen: true}

// Map creates an immutable dgo.Map from the given argument which must be a go map.
func Map(m interface{}) dgo.Map {
	rm := reflect.ValueOf(m)
	if rm.Kind() != reflect.Map {
		panic(fmt.Errorf(`illegal argument: %tst is not a map`, m))
	}
	return MapFromReflected(rm, true).(dgo.Map)
}

// MapFromReflected creates a Map from a reflected map. If frozen is true, the created Map will be
// immutable and the type will reflect exactly that map and nothing else. If frozen is false, the
// created Map will be mutable and its type will be derived from the reflected map.
func MapFromReflected(rm reflect.Value, frozen bool) dgo.Value {
	if rm.IsNil() {
		return Nil
	}
	keys := rm.MapKeys()
	top := len(keys)
	ic := top
	if top == 0 {
		if frozen {
			return emptyMap
		}
		ic = initialCapacity
	}
	tbl := make([]*hashNode, tableSizeFor(int(float64(ic)/loadFactor)))
	hl := len(tbl) - 1
	se := make([][2]dgo.Value, len(keys))
	for i := range keys {
		key := keys[i]
		se[i] = [2]dgo.Value{ValueFromReflected(key), ValueFromReflected(rm.MapIndex(key))}
	}

	// Sort by key to always get predictable order
	sort.Slice(se, func(i, j int) bool {
		if cm, ok := se[i][0].(dgo.Comparable); ok {
			var c int
			if c, ok = cm.CompareTo(se[j][0]); ok {
				return c < 0
			}
		}
		return false
	})

	m := &hashMap{table: tbl, len: top, frozen: frozen}
	for i := range se {
		e := se[i]
		k := e[0]
		hk := hl & hash(k.HashCode())
		hn := &hashNode{mapEntry: mapEntry{key: k, value: e[1]}, hashNext: tbl[hk], prev: m.last}
		if m.last == nil {
			m.first = hn
		} else {
			m.last.next = hn
		}
		m.last = hn
		tbl[hk] = hn
	}
	var typ dgo.MapType
	if !frozen {
		typ = TypeFromReflected(rm.Type()).(dgo.MapType)
	}
	m.typ = typ
	return m
}

// MutableMap creates an empty dgo.Map with the given capacity. The map can be optionally constrained
// by the given type which can be nil, the zero value of a go map, or a dgo.MapType
func MutableMap(capacity int, typ interface{}) dgo.Map {
	if capacity <= 0 {
		capacity = initialCapacity
	}

	parseMapType := func(s string) dgo.MapType {
		if t, ok := Parse(s).(dgo.MapType); ok {
			return t
		}
		panic(fmt.Errorf("expression '%s' does not evaluate to a MapType", s))
	}
	var mt dgo.MapType
	if typ != nil {
		switch typ := typ.(type) {
		case dgo.MapType:
			mt = typ
		case string:
			mt = parseMapType(typ)
		case dgo.String:
			mt = parseMapType(typ.GoString())
		default:
			mt = TypeFromReflected(reflect.TypeOf(typ)).(dgo.MapType)
		}
	}
	return &hashMap{table: make([]*hashNode, tableSizeFor(capacity)), len: 0, typ: mt, frozen: false}
}

func (g *hashMap) All(predicate dgo.EntryPredicate) bool {
	for e := g.first; e != nil; e = e.next {
		if !predicate(e) {
			return false
		}
	}
	return true
}

func (g *hashMap) AllKeys(predicate dgo.Predicate) bool {
	for e := g.first; e != nil; e = e.next {
		if !predicate(e.key) {
			return false
		}
	}
	return true
}

func (g *hashMap) AllValues(predicate dgo.Predicate) bool {
	for e := g.first; e != nil; e = e.next {
		if !predicate(e.value) {
			return false
		}
	}
	return true
}

func (g *hashMap) Any(predicate dgo.EntryPredicate) bool {
	for e := g.first; e != nil; e = e.next {
		if predicate(e) {
			return true
		}
	}
	return false
}

func (g *hashMap) AnyKey(predicate dgo.Predicate) bool {
	for e := g.first; e != nil; e = e.next {
		if predicate(e.key) {
			return true
		}
	}
	return false
}

func (g *hashMap) AnyValue(predicate dgo.Predicate) bool {
	for e := g.first; e != nil; e = e.next {
		if predicate(e.value) {
			return true
		}
	}
	return false
}

func (g *hashMap) AppendTo(w util.Indenter) {
	w.AppendRune('{')
	ew := w.Indent()
	first := true
	for e := g.first; e != nil; e = e.next {
		if first {
			first = false
		} else {
			ew.AppendRune(',')
		}
		ew.NewLine()
		ew.AppendValue(e)
	}
	w.NewLine()
	w.AppendRune('}')
}

func (g *hashMap) Copy(frozen bool) dgo.Map {
	if frozen && g.frozen {
		return g
	}

	c := &hashMap{len: g.len, frozen: frozen}
	g.resize(c, 0)
	if frozen {
		for e := c.first; e != nil; e = e.next {
			e.copyFreeze()
		}
	}
	return c
}

func (g *hashMap) Each(doer dgo.EntryDoer) {
	for e := g.first; e != nil; e = e.next {
		doer(e)
	}
}

func (g *hashMap) EachKey(doer dgo.Doer) {
	for e := g.first; e != nil; e = e.next {
		doer(e.key)
	}
}

func (g *hashMap) EachValue(doer dgo.Doer) {
	for e := g.first; e != nil; e = e.next {
		doer(e.value)
	}
}

func (g *hashMap) Entries() dgo.Array {
	es := make([]dgo.Value, g.len)
	ei := 0
	for e := g.first; e != nil; e = e.next {
		es[ei] = e.FrozenCopy()
		ei++
	}
	return &array{slice: es, frozen: true}
}

func (g *hashMap) Equals(other interface{}) bool {
	return equals(nil, g, other)
}

func (g *hashMap) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if om, ok := other.(*hashMap); ok && g.len == om.len {
		var ov dgo.Value
		for e := g.first; e != nil; e = e.next {
			ov, ok = om.Get(e.key)
			if !(ok && equals(seen, e.value, ov)) {
				return false
			}
		}
		return true
	}
	return false
}

func (g *hashMap) Freeze() {
	g.frozen = true
	for e := g.first; e != nil; e = e.next {
		e.Freeze()
	}
}

func (g *hashMap) Frozen() bool {
	return g.frozen
}

func (g *hashMap) FrozenCopy() dgo.Value {
	return g.Copy(true)
}

func (g *hashMap) Get(key interface{}) (dgo.Value, bool) {
	tbl := g.table
	tl := len(tbl) - 1
	if tl >= 0 {
		// This switch increases performance a great deal because using the direct implementation
		// instead of the dgo.Value enables inlining of the HashCode() method
		switch k := key.(type) {
		case *hstring:
			for e := tbl[tl&hash(k.HashCode())]; e != nil; e = e.hashNext {
				if k.Equals(e.key) {
					return e.value, true
				}
			}
		case intVal:
			for e := tbl[tl&hash(k.HashCode())]; e != nil; e = e.hashNext {
				if k == e.key {
					return e.value, true
				}
			}
		case string:
			gk := makeHString(k)
			for e := tbl[tl&hash(gk.HashCode())]; e != nil; e = e.hashNext {
				if gk.Equals(e.key) {
					return e.value, true
				}
			}
		default:
			gk := Value(k)
			for e := tbl[tl&hash(gk.HashCode())]; e != nil; e = e.hashNext {
				if gk.Equals(e.key) {
					return e.value, true
				}
			}
		}
	}
	return nil, false
}

func (g *hashMap) HashCode() int {
	return deepHashCode(nil, g)
}

func (g *hashMap) deepHashCode(seen []dgo.Value) int {
	h := 1
	for e := g.first; e != nil; e = e.next {
		h = h*31 + deepHashCode(seen, e)
	}
	return h
}

func (g *hashMap) Keys() dgo.Array {
	return &array{slice: g.keys(), frozen: g.frozen}
}

func (g *hashMap) keys() []dgo.Value {
	ks := make([]dgo.Value, g.len)
	p := 0
	for e := g.first; e != nil; e = e.next {
		ks[p] = e.key
		p++
	}
	return ks
}

func (g *hashMap) Len() int {
	return g.len
}

func (g *hashMap) Map(mapper dgo.EntryMapper) dgo.Map {
	c := &hashMap{len: g.len, frozen: g.frozen}
	g.resize(c, 0)
	for e := c.first; e != nil; e = e.next {
		e.value = Value(mapper(e))
	}
	return c
}

func (g *hashMap) MarshalJSON() ([]byte, error) {
	// JSON is the default Indentable string format, so this one is easy
	return []byte(util.ToStringERP(g)), nil
}

func (g *hashMap) UnmarshalJSON(b []byte) error {
	if g.frozen {
		panic(frozenMap(`UnmarshalJSON`))
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	t, err := dec.Token()
	if err == nil {
		if delim, ok := t.(json.Delim); !ok || delim != '{' {
			return errors.New("expecting data to be an object")
		}
		var m *hashMap
		m, err = jsonDecodeMap(dec)
		if err == nil {
			g.assignFrom(m)
		}
	}
	return err
}

func (g *hashMap) UnmarshalYAML(n *yaml.Node) error {
	if g.frozen {
		panic(frozenMap(`UnmarshalYAML`))
	}
	if n.Kind != yaml.MappingNode {
		return errors.New("expecting data to be an object")
	}
	m, err := yamlDecodeMap(n)
	if err == nil {
		g.assignFrom(m)
	}
	return err
}

func (g *hashMap) assignFrom(m *hashMap) {
	if g.typ != nil {
		*g = hashMap{typ: g.typ}
		g.PutAll(m)
		g.frozen = m.frozen
	} else {
		*g = *m
	}
}

// MarshalYAML returns a *yaml.Node that represents this Map.
func (g *hashMap) MarshalYAML() (interface{}, error) {
	s := make([]*yaml.Node, g.len*2)
	i := 0
	for e := g.first; e != nil; e = e.next {
		n, err := yamlEncodeValue(e.key)
		if err != nil {
			return nil, err
		}
		s[i] = n
		i++
		n, err = yamlEncodeValue(e.value)
		if err != nil {
			return nil, err
		}
		s[i] = n
		i++
	}
	return &yaml.Node{Kind: yaml.MappingNode, Tag: `!!map`, Content: s}, nil
}

func (g *hashMap) Merge(associations dgo.Map) dgo.Map {
	if associations.Len() == 0 || g == associations {
		return g
	}
	l := g.len
	if l == 0 {
		return associations
	}
	c := &hashMap{len: l, typ: g.typ}
	g.resize(c, l+associations.Len())
	c.PutAll(associations)
	c.frozen = g.frozen
	return c
}

func (g *hashMap) Put(ki, vi interface{}) dgo.Value {
	if g.frozen {
		panic(frozenMap(`Put`))
	}
	k := Value(ki)
	v := Value(vi)
	hs := hash(k.HashCode())
	var hk int

	tbl := g.table
	hk = (len(tbl) - 1) & hs
	for e := tbl[hk]; e != nil; e = e.hashNext {
		if k.Equals(e.key) {
			g.assertType(k, v, 0)
			old := e.value
			e.value = v
			return old
		}
	}

	g.assertType(k, v, 1)
	if float64(g.len+1) > float64(len(g.table))*loadFactor {
		g.resize(g, 1)
		tbl = g.table
		hk = (len(tbl) - 1) & hs
	}

	nd := &hashNode{mapEntry: mapEntry{key: frozenCopy(k), value: v}, hashNext: tbl[hk], prev: g.last}
	if g.first == nil {
		g.first = nd
	} else {
		g.last.next = nd
	}
	g.last = nd
	tbl[hk] = nd
	g.len++
	return nil
}

func (g *hashMap) PutAll(associations dgo.Map) {
	al := associations.Len()
	if al == 0 {
		return
	}
	if g.frozen {
		panic(frozenMap(`PutAll`))
	}

	l := g.len
	if float64(l+al) > float64(l)*loadFactor {
		g.resize(g, al)
	}
	tbl := g.table

	associations.Each(func(entry dgo.MapEntry) {
		key := entry.Key()
		val := entry.Value()
		hk := (len(tbl) - 1) & hash(key.HashCode())
		for e := tbl[hk]; e != nil; e = e.hashNext {
			if key.Equals(e.key) {
				g.assertType(key, val, 0)
				e.value = val
				return
			}
		}
		g.assertType(key, val, 1)
		nd := &hashNode{mapEntry: mapEntry{key: frozenCopy(key), value: val}, hashNext: tbl[hk], prev: g.last}
		if g.first == nil {
			g.first = nd
		} else {
			g.last.next = nd
		}
		g.last = nd
		tbl[hk] = nd
		l++
	})
	g.len = l
}

func (g *hashMap) Remove(ki interface{}) dgo.Value {
	if g.frozen {
		panic(frozenMap(`Remove`))
	}
	key := Value(ki)
	hk := (len(g.table) - 1) & hash(key.HashCode())

	var p *hashNode
	for e := g.table[hk]; e != nil; e = e.hashNext {
		if key.Equals(e.key) {
			old := e.value
			if p == nil {
				g.table[hk] = e.hashNext
			} else {
				p.hashNext = e.hashNext
			}
			if e.prev == nil {
				g.first = e.next
			} else {
				e.prev.next = e.next
			}
			if e.next == nil {
				g.last = e.prev
			} else {
				e.next.prev = e.prev
			}
			g.len--
			return old
		}
		p = e
	}
	return nil
}

func (g *hashMap) RemoveAll(keys dgo.Array) {
	if g.frozen {
		panic(frozenMap(`RemoveAll`))
	}
	if g.len == 0 || keys.Len() == 0 {
		return
	}

	tbl := g.table
	kl := len(tbl) - 1
	keys.Each(func(k dgo.Value) {
		hk := kl & hash(k.HashCode())
		var p *hashNode
		for e := tbl[hk]; e != nil; e = e.hashNext {
			if k.Equals(e.key) {
				if p == nil {
					tbl[hk] = e.hashNext
				} else {
					p.hashNext = e.hashNext
				}
				if e.prev == nil {
					g.first = e.next
				} else {
					e.prev.next = e.next
				}
				if e.next == nil {
					g.last = e.prev
				} else {
					e.next.prev = e.prev
				}
				g.len--
				break
			}
			p = e
		}
	})
}

func (g *hashMap) SetType(t dgo.MapType) {
	if g.frozen {
		panic(frozenMap(`SetType`))
	}
	if t.Instance(g) {
		g.typ = t
		return
	}
	panic(IllegalAssignment(t, g))
}

func (g *hashMap) String() string {
	return util.ToStringERP(g)
}

func (g *hashMap) With(ki, vi interface{}) dgo.Map {
	var c *hashMap
	key := Value(ki)
	val := Value(vi)
	if g.table == nil {
		c = &hashMap{table: make([]*hashNode, tableSizeFor(initialCapacity)), len: g.len, typ: g.typ}
	} else {
		if ov, ok := g.Get(key); ok && ov.Equals(val) {
			return g
		}
		c = &hashMap{len: g.len, typ: g.typ}
		g.resize(c, 1)
	}
	c.Put(key, val)
	c.frozen = g.frozen
	return c
}

func (g *hashMap) Without(ki interface{}) dgo.Map {
	key := Value(ki)
	if _, ok := g.Get(key); !ok {
		return g
	}
	c := &hashMap{typ: g.typ, len: g.len}
	g.resize(c, 0)
	c.Remove(key)
	c.frozen = g.frozen
	return c
}

func (g *hashMap) WithoutAll(keys dgo.Array) dgo.Map {
	if g.len == 0 || keys.Len() == 0 {
		return g
	}
	c := &hashMap{typ: g.typ, len: g.len}
	g.resize(c, 0)
	c.RemoveAll(keys)
	if g.len == c.len {
		return g
	}
	c.frozen = g.frozen
	return c
}

func (g *hashMap) Type() dgo.Type {
	if g.typ == nil {
		return (*exactMapType)(g)
	}
	return g.typ
}

func (g *hashMap) Values() dgo.Array {
	return &array{slice: g.values(), frozen: g.frozen}
}

func (g *hashMap) assertType(k, v dgo.Value, addedSize int) {
	if t := g.typ; t != nil {
		kt := t.KeyType()
		if !kt.Instance(k) {
			panic(IllegalAssignment(kt, k))
		}
		vt := t.ValueType()
		if !vt.Instance(v) {
			panic(IllegalAssignment(vt, v))
		}
		if addedSize > 0 {
			sz := g.len + addedSize
			if sz > t.Max() {
				panic(IllegalSize(t, sz))
			}
		}
	}
}

func (g *hashMap) values() []dgo.Value {
	ks := make([]dgo.Value, g.len)
	p := 0
	for e := g.first; e != nil; e = e.next {
		ks[p] = e.value
		p++
	}
	return ks
}

func (g *hashMap) resize(c *hashMap, capInc int) {
	tbl := g.table
	tl := tableSizeFor(len(tbl) + capInc)
	nt := make([]*hashNode, tl)
	c.table = nt
	tl--
	var prev *hashNode
	for oe := g.first; oe != nil; oe = oe.next {
		hk := tl & hash(oe.key.HashCode())
		ne := &hashNode{mapEntry: mapEntry{key: oe.key, value: oe.value}, hashNext: nt[hk]}
		if prev == nil {
			c.first = ne
		} else {
			prev.next = ne
			ne.prev = prev
		}
		nt[hk] = ne
		prev = ne
	}
	c.last = prev
}

func tableSizeFor(cap int) int {
	if cap < 1 {
		return 1
	}
	n := (uint)(cap - 1)
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	if n > maximumCapacity {
		return maximumCapacity
	}
	return int(n)
}

func frozenCopy(v dgo.Value) dgo.Value {
	if f, ok := v.(dgo.Freezable); ok {
		v = f.FrozenCopy()
	}
	return v
}

func hash(h int) int {
	return h ^ (h >> 16)
}

func mapTypeOne(args []interface{}) dgo.MapType {
	// min integer
	a0, ok := Value(args[0]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Map`, `Integer`, args, 0))
	}
	return newMapType(nil, nil, int(a0.GoInt()), math.MaxInt64)
}

func mapTypeTwo(args []interface{}) dgo.MapType {
	// key and value types or min and max integers
	switch a0 := Value(args[0]).(type) {
	case dgo.Type:
		a1, ok := Value(args[1]).(dgo.Type)
		if !ok {
			panic(illegalArgument(`Map`, `Type`, args, 1))
		}
		return newMapType(a0, a1, 0, math.MaxInt64)
	case dgo.Integer:
		a1, ok := Value(args[1]).(dgo.Integer)
		if !ok {
			panic(illegalArgument(`Map`, `Integer`, args, 1))
		}
		return newMapType(nil, nil, int(a0.GoInt()), int(a1.GoInt()))
	default:
		panic(illegalArgument(`Map`, `Type or Integer`, args, 0))
	}
}

func mapTypeThree(args []interface{}) dgo.MapType {
	// key and value types, and min integer
	a0, ok := Value(args[0]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Map`, `Type`, args, 0))
	}
	a1, ok := Value(args[1]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Map`, `Type`, args, 1))
	}
	a2, ok := Value(args[2]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Map`, `Integer`, args, 2))
	}
	return newMapType(a0, a1, int(a2.GoInt()), math.MaxInt64)
}

func mapTypeFour(args []interface{}) dgo.MapType {
	// key and value types, and min and max integers
	a0, ok := Value(args[0]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Map`, `Type`, args, 0))
	}
	a1, ok := Value(args[1]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Map`, `Type`, args, 1))
	}
	a2, ok := Value(args[2]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Map`, `Integer`, args, 2))
	}
	a3, ok := Value(args[3]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Map`, `Integer`, args, 3))
	}
	return newMapType(a0, a1, int(a2.GoInt()), int(a3.GoInt()))
}

// MapType returns a type that represents an Map value
func MapType(args ...interface{}) dgo.MapType {
	switch len(args) {
	case 0:
		return DefaultMapType
	case 1:
		return mapTypeOne(args)
	case 2:
		return mapTypeTwo(args)
	case 3:
		return mapTypeThree(args)
	case 4:
		return mapTypeFour(args)
	default:
		panic(fmt.Errorf(`illegal number of arguments for MapType. Expected 0 - 4, got %d`, len(args)))
	}
}

func newMapType(kt, vt dgo.Type, min, max int) dgo.MapType {
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
	if kt == nil {
		kt = DefaultAnyType
	}
	if vt == nil {
		vt = DefaultAnyType
	}
	if min == 0 && max == math.MaxInt64 {
		// Unbounded
		if kt == DefaultAnyType && vt == DefaultAnyType {
			return DefaultMapType
		}
	}
	return &sizedMapType{keyType: kt, valueType: vt, min: min, max: max}
}

func (t *sizedMapType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *sizedMapType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	if ot, ok := other.(*sizedMapType); ok {
		return t.min <= ot.min && ot.max <= t.max && Assignable(guard, t.keyType, ot.keyType) && Assignable(guard, t.valueType, ot.valueType)
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *sizedMapType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *sizedMapType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*sizedMapType); ok {
		return t.min == ot.min && t.max == ot.max && equals(seen, t.keyType, ot.keyType) && equals(seen, t.valueType, ot.valueType)
	}
	return false
}

func (t *sizedMapType) HashCode() int {
	return deepHashCode(nil, t)
}

func (t *sizedMapType) deepHashCode(seen []dgo.Value) int {
	h := int(dgo.TiMap)
	if t.min > 0 {
		h = h*31 + t.min
	}
	if t.max < math.MaxInt64 {
		h = h*31 + t.max
	}
	if DefaultAnyType != t.keyType {
		h = h*31 + deepHashCode(seen, t.keyType)
	}
	if DefaultAnyType != t.valueType {
		h = h*31 + deepHashCode(seen, t.keyType)
	}
	return h
}

func (t *sizedMapType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *sizedMapType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	if ov, ok := value.(*hashMap); ok {
		l := ov.len
		if t.min <= l && l <= t.max {
			kt := t.keyType
			vt := t.valueType
			if DefaultAnyType == kt {
				if DefaultAnyType == vt {
					return true
				}
				return ov.AllValues(func(v dgo.Value) bool { return Instance(guard, vt, v) })
			}
			if DefaultAnyType == vt {
				return ov.AllKeys(func(k dgo.Value) bool { return kt.Instance(k) })
			}
			return ov.All(func(e dgo.MapEntry) bool { return Instance(guard, kt, e.Key()) && Instance(guard, vt, e.Value()) })
		}
	}
	return false
}

func (t *sizedMapType) KeyType() dgo.Type {
	return t.keyType
}

func (t *sizedMapType) Max() int {
	return t.max
}

func (t *sizedMapType) Min() int {
	return t.min
}

func (t *sizedMapType) Resolve(ap dgo.AliasProvider) {
	t.keyType = ap.Replace(t.keyType)
	t.valueType = ap.Replace(t.valueType)
}

func (t *sizedMapType) String() string {
	return TypeString(t)
}

func (t *sizedMapType) Type() dgo.Type {
	return &metaType{t}
}

func (t *sizedMapType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMapSized
}

func (t *sizedMapType) ValueType() dgo.Type {
	return t.valueType
}

func (t *sizedMapType) Unbounded() bool {
	return t.min == 0 && t.max == math.MaxInt64
}

// DefaultMapType is the unconstrained Map type
const DefaultMapType = defaultMapType(0)

func (t defaultMapType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case defaultMapType, *sizedMapType:
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t defaultMapType) Equals(other interface{}) bool {
	return t == other
}

func (t defaultMapType) HashCode() int {
	return int(dgo.TiMap)
}

func (t defaultMapType) Instance(value interface{}) bool {
	_, ok := value.(dgo.Map)
	return ok
}

func (t defaultMapType) KeyType() dgo.Type {
	return DefaultAnyType
}

func (t defaultMapType) Max() int {
	return math.MaxInt64
}

func (t defaultMapType) Min() int {
	return 0
}

func (t defaultMapType) String() string {
	return TypeString(t)
}

func (t defaultMapType) Type() dgo.Type {
	return &metaType{t}
}

func (t defaultMapType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMap
}

func (t defaultMapType) ValueType() dgo.Type {
	return DefaultAnyType
}

func (t defaultMapType) Unbounded() bool {
	return true
}

func (t *exactMapType) Assignable(other dgo.Type) bool {
	return CheckAssignableTo(nil, other, t)
}

func (t *exactMapType) AssignableTo(guard dgo.RecursionGuard, other dgo.Type) bool {
	switch ot := other.(type) {
	case defaultMapType:
		return true
	case *exactMapType:
		return t.Equals(ot)
	case *sizedMapType:
		m := (*hashMap)(t)
		return ot.Instance(m)
	}
	return false
}

func (t *exactMapType) Equals(other interface{}) bool {
	if ot, ok := other.(*exactMapType); ok {
		return (*hashMap)(t).Equals((*hashMap)(ot))
	}
	return false
}

func (t *exactMapType) HashCode() int {
	return (*hashMap)(t).HashCode()*31 + int(dgo.TiMapExact)
}

func (t *exactMapType) Instance(value interface{}) bool {
	if ov, ok := value.(*hashMap); ok {
		return (*hashMap)(t).Equals(ov)
	}
	return false
}

func (t *exactMapType) KeyType() dgo.Type {
	return &allOfValueType{slice: (*hashMap)(t).keys(), frozen: true}
}

func (t *exactMapType) Max() int {
	return t.len
}

func (t *exactMapType) Min() int {
	return t.len
}

func (t *exactMapType) String() string {
	return TypeString(t)
}

func (t *exactMapType) Type() dgo.Type {
	return &metaType{t}
}

func (t *exactMapType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMapExact
}

func (t *exactMapType) Value() dgo.Value {
	v := (*hashMap)(t)
	return v
}

func (t *exactMapType) ValueType() dgo.Type {
	return &allOfValueType{slice: (*hashMap)(t).values(), frozen: true}
}

func frozenMap(f string) error {
	return fmt.Errorf(`%s called on a frozen Map`, f)
}

func (t *exactMapType) Unbounded() bool {
	return false
}
