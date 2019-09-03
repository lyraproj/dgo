package internal_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func ExampleUnmarshalJSON() {
	v, err := vf.UnmarshalJSON([]byte(`["hello",true,1,3.14,null,{"a":1}]`))
	if err == nil {
		fmt.Println(v.Equals(vf.Values(`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1})))
	}
	// Output: true
}

func ExampleMarshalJSON() {
	v, err := vf.MarshalJSON(vf.Values(
		`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1}))
	if err == nil {
		fmt.Println(string(v))
	}
	// Output: ["hello",true,1,3.14,null,{"a":1}]
}

func TestArray_MarshalJSON(t *testing.T) {
	a := vf.Values(
		`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1, "b": "two", "c": 2.17})
	b, err := a.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, `["hello",true,1,3.14,null,{"a":1,"b":"two","c":2.17}]`, string(b))
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
	m := vf.Map(map[string]interface{}{"a": 1, "b": "two", "c": vf.Values(`hello`, true, 1, 3.14, nil)})
	b, err := m.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, `{"a":1,"b":"two","c":["hello",true,1,3.14,null]}`, string(b))
}

func TestMap_UnmarshalJSON(t *testing.T) {
	m := vf.MutableMap(7, nil)
	err := m.UnmarshalJSON([]byte(`{"a":1,"b":"two","c":["hello",true,1,3.14,null]}`))
	require.Nil(t, err)
	require.Equal(t, vf.Map(map[string]interface{}{"a": 1, "b": "two", "c": vf.Values(`hello`, true, 1, 3.14, nil)}), m)
}

func TestMap_UnmarshalJSON_frozen(t *testing.T) {
	m := vf.Map(map[string]int{})
	require.Panic(t, func() { _ = m.UnmarshalJSON([]byte(`{"a":1,"b":"two","c":["hello",true,1,3.14,null]}`)) }, `frozen`)
}

func TestMap_UnmarshalJSON_notObject(t *testing.T) {
	m := vf.MutableMap(7, nil)
	require.Equal(t, errors.New(`expecting data to be an object`), m.UnmarshalJSON([]byte(`["hello",true,1,3.14,null]`)))
}

func TestMap_UnmarshalJSON_unbalancedKey(t *testing.T) {
	m := vf.MutableMap(7, nil)
	require.Equal(t, errors.New(`invalid character ']' after object key`), m.UnmarshalJSON([]byte(`{"a"]`)))
}

func TestMap_UnmarshalJSON_unbalancedEntry(t *testing.T) {
	m := vf.MutableMap(7, nil)
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
	m := vf.MutableMap(7, nil)
	err := m.UnmarshalJSON(b.Bytes())
	require.Nil(t, err)
	require.Equal(t, 100, m.Len())
}
