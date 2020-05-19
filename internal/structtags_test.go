package internal_test

import (
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/vf"
)

func TestStructTags_intRange(t *testing.T) {
	type structA struct {
		A int `dgo:"0..10"`
	}
	assert.Equal(t, tf.Parse(`{"A":0..10}`), tf.StructMapTypeFromReflected(reflect.TypeOf(structA{})))

	s := &structA{}
	m := vf.MutableMap(s)
	m.Put(`A`, 4)
	assert.Equal(t, 4, s.A)
	assert.Panic(t, func() { m.Put(`A`, 11) }, `the value 11 cannot be assigned`)
}

func TestStructTags_tagInNested(t *testing.T) {
	type structA struct {
		A int `dgo:"0..10"`
	}
	type structB struct {
		B int `dgo:"20..30"`
	}
	type structC struct {
		structA
		structB
		C map[string]string `dgo:"map[string]/^[a-z].*/"`
	}
	assert.Equal(t, tf.Parse(`{"A":0..10,"B":20..30,"C":map[string]/^[a-z].*/}`),
		tf.StructMapTypeFromReflected(reflect.TypeOf(structC{})))

	s := &structC{}
	m := vf.MutableMap(s)
	m.Put(`A`, 4)
	m.Put(`B`, 24)
	m.Put(`C`, map[string]string{`x`: `text`})
	assert.Equal(t, 4, s.A)
	assert.Equal(t, 24, s.B)
	assert.Equal(t, map[string]string{`x`: `text`}, s.C)
	assert.Panic(t, func() { m.Put(`A`, 11) }, `the value 11 cannot be assigned`)
	assert.Panic(t, func() { m.Put(`B`, 31) }, `the value 31 cannot be assigned`)
	assert.Panic(t, func() { m.Put(`C`, map[string]string{`x`: `Text`}) },
		`the value map\[x:Text\] cannot be assigned to a variable of type map\[string\]/\^\[a-z\]\.\*/`)
	assert.Panic(t, func() { m.Put(23, ``) }, `internal_test.structC has no field named '23'`)
}

func TestStructTags_typeMismatch(t *testing.T) {
	type structA struct {
		A int `dgo:"float"`
	}
	assert.Panic(t, func() { tf.StructMapTypeFromReflected(reflect.TypeOf(structA{})) },
		`type float declared in field tag of internal_test.structA.A is not assignable to actual type int`)
}

func TestStructTags_badSyntax(t *testing.T) {
	type structA struct {
		A int `dgo:"1+2"`
	}
	assert.Panic(t, func() { tf.StructMapTypeFromReflected(reflect.TypeOf(structA{})) },
		`unable to parse dgo tag structA.A "1\+2": `)
}
