package internal_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/util"
	"github.com/tada/dgo/vf"
)

func ExampleMap_structType() {
	m := vf.Map(`host`, `example.com`, `port`, 22).Copy(false)
	fmt.Println(m)
	// Output: map[host:example.com port:22]
}

func ExampleMap_Put_structType() {
	m := vf.Map(`host`, `example.com`, `port`, 22).Copy(false)
	m.Put(`port`, 465)
	fmt.Println(m)
	// Output: map[host:example.com port:465]
}

func ExampleMap_Put_structTypeAdditionalKey() {
	m := vf.Map(`host`, `example.com`, `port`, 22).Copy(false)
	m.Put(`login`, `bob`)
	fmt.Println(m)
	// Output: map[host:example.com port:22 login:bob]
}

func ExampleMap_Put_structTypeIllegalKey() {
	type st struct {
		Host string
		Port int
	}
	m := vf.Map(&st{Host: `example.com`, Port: 22}).Copy(false)
	if err := catch.Do(func() { m.Put(`Login`, `bob`) }); err != nil {
		fmt.Println(err)
	}
	// Output: internal_test.st has no field named 'Login'
}

func TestMap_Put_structTypeIllegalValue(t *testing.T) {
	type st struct {
		Host string
		Port int
	}
	m := vf.Map(&st{Host: `example.com`, Port: 22}).Copy(false)
	assert.Panic(t, func() { m.Put(`Port`, `22`) }, `reflect: call of reflect.Value.SetString on int Value`)
}

func TestStructType_Get(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.StructMapType)
	assert.Equal(t, tp.GetEntryType(`a`).Value(), typ.Integer)
	assert.Nil(t, tp.GetEntryType(`c`))
}

func TestStructType_Validate(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.MapValidation)
	es := tp.Validate(nil, vf.Map(`a`, 1, `b`, `yes`))
	assert.Equal(t, 0, len(es))
}

func TestStructType_Validate_valueType(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.MapValidation)
	es := tp.Validate(nil, vf.Map(`a`, `no`, `b`, `yes`))
	assert.Equal(t, 1, len(es))
	assert.Equal(t, `parameter 'a' is not an instance of type int`, es[0].Error())
}

func TestStructType_Validate_missingKey(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.MapValidation)
	es := tp.Validate(nil, vf.Map(`a`, 1))
	assert.Equal(t, 1, len(es))
	assert.Equal(t, `missing required parameter 'b'`, es[0].Error())
}

func TestStructType_Validate_unknownKey(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.MapValidation)
	es := tp.Validate(nil, vf.Map(`a`, 1, `b`, `yes`, `c`, `no`))
	assert.Equal(t, 1, len(es))
	assert.Equal(t, `unknown parameter 'c'`, es[0].Error())
}

func TestStructType_Validate_notMap(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.MapValidation)
	es := tp.Validate(nil, vf.Values(1, 2))
	assert.Equal(t, 1, len(es))
	assert.Error(t, `value is not a Map`, es[0])
}

func TestStructType_ValidateVerbose_valueType(t *testing.T) {
	tp := tf.ParseType(`{a:int}`).(dgo.MapValidation)
	out := util.NewIndenter(`  `)
	ok := tp.ValidateVerbose(vf.Map(`a`, `no`), out)
	es := out.String()
	assert.False(t, ok)
	assert.Equal(t, `Validating 'a' against definition int
  'a' FAILED!
  Reason: expected a value of type int, got "no"
`, es)
}

func TestStructType_ValidateVerbose_missingKey(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.MapValidation)
	out := util.NewIndenter(`  `)
	ok := tp.ValidateVerbose(vf.Map(`a`, 1), out)
	es := out.String()
	assert.False(t, ok)
	assert.Equal(t, `Validating 'a' against definition int
  'a' OK!
Validating 'b' against definition string
  'b' FAILED!
  Reason: required key not found in input
`, es)
}

func TestStructType_ValidateVerbose_unknownKey(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.MapValidation)
	out := util.NewIndenter(`  `)
	ok := tp.ValidateVerbose(vf.Map(`a`, 1, `b`, `yes`, `c`, `no`), out)
	es := out.String()
	assert.False(t, ok)
	assert.Equal(t, `Validating 'a' against definition int
  'a' OK!
Validating 'b' against definition string
  'b' OK!
Validating 'c'
  'c' FAILED!
  Reason: key is not found in definition
`, es)
}

