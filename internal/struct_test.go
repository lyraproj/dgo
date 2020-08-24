package internal_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func TestStruct(t *testing.T) {
	type structA struct {
		A int
		B int
	}
	m := vf.Map(&structA{1, 2})
	tp := m.(dgo.StructMapType)
	assert.Same(t, m, tp)
	assert.Assignable(t, tp, m.Type())
	assert.False(t, tp.Additional())
	assert.Instance(t, tp, m)
	assert.Equal(t, tf.Map(typ.String, typ.Integer), typ.Generic(tp))
	assert.Equal(t, 2, tp.Max())
	assert.Equal(t, 2, tp.Min())
	assert.False(t, tp.Unbounded())
	assert.Equal(t, tf.StructMapEntry(`A`, 1, true), tp.GetEntryType(`A`))

	m2 := tp.(dgo.Factory).New(vf.Map(`A`, 1, `B`, 2))
	assert.Equal(t, m, m2)
	assert.Equal(t, &structA{1, 2}, m2.(dgo.Struct).GoStruct())
	assert.Equal(t, `internal_test.structA`, tp.ReflectType().String())

	c := 0
	tp.EachEntryType(func(e dgo.StructMapEntry) {
		c++
		assert.Equal(t, c, e.Value())
	})
	assert.Equal(t, 2, c)
}

func TestStruct_Any(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m := vf.Map(&structA{1, 2.0, `three`})
	assert.False(t, m.Any(func(e dgo.MapEntry) bool {
		return e.Key().Equals(`Fourth`)
	}))
	assert.True(t, m.Any(func(e dgo.MapEntry) bool {
		return e.Key().Equals(`Second`) && e.Value().Equals(2.0)
	}))
}

func TestStruct_AllKeys(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m := vf.Map(&structA{1, 2.0, `three`})
	assert.False(t, m.AllKeys(func(k dgo.Value) bool {
		return len(k.String()) == 5
	}))
	assert.True(t, m.AnyKey(func(k dgo.Value) bool {
		return len(k.String()) >= 5
	}))
}

func TestStruct_AnyKey(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m := vf.Map(&structA{1, 2.0, `three`})
	assert.False(t, m.AnyKey(func(k dgo.Value) bool {
		return k.Equals(`Fourth`)
	}))
	assert.True(t, m.AnyKey(func(k dgo.Value) bool {
		return k.Equals(`Second`)
	}))
}

func TestStruct_AnyValue(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m := vf.Map(&structA{1, 2.0, `three`})
	assert.False(t, m.AnyValue(func(v dgo.Value) bool {
		return v.Equals(`four`)
	}))
	assert.True(t, m.AnyValue(func(v dgo.Value) bool {
		return v.Equals(`three`)
	}))
}

func TestStruct_ContainsKey(t *testing.T) {
	type structA struct {
		First  int
		Second float64
	}
	m := vf.Map(&structA{})
	assert.True(t, m.ContainsKey(`First`))
	assert.False(t, m.ContainsKey(`Third`))
	assert.False(t, m.ContainsKey(1))
}

func TestStruct_Copy(t *testing.T) {
	type structA struct {
		A string
		E dgo.Array
	}
	s := structA{A: `Alpha`}
	m := vf.Map(&s)
	m.Put(`E`, vf.MutableValues(`Echo`, `Foxtrot`))

	c := m.Copy(true).(dgo.Map)
	assert.False(t, m.Frozen())
	assert.False(t, m.Get(`E`).(dgo.Freezable).Frozen())
	assert.True(t, c.Frozen())
	assert.True(t, c.Get(`E`).(dgo.Freezable).Frozen())

	m.Put(`A`, `Adam`)
	assert.Equal(t, `Adam`, m.Get(`A`))
	assert.Equal(t, `Alpha`, c.Get(`A`))
	assert.Same(t, c, c.Copy(true))

	assert.Panic(t, func() { c.Put(`A`, `Adam`) }, `frozen`)

	d := c.Copy(false).(dgo.Map)
	assert.NotSame(t, c, d)
	d.Put(`A`, `Adam`)
	assert.Equal(t, `Adam`, d.Get(`A`))
	assert.Equal(t, `Alpha`, c.Get(`A`))
}

