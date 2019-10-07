package internal_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
	"gopkg.in/yaml.v3"
)

type failingMarshaler struct {
	dgo.Value
}

var errFailing = errors.New("errFailing")

func (ft *failingMarshaler) MarshalYAML() (interface{}, error) {
	return nil, errFailing
}

func TestMap_MarshalYaml(t *testing.T) {
	m := vf.Map("a", 1, "b", "two", "c", vf.Values(`hello`, true, 1, 3.14, nil))
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `a: 1
b: two
c:
  - hello
  - true
  - 1
  - 3.14
  - null
`, string(b))
}

func TestMap_MarshalYaml_binary(t *testing.T) {
	m := vf.Map("b", vf.BinaryFromString(`AQQD`))
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `b: !!binary AQQD
`, string(b))
}

func TestMap_MarshalYaml_timestamp(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	m := vf.Map("t", vf.Time(ts))
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `t: !!timestamp 2019-10-06T07:15:00-07:00
`, string(b))
}

func TestMap_MarshalYaml_type(t *testing.T) {
	m := vf.Map("t", typ.String)
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `t: !puppet.com,2019:dgo/type string
`, string(b))
}

func TestMap_MarshalYaml_fail(t *testing.T) {
	m := vf.MutableMap(nil)
	m.Put("a", 1)
	m.Put("b", vf.MutableValues(nil, &failingMarshaler{vf.String(`bv`)}))
	_, err := yaml.Marshal(m)
	if err == nil {
		t.Fatalf(`no error was returned`)
	}
	require.Equal(t, errFailing.Error(), err.Error())

	m = vf.MutableMap(nil)
	m.Put(&failingMarshaler{vf.String(`b`)}, "bv")
	_, err = yaml.Marshal(m)
	if err == nil {
		t.Fatalf(`no error was returned`)
	}
	require.Equal(t, errFailing.Error(), err.Error())
}

func TestMap_UnmarshalYAML_Map(t *testing.T) {
	m := vf.MutableMap(nil)
	require.Ok(t, yaml.Unmarshal([]byte(`
a: 1
b: two
c: 
  - hello
  - true
  - 1
  - 3.14
  - null
`), m))
	require.Equal(t, vf.Map("a", 1, "b", "two", "c", vf.Values(`hello`, true, 1, 3.14, nil)), m)
}

func TestMap_UnmarshalYAML_binary(t *testing.T) {
	m := vf.MutableMap(nil)
	require.Ok(t, yaml.Unmarshal([]byte("b: !!binary AQQD\n"), m))
	require.Equal(t, vf.Map("b", vf.BinaryFromString(`AQQD`)), m)
}

func TestMap_UnmarshalYAML_timestamp(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	m := vf.MutableMap(nil)
	require.Ok(t, yaml.Unmarshal([]byte("t: !!timestamp 2019-10-06T07:15:00-07:00\n"), m))
	require.Equal(t, vf.Map("t", vf.Time(ts)), m)
}

func TestMap_UnmarshalYAML_bad_timestamp(t *testing.T) {
	m := vf.MutableMap(nil)
	require.Panic(t, func() { _ = yaml.Unmarshal([]byte("t: !!timestamp 2019-13-06T07:15:00-07:00\n"), m) }, `cannot decode`)
}

func TestMap_UnmarshalYAML_type(t *testing.T) {
	m := vf.MutableMap(nil)
	require.Ok(t, yaml.Unmarshal([]byte("t: !puppet.com,2019:dgo/type string\n"), m))
	require.Equal(t, vf.Map("t", typ.String), m)
}

func TestMap_UnmarshalYAML_TypedMap(t *testing.T) {
	m := vf.MutableMap(`map[string](int|float)`)
	require.Ok(t, yaml.Unmarshal([]byte(`
int: 23
float: 3.14`), m))

	require.Panic(t, func() {
		m := vf.MutableMap(`map[string](int|float)`)
		_ = yaml.Unmarshal([]byte(`
int: 23
string: hello`), m)
	}, `the string "hello" cannot be assigned to a variable of type int|float`)
}

type testNoMarshaler struct {
	A string
}

type testMarshaler struct {
	A string
}

type marshalTestNode struct {
	testMarshaler
}

type marshalTestFail struct {
	testMarshaler
}

func (m *testMarshaler) MarshalYAML() (interface{}, error) {
	return internal.Map([]interface{}{`A`, m.A}), nil
}

func (m *marshalTestNode) MarshalYAML() (interface{}, error) {
	return internal.Map([]interface{}{`A`, m.A}).MarshalYAML()
}

func (m *marshalTestFail) MarshalYAML() (interface{}, error) {
	return nil, errFailing
}

func TestNative_MarshalYAML(t *testing.T) {
	m := vf.MutableValues(nil, vf.Value(&testMarshaler{A: `hello`}))
	bs, err := yaml.Marshal(m)
	require.Ok(t, err)
	require.Equal(t, "- A: hello\n", string(bs))
}

func TestNative_MarshalYAML_node(t *testing.T) {
	m := vf.MutableValues(nil, vf.Value(&marshalTestNode{testMarshaler{A: `hello`}}))
	bs, err := yaml.Marshal(m)
	require.Ok(t, err)
	require.Equal(t, "- A: hello\n", string(bs))
}

func TestNative_MarshalYAML_failNoMarshaler(t *testing.T) {
	m := vf.MutableValues(nil, vf.Value(&testNoMarshaler{}))
	_, err := yaml.Marshal(m)
	require.NotNil(t, err)
	require.Equal(t, `unable to marshal into value of type *internal_test.testNoMarshaler`, err.Error())
}

func TestNative_MarshalYAML_fail(t *testing.T) {
	m := vf.MutableValues(nil, vf.Value(&marshalTestFail{}))
	_, err := yaml.Marshal(m)
	require.NotNil(t, err)
	require.Equal(t, `errFailing`, err.Error())
}

func TestMap_UnmarshalYAML_Map_wrongType(t *testing.T) {
	m := vf.MutableMap(nil)
	err := yaml.Unmarshal([]byte(`["hello",true,1,3.14,null]`), m)
	if err == nil {
		t.Fatalf(`no error was returned`)
	}
	require.Equal(t, `expecting data to be a map`, err.Error())
}

func TestMap_UnmarshalYAML_Map_frozen(t *testing.T) {
	m := vf.Map(`a`, 1)
	require.Panic(t, func() { _ = yaml.Unmarshal([]byte(`{"b":"two"}`), m) }, `UnmarshalYAML .* frozen`)
}

func TestMap_UnmarshalYAML_Array(t *testing.T) {
	a := vf.MutableValues(nil)
	err := yaml.Unmarshal([]byte(`["hello",true,1,3.14,null]`), a)
	require.Ok(t, err)
	require.Equal(t, vf.Values(`hello`, true, 1, 3.14, nil), a)
}

func TestMap_UnmarshalYAML_Array_wrongType(t *testing.T) {
	a := vf.MutableValues(nil)
	err := yaml.Unmarshal([]byte(`{"a":1,"b":"two"}`), a)
	if err == nil {
		t.Fatalf(`no error was returned`)
	}
	require.Equal(t, `expecting data to be an array`, err.Error())
}

func TestMap_UnmarshalYAML_Array_frozen(t *testing.T) {
	a := vf.Values()
	require.Panic(t, func() { _ = yaml.Unmarshal([]byte(`["hello",true,1,3.14,null]`), a) }, `UnmarshalYAML .* frozen`)
}