func TestStructType_ValidateVerbose(t *testing.T) {
	tp := tf.ParseType(`{a:int,b:string}`).(dgo.MapValidation)
	out := util.NewIndenter(``)
	assert.False(t, tp.ValidateVerbose(vf.Values(1, 2), out))
	assert.Equal(t, `value is not a Map`, out.String())
}

func TestStructType_alias(t *testing.T) {
	tp := tf.ParseType(`person={name:string,mom:person,dad:person}`).(dgo.StructMapType)
	assert.Same(t, tp, tp.GetEntryType(`mom`).Value())
}

func TestStructType(t *testing.T) {
	tp := tf.StructMap(false)
	assert.Equal(t, tp, tp)
	assert.Equal(t, tp, tf.ParseType(`{}`))
	assert.Equal(t, tp.KeyType(), typ.Any)
	assert.Equal(t, tp.ValueType(), typ.Any)
	assert.False(t, tp.Unbounded())

	tp = tf.StructMap(false,
		tf.StructMapEntry(`a`, typ.Integer, true))
	assert.Equal(t, tp, tf.ParseType(`{a:int}`))
	assert.NotEqual(t, tp, tf.ParseType(`{a?:int}`))
	assert.NotEqual(t, tp, tf.ParseType(`{a:int,b:int}`))
	assert.Assignable(t, tp, tf.ParseType(`{a:0..5}`))
	assert.Assignable(t, tp, tf.AnyOf(tf.ParseType(`{a:0..5}`), tf.ParseType(`{a:10..15}`)))
	assert.NotAssignable(t, tp, tf.ParseType(`{a:float}`))
	assert.Equal(t, tp.KeyType(), vf.String(`a`).Type())
	assert.Equal(t, tp.ValueType(), typ.Integer)
	assert.Equal(t, tp.Min(), 1)
	assert.Equal(t, tp.Max(), 1)

	tp = tf.StructMap(false,
		tf.StructMapEntry(`a`, typ.Integer, true),
		tf.StructMapEntry(`b`, typ.String, false))
	assert.Equal(t, tf.Map(typ.String, typ.Any), typ.Generic(tp))
	assert.Equal(t, tp, tp)
	assert.Equal(t, tp, tf.ParseType(`{a:int,b?:string}`))
	assert.Equal(t, tp.KeyType(), tf.ParseType(`"a"&"b"`))
	assert.Equal(t, tp.ValueType(), tf.ParseType(`int&string`))
	assert.NotEqual(t, tp, tf.ParseType(`map[string](int|string)`))
	assert.Equal(t, tp.Min(), 1)
	assert.Equal(t, tp.Max(), 2)
	assert.False(t, tp.Additional())

	m := vf.Map(`a`, 3, `b`, `yes`)
	assert.Assignable(t, tp, tp)
	assert.Assignable(t, tp, m.Type())
	assert.Instance(t, tp, m)

	m = vf.Map(`a`, 3)
	assert.Assignable(t, tp, m.Type())
	assert.Instance(t, tp, m)

	m = vf.Map(`b`, `yes`)
	assert.NotAssignable(t, tp, m.Type())
	assert.NotInstance(t, tp, m)

	m = vf.Map(`a`, 3, `b`, 4)
	assert.NotAssignable(t, tp, m.Type())
	assert.NotInstance(t, tp, m)
	assert.NotInstance(t, tp, vf.Values(`a`, `b`))
	assert.Instance(t, tp.Type(), tp)

	tps := tf.StructMap(false,
		tf.StructMapEntry(`a`, tf.Integer64(0, 10, true), true),
		tf.StructMapEntry(`b`, tf.String(20), false))
	assert.Equal(t, tps, tf.ParseType(`{a:0..10,b?:string[20]}`))
	assert.Assignable(t, tp, tps)

	tps = tf.StructMap(false,
		tf.StructMapEntry(`a`, typ.Integer, true),
		tf.StructMapEntry(`b`, typ.String, true))
	assert.Equal(t, tps, tf.ParseType(`{a:int,b:string}`))
	assert.Assignable(t, tp, tps)

	tps = tf.StructMap(false,
		tf.StructMapEntry(`a`, typ.Integer, false),
		tf.StructMapEntry(`b`, typ.String, false))
	assert.Equal(t, tps, tf.ParseType(`{a?:int,b?:string}`))
	assert.NotAssignable(t, tp, tps)

	tps = tf.StructMap(false,
		tf.StructMapEntry(`a`, typ.Integer, true))
	assert.Equal(t, tps, tf.ParseType(`{a:int}`))
	assert.Assignable(t, tp, tps)

	tps = tf.StructMap(false,
		tf.StructMapEntry(`a`, 3, true))
	assert.Equal(t, tps, tf.ParseType(`{a:3}`))
	assert.Assignable(t, tp, tps)

	tps = tf.StructMap(false,
		tf.StructMapEntry(`b`, typ.String, false))
	assert.Equal(t, tps, tf.ParseType(`{b?:string}`))
	assert.NotAssignable(t, tp, tps)

	tps = tf.StructMap(true,
		tf.StructMapEntry(`a`, typ.Integer, true))
	assert.Equal(t, tps, tf.ParseType(`{a:int,...}`))
	assert.NotAssignable(t, tp, tps)
	assert.Equal(t, tps.Min(), 1)
	assert.Equal(t, tps.Max(), dgo.UnboundedSize)
	assert.True(t, tps.Additional())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.NotEqual(t, tp.HashCode(), tf.ParseType(`{a:int,b?:string,...}`).HashCode())
	assert.Panic(t, func() {
		tf.StructMap(false,
			tf.StructMapEntry(tf.Pattern(regexp.MustCompile(`a*`)), typ.Integer, true))
	}, `non exact key types`)

	tps = tf.ParseType(`{a:0..10,b?:int}`).(dgo.StructMapType)
	assert.True(t, reflect.ValueOf(map[string]int64{}).Type().AssignableTo(tps.ReflectType()))
}