func TestStruct_EachKey(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m := vf.Map(&structA{1, 2.0, `three`})
	var vs []dgo.Value
	m.EachKey(func(v dgo.Value) {
		vs = append(vs, v)
	})

	assert.Equal(t, 3, len(vs))
	assert.Equal(t, vf.Values(`First`, `Second`, `Third`), vs)
}

func TestStruct_EachValue(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m := vf.Map(&structA{1, 2.0, `three`})
	var vs []dgo.Value
	m.EachValue(func(v dgo.Value) {
		vs = append(vs, v)
	})

	assert.Equal(t, 3, len(vs))
	assert.Equal(t, vf.Values(1, 2.0, `three`), vs)
}

func TestStruct_Each(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m := vf.Map(&structA{1, 2.0, `three`})
	var vs []dgo.Value
	m.Each(func(v dgo.Value) {
		e := v.(dgo.MapEntry)
		vs = append(vs, e.Key())
		vs = append(vs, e.Value())
	})

	assert.Equal(t, 6, len(vs))
	assert.Equal(t, vf.Values(`First`, 1, `Second`, 2.0, `Third`, `three`), vs)
}

func TestStruct_Get(t *testing.T) {
	type structA struct {
		A string
		B int
		C *string
		D *int
		E []string
	}
	c := `Charlie`
	s := structA{A: `Alpha`, B: 32, C: &c, E: []string{`Echo`, `Foxtrot`}}
	m := vf.Map(&s)
	assert.Equal(t, m.Get(`A`), `Alpha`)
	assert.Equal(t, m.Get(`B`), 32)
	assert.Equal(t, m.Get(`C`), c)
	assert.Equal(t, m.Get(`D`), vf.Nil)
	assert.Equal(t, m.Get(`E`), []string{`Echo`, `Foxtrot`})
	assert.Equal(t, m.Get(`F`), nil)
	assert.Equal(t, m.Get(10), nil)
}

func TestStruct_Find(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	s := structA{A: `Alpha`, B: 32}
	m := vf.Map(&s)
	found := m.Find(func(e dgo.MapEntry) bool {
		return e.Key().Equals(`B`)
	})
	assert.Equal(t, found.Value(), 32)
	found = m.Find(func(e dgo.MapEntry) bool { return false })
	assert.Nil(t, found)
}

func TestStruct_Freeze(t *testing.T) {
	type structA struct {
		A string
	}
	s := structA{}
	m := vf.Map(&s)
	m.Freeze()
	assert.Panic(t, func() { m.Put(`A`, `Alpha`) }, `frozen`)
}

func TestStruct_FrozenCopy(t *testing.T) {
	type structA struct {
		A string
		E dgo.Array
	}
	s := structA{A: `Alpha`}
	m := vf.Map(&s)
	m.Put(`E`, vf.MutableValues(`Echo`, `Foxtrot`))

	c := m.FrozenCopy().(dgo.Map)
	assert.False(t, m.Frozen())
	assert.False(t, m.Get(`E`).(dgo.Freezable).Frozen())
	assert.True(t, c.Frozen())
	assert.True(t, c.Get(`E`).(dgo.Freezable).Frozen())

	m.Put(`A`, `Adam`)
	assert.Equal(t, `Adam`, m.Get(`A`))
	assert.Equal(t, `Alpha`, c.Get(`A`))
	assert.Same(t, c, c.FrozenCopy())

	assert.Panic(t, func() { c.Put(`A`, `Adam`) }, `frozen`)
}

