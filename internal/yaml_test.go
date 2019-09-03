package internal_test

import (
	"errors"
	"testing"

	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"

	require "github.com/lyraproj/dgo/dgo_test"
	"gopkg.in/yaml.v3"
)

type failingMarshaler struct{}

var failingErr = errors.New("failingErr")

func (ft *failingMarshaler) MarshalYAML() (interface{}, error) {
	return nil, failingErr
}

func TestMap_MarshalYaml(t *testing.T) {
	m := vf.Map(map[string]interface{}{"a": 1, "b": "two", "c": vf.Values(`hello`, true, 1, 3.14, typ.String, nil)})
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `a: 1
b: two
c:
  - hello
  - true
  - 1
  - 3.14
  - !puppet.com,2019:dgo/type string
  - null
`, string(b))
}

func TestMap_MarshalYaml_fail(t *testing.T) {
	m := vf.MutableMap(2, nil)
	m.Put("a", 1)
	m.Put("b", vf.MutableValues(nil, &failingMarshaler{}))
	_, err := yaml.Marshal(m)
	if err == nil {
		t.Fatalf(`no error was returned`)
	}
	require.Equal(t, failingErr.Error(), err.Error())
}

func TestMap_UnmarshalYAML_Map(t *testing.T) {
	m := vf.MutableMap(7, nil)
	err := yaml.Unmarshal([]byte(`{"a":1,"b":"two","c":["hello",true,1,3.14,null]}`), m)
	require.Nil(t, err)
	require.Equal(t, vf.Map(map[string]interface{}{"a": 1, "b": "two", "c": vf.Values(`hello`, true, 1, 3.14, nil)}), m)
}

func TestMap_UnmarshalYAML_Array(t *testing.T) {
	a := vf.MutableValues(nil)
	err := yaml.Unmarshal([]byte(`["hello",true,1,3.14,null]`), a)
	require.Nil(t, err)
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