func TestStructEntry(t *testing.T) {
	tp := tf.StructMapEntry(`a`, typ.String, true)
	assert.Equal(t, tp, tf.StructMapEntry(`a`, typ.String, true))
	assert.NotEqual(t, tp, tf.StructMapEntry(`a`, typ.String, false))
	assert.NotEqual(t, tp, vf.Values(`a`, typ.String))
	assert.Equal(t, `"a":string`, tp.String())
}

func TestStructFromMap(t *testing.T) {
	assert.Panic(t, func() {
		tf.StructMapFromMap(false, vf.Map(3, `dope`))
	}, `cannot be assigned to a variable of type map`)

	tp := tf.StructMapFromMap(false, vf.Map(`first`, typ.String))
	assert.Equal(t, tp, tf.StructMap(false, tf.StructMapEntry(`first`, typ.String, true)))

	tp = tf.StructMapFromMap(false, vf.Map(`first`, vf.Map(`type`, typ.String)))
	assert.Equal(t, tp, tf.StructMap(false, tf.StructMapEntry(`first`, typ.String, true)))

	tp = tf.StructMapFromMap(false, vf.Map(`first`, vf.Map(`type`, typ.String, `required`, true)))
	assert.Equal(t, tp, tf.StructMap(false, tf.StructMapEntry(`first`, typ.String, true)))

	tp = tf.StructMapFromMap(false, vf.Map(`first`, vf.Map(`type`, typ.String, `required`, false)))
	assert.Equal(t, tp, tf.StructMap(false, tf.StructMapEntry(`first`, typ.String, false)))

	tp = tf.StructMapFromMap(false, vf.Map(`first`, vf.Map(`type`, `string`, `required`, false)))
	assert.Equal(t, tp, tf.StructMap(false, tf.StructMapEntry(`first`, typ.String, false)))

	tp = tf.StructMapFromMap(false, vf.Map(`first`, `string`))
	assert.Equal(t, tp, tf.StructMap(false, tf.StructMapEntry(`first`, typ.String, true)))

	tp = tf.StructMapFromMap(false, vf.Map(`first`, `"x"`))
	assert.Equal(t, tp, tf.StructMap(false, tf.StructMapEntry(`first`, vf.String(`x`).Type(), true)))

	tp = tf.StructMapFromMap(false, vf.Map())
	assert.Equal(t, tp, tf.StructMap(false))
	assert.False(t, tp.Unbounded())

	tp = tf.StructMapFromMap(true, vf.Map())
	assert.Equal(t, typ.Any, tp.KeyType())
	assert.Equal(t, typ.Any, tp.ValueType())
	assert.True(t, tp.Unbounded())
}