func TestStruct_ThawedCopy(t *testing.T) {
	type structA struct {
		A string
		E dgo.Array
	}
	s := structA{A: `Alpha`}
	m := vf.Map(&s)
	m.Put(`E`, vf.MutableValues(`Echo`, `Foxtrot`))
	m.Freeze()

	c := m.ThawedCopy().(dgo.Map)
	assert.True(t, m.Frozen())
	assert.True(t, m.Get(`E`).(dgo.Freezable).Frozen())
	assert.False(t, c.Frozen())
	assert.False(t, c.Get(`E`).(dgo.Freezable).Frozen())

	c.Put(`A`, `Adam`)
	assert.Equal(t, `Adam`, c.Get(`A`))
	assert.Equal(t, `Alpha`, m.Get(`A`))
	assert.NotSame(t, c, c.FrozenCopy())
}

func TestStruct_HashCode(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	s := structA{A: `Alpha`, B: 32}
	m := vf.Map(&s)
	assert.NotEqual(t, 0, m.HashCode())
	assert.Equal(t, m.HashCode(), m.HashCode())
}

func TestStruct_GoStruct(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	s := structA{A: `Alpha`, B: 32}
	m := vf.Map(&s).(dgo.Struct)
	_, ok := m.GoStruct().(*structA)
	require.True(t, ok)
}

func TestStruct_Keys(t *testing.T) {
	type structA struct {
		A string
		B int
		C *string
	}
	s := structA{}
	m := vf.Map(&s)
	assert.Equal(t, vf.Strings(`A`, `B`, `C`), m.Keys())
}

func TestStruct_Map(t *testing.T) {
	type structA struct {
		A string
		B string
		C string
	}
	a := vf.Map(&structA{A: `value a`, B: `value b`, C: `value c`})
	assert.Equal(t,
		vf.Map(map[string]string{`A`: `the a`, `B`: `the b`, `C`: `the c`}),
		a.Map(func(e dgo.MapEntry) interface{} {
			return strings.Replace(fmt.Sprintf("%v", e.Value()), `value`, `the`, 1)
		}))
	assert.Equal(t, vf.Map(`A`, nil, `B`, vf.Nil, `C`, nil), a.Map(func(e dgo.MapEntry) interface{} {
		return nil
	}))
}

func TestStruct_Merge(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m1 := vf.Map(&structA{1, 2.0, `three`})

	m2 := vf.Map(
		`Third`, `tres`,
		`Fourth`, `cuatro`)

	assert.Equal(t, m1.Merge(m2), vf.Map(
		`First`, 1,
		`Second`, 2.0,
		`Third`, `tres`,
		`Fourth`, `cuatro`))

	assert.Same(t, m1, m1.Merge(m1))
	assert.Same(t, m1, m1.Merge(vf.Map()))
	assert.Same(t, m1, vf.Map().Merge(m1))
}

func TestStruct_Put(t *testing.T) {
	type structA struct {
		A string
		B int
		C *string
		D *int
		E []string
	}
	s := structA{}
	m := vf.Map(&s)
	m.Put(`A`, `Alpha`)
	m.Put(`B`, 32)
	m.Put(`C`, `Charlie`)
	m.Put(`D`, 42)
	m.Put(`E`, []string{`Echo`, `Foxtrot`})

	assert.Panic(t, func() { m.Put(`F`, `nope`) }, `no field named 'F'`)

	assert.Equal(t, s.A, `Alpha`)
	assert.Equal(t, s.B, 32)
	assert.Equal(t, *s.C, `Charlie`)
	assert.Equal(t, *s.D, 42)
	assert.Equal(t, s.E, []string{`Echo`, `Foxtrot`})
}

func TestStruct_PutAll(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	s := structA{}
	m := vf.Map(&s)
	m.PutAll(vf.Map(`A`, `Alpha`, `B`, 32))
	assert.Equal(t, s.A, `Alpha`)
	assert.Equal(t, s.B, 32)
}

func TestStruct_ReflectTo(t *testing.T) {
	type structA struct {
		A string
		B int
		C *string
		D *int
		E []string
	}

	type structB struct {
		A string
		B *structA
		C structA
	}

	c := `Charlie`
	d := 42
	s := structA{A: `Alpha`, B: 32, C: &c, D: &d, E: []string{`Echo`, `Foxtrot`}}

	// Pass pointer to struct
	m := vf.Map(&s)

	x := structB{}
	rv := reflect.ValueOf(&x).Elem()

	// By pointer
	xb := rv.FieldByName(`B`)
	m.ReflectTo(xb)
	assert.Equal(t, x.B, &s)
	assert.Same(t, x.B, &s)

	// By value
	m.ReflectTo(rv.FieldByName(`C`))
	assert.Equal(t, &x.C, &s)
	assert.NotSame(t, &x.C, &s)

	m.Freeze()
	m.ReflectTo(xb)
	assert.NotSame(t, &x.B, &s)
}

