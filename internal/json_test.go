package internal_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func TestUnmarshalJSON(t *testing.T) {
	v, err := vf.UnmarshalJSON([]byte(`["hello",true,1,3.14,null,{"a":1}]`))
	require.Ok(t, err)
	require.Equal(t, vf.Values(`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1}), v)
}

func TestArray_MarshalJSON(t *testing.T) {
	a := vf.Values(
		`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1, "b": "two", "c": 2.17})
	b, err := a.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, `["hello",true,1,3.14,null,{"a":1,"b":"two","c":2.17}]`, string(b))
}

func TestNil_MarshalJSON(t *testing.T) {
	b, err := vf.Nil.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, `null`, string(b))
}

func TestString_MarshalJSON(t *testing.T) {
	b, err := vf.String(`with"in it`).(json.Marshaler).MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, `"with\"in it"`, string(b))
}

func TestArray_UnmarshalJSON(t *testing.T) {
	a := vf.MutableValues(nil)
	err := a.UnmarshalJSON([]byte(`["hello",true,1,3.14,null,{"a":1,"b":"two","c":2.17}]`))
	require.Nil(t, err)
	require.Equal(t, a, vf.Values(`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1, "b": "two", "c": 2.17}))
}

func TestArray_UnmarshalJSON_frozen(t *testing.T) {
	a := vf.Values()
	require.Panic(t, func() { _ = a.UnmarshalJSON([]byte(`["hello",true,1,3.14,null,{"a":1,"b":"two","c":2.17}]`)) }, `frozen`)
}

func TestMap_UnmarshalJSON_notArray(t *testing.T) {
	a := vf.MutableValues(nil)
	require.Equal(t, errors.New(`expecting data to be an array`), a.UnmarshalJSON([]byte(`{"a":1,"b":"two","c":2.17}`)))
}

func TestMap_UnmarshalJSON_unbalancedArray(t *testing.T) {
	a := vf.MutableValues(nil)
	require.Equal(t, errors.New(`invalid character '}' after array element`), a.UnmarshalJSON([]byte(`["hello",true}`)))
}

func TestMap_MarshalJSON(t *testing.T) {
	m := vf.Map("a", 1, "b", "two", "c", vf.Values(`hello`, true, 1, 3.14, nil))
	b, err := m.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, `{"a":1,"b":"two","c":["hello",true,1,3.14,null]}`, string(b))
}

func TestMap_UnmarshalJSON(t *testing.T) {
	m := vf.MutableMap(nil)
	err := m.UnmarshalJSON([]byte(`{"a":1,"b":"two","c":["hello",true,1,3.14,null]}`))
	require.Nil(t, err)
	require.Equal(t, vf.Map("a", 1, "b", "two", "c", vf.Values(`hello`, true, 1, 3.14, nil)), m)
}

func TestMap_UnmarshalJSON_frozen(t *testing.T) {
	m := vf.Map(map[string]int{})
	require.Panic(t, func() { _ = m.UnmarshalJSON([]byte(`{"a":1,"b":"two","c":["hello",true,1,3.14,null]}`)) }, `frozen`)
}

func TestMap_UnmarshalJSON_notObject(t *testing.T) {
	m := vf.MutableMap(nil)
	require.Equal(t, errors.New(`expecting data to be a map`), m.UnmarshalJSON([]byte(`["hello",true,1,3.14,null]`)))
}

func TestMap_UnmarshalJSON_unbalancedKey(t *testing.T) {
	m := vf.MutableMap(nil)
	require.Equal(t, errors.New(`invalid character ']' after object key`), m.UnmarshalJSON([]byte(`{"a"]`)))
}

func TestMap_UnmarshalJSON_unbalancedEntry(t *testing.T) {
	m := vf.MutableMap(nil)
	require.Equal(t, errors.New(`invalid character ']' after object key:value pair`), m.UnmarshalJSON([]byte(`{"a":1]`)))
}

func TestMap_UnmarshalLong(t *testing.T) {
	b := bytes.NewBufferString(``)
	util.WriteByte(b, '{')
	for i := 0; i < 100; i++ {
		if i > 0 {
			util.WriteByte(b, ',')
		}
		util.Fprintf(b, `"%d":%d`, i, i)
	}
	util.WriteByte(b, '}')
	m := vf.MutableMap(nil)
	err := m.UnmarshalJSON(b.Bytes())
	require.Nil(t, err)
	require.Equal(t, 100, m.Len())
}
