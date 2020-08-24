package stringer_test

import (
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestType_String_floatExact(t *testing.T) {
	assert.Equal(t, `1.0`, vf.Value(1.0).Type().String())
}

func TestType_String_stringExact(t *testing.T) {
	assert.Equal(t, `"with \"quotes\" in it"`, vf.Value("with \"quotes\" in it").Type().String())
}

func TestType_String_arrayExact(t *testing.T) {
	assert.Equal(t, `{1.0,2.3}`, vf.Values(1.0, 2.3).Type().String())
}

func TestType_String_arrayUnboundedElementInt(t *testing.T) {
	assert.Equal(t, `[]int`, tf.Array(typ.Integer).String())
}

func TestType_String_arrayBoundedElementInt(t *testing.T) {
	assert.Equal(t, `[0,4]int`, tf.Array(typ.Integer, 0, 4).String())
}

func TestType_String_arrayBoundedElementRange(t *testing.T) {
	assert.Equal(t, `[0,4]3..8`, tf.Array(tf.Integer64(3, 8, true), 0, 4).String())
}

func TestType_String_arrayBoundedElementEnum(t *testing.T) {
	assert.Equal(t, `[0,4]("a"|"b")`, tf.Array(tf.Enum(`a`, `b`), 0, 4).String())
}

func TestType_String_mapUnboundedEntryStringInt(t *testing.T) {
	assert.Equal(t, `map[string]int`, tf.Map(typ.String, typ.Integer).String())
}

func TestType_String_arrayBoundedEntryStringInt(t *testing.T) {
	assert.Equal(t, `map[string,0,4]int`, tf.Map(typ.String, typ.Integer, 0, 4).String())
}

func TestType_String_arrayBoundedEntryStringRange(t *testing.T) {
	assert.Equal(t, `map[string,0,4]3..8`, tf.Map(typ.String, tf.Integer64(3, 8, true), 0, 4).String())
}

func TestType_String_arrayBoundedEntryStringEnum(t *testing.T) {
	assert.Equal(t, `map[string,0,4]("a"|"b")`, tf.Map(typ.String, tf.Enum(`a`, `b`), 0, 4).String())
}

func TestType_String_arrayBoundedEntryIntRangeEnum(t *testing.T) {
	assert.Equal(t, `map[3..8,0,4]("a"|"b")`, tf.Map(tf.Integer64(3, 8, true), tf.Enum(`a`, `b`), 0, 4).String())
}

func TestType_String_priorities(t *testing.T) {
	var tp dgo.Type = tf.Array(tf.String(1), 2, 2)

	assert.Equal(t, `[2,2]string[1]`, tp.String())

	tp = tf.Array(vf.Value(regexp.MustCompile(`a`)).Type(), 2, 2)
	assert.Equal(t, `[2,2]regexp "a"`, tp.String())

	tp = tf.Array(tf.Not(tf.String(1)), 2, 2)
	assert.Equal(t, `[2,2]!string[1]`, tp.String())

	tp = tf.Not(tf.String(1))
	assert.Equal(t, `!string[1]`, tp.String())

	tp = tf.Array(tf.Map(typ.Integer, typ.String, 1), 2, 2)
	assert.Equal(t, `[2,2]map[int,1]string`, tp.String())

	tp = tf.Array(tf.Map(typ.Integer, tf.String(1), 2), 3, 4)
	assert.Equal(t, `[3,4]map[int,2]string[1]`, tp.String())

	tp = tf.Array(tf.AllOf(tf.Enum(`a`, `b`), tf.Enum(`b`, `c`)), 2, 2)
	assert.Equal(t, `[2,2](("a"|"b")&("b"|"c"))`, tp.String())

	tp = tf.Array(tf.OneOf(typ.Integer, typ.String), 2, 2)
	assert.Equal(t, `[2,2](int^string)`, tp.String())

	tp = tf.Array(vf.Values(`a`, `b`).Type().(dgo.ArrayType).ElementType(), 2, 2)
	assert.Equal(t, `[2,2]("a"&"b")`, tp.String())
}

func TestTypeIdentifier_String(t *testing.T) {
	assert.Equal(t, `pattern`, dgo.TiStringPattern.String())

	assert.Panic(t, func() { _ = dgo.TypeIdentifier(0x1000).String() }, `unhandled TypeIdentifier 4096`)
}