func TestStruct_Remove(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	s := structA{}
	m := vf.Map(&s)
	assert.Panic(t, func() { m.Remove(`B`) }, `cannot be removed`)
	assert.Panic(t, func() { m.RemoveAll(vf.Values(`A`, `B`)) }, `cannot be removed`)
}

func TestStruct_AppendTo(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	v := vf.Map(&structA{A: `hello`, B: 2})
	assert.Equal(t, `{
  "A": "hello",
  "B": 2
}`, util.ToIndentedStringERP(v))
}

func TestStruct_String(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	s := structA{A: `Alpha`, B: 32}
	m := vf.Map(&s)
	assert.Equal(t, `{"A":"Alpha","B":32}`, m.String())
}

func TestStruct_StringKeys(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	s := structA{A: `Alpha`, B: 32}
	m := vf.Map(&s)
	assert.True(t, m.StringKeys())
}

func TestStruct_Type(t *testing.T) {
	type structA struct {
		A string
		B int
	}
	s := structA{A: `Alpha`, B: 32}
	m := vf.Map(&s)
	tp := m.Type()
	assert.Assignable(t, tp, tp)
	assert.Instance(t, tp, m)
}

func TestStruct_Validate(t *testing.T) {
	type structA struct {
		A int
		B int
	}
	s := structA{A: 3, B: 4}
	tp := vf.Map(&s).Type().(dgo.MapValidation)
	es := tp.Validate(nil, vf.Map(`A`, 3, `B`, 4))
	assert.Equal(t, 0, len(es))

	es = tp.Validate(nil, vf.Map(`A`, 2, `B`, 4))
	assert.Equal(t, 1, len(es))
}

func TestStruct_ValidateVerbose(t *testing.T) {
	type structA struct {
		A int
		B int
	}
	s := structA{A: 3, B: 4}
	tp := vf.Map(&s).Type().(dgo.MapValidation)
	out := util.NewIndenter(``)
	assert.False(t, tp.ValidateVerbose(vf.Values(1, 2), out))
	assert.Equal(t, `value is not a Map`, out.String())
}

func TestStruct_Values(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	m := vf.Map(&structA{1, 2.0, `three`})

	assert.True(t, m.Values().SameValues(vf.Values(1, 2.0, `three`)))
}

func TestStruct_With(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	om := vf.Map(&structA{1, 2.0, `three`})
	m := om.With(`Fourth`, `quad`)
	assert.Equal(t, m, vf.Map(`First`, 1, `Second`, 2.0, `Third`, `three`, `Fourth`, `quad`))
}

func TestStruct_Without(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	om := vf.Map(&structA{1, 2.0, `three`})
	m := om.Without(`Second`)
	assert.Equal(t, m, map[string]interface{}{
		`First`: 1,
		`Third`: `three`,
	})

	// Original is not modified
	assert.Equal(t, om, map[string]interface{}{
		`First`:  1,
		`Second`: 2.0,
		`Third`:  `three`,
	})

	m = om.Without(`NotPresent`)
	assert.Same(t, m, om)
}

func TestStruct_WithoutAll(t *testing.T) {
	type structA struct {
		First  int
		Second float64
		Third  string
	}
	om := vf.Map(&structA{1, 2.0, `three`})
	m := om.WithoutAll(vf.Strings(`First`, `Second`))
	assert.Equal(t, m, map[string]interface{}{
		`Third`: `three`,
	})

	// Original is not modified
	assert.Equal(t, om, map[string]interface{}{
		`First`:  1,
		`Second`: 2.0,
		`Third`:  `three`,
	})
}
