package streamer_test

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/vf"
)

type badWriter int

func (b badWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New(`bang`)
}

func ExampleJSON() {
	s := streamer.New(nil, nil)
	s.Stream(vf.Map(`a`, 1, `b`, []int{1, 2}), streamer.JSON(os.Stdout))
	// Output: {"a":1,"b":[1,2]}
}

func TestJSON_AddRef(t *testing.T) {
	v := vf.Strings(`a`, `b`)
	a := vf.Values(v, v)
	s := streamer.New(nil, nil)
	b := bytes.Buffer{}
	s.Stream(a, streamer.JSON(&b))
	assert.Equal(t, `[["a","b"],{"__ref":1}]`, b.String())
}

func TestJSON_primitives(t *testing.T) {
	v := vf.Values(true, nil, 1, 2.1, `string`)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(v, streamer.JSON(&b))
	assert.Equal(t, `[true,null,1,2.1,"string"]`, b.String())
}

func TestJSON_badWrite(t *testing.T) {
	assert.Panic(t, func() { streamer.New(nil, nil).Stream(vf.Int64(3), streamer.JSON(badWriter(0))) }, `bang`)
}

func TestJSON_CanDoBinary(t *testing.T) {
	v := vf.Values(vf.BinaryFromString(`AQID`))
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(v, streamer.JSON(&b))
	assert.Equal(t, `[{"__type":"binary","__value":"AQID"}]`, b.String())
}

func TestJSON_CanDoDuration(t *testing.T) {
	d, _ := time.ParseDuration(`4m3s`)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(vf.Duration(d), streamer.JSON(&b))
	require.Equal(t, `{"__type":"duration","__value":"4m3s"}`, b.String())
}

func TestJSON_CanDoTime(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(vf.Time(ts), streamer.JSON(&b))
	assert.Equal(t, `{"__type":"time","__value":"2019-10-06T07:15:00-07:00"}`, b.String())
}

func TestJSON_CanDoBigInt(t *testing.T) {
	bi, _ := new(big.Int).SetString(`0x10000000000000000`, 0)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(vf.BigInt(bi), streamer.JSON(&b))
	assert.Equal(t, `{"__type":"bigint","__value":"10000000000000000"}`, b.String())

	bv, ok := streamer.UnmarshalJSON(b.Bytes(), nil).(dgo.BigInt)
	assert.True(t, ok)
	assert.True(t, bi.Cmp(bv.GoBigInt()) == 0)
}

func TestJSON_CanDoBigFloat(t *testing.T) {
	bf, _ := new(big.Float).SetString(`1e10000`)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(vf.BigFloat(bf), streamer.JSON(&b))
	assert.Equal(t, `{"__type":"bigfloat","__value":"0x.9b84ea28556bf269p+33220"}`, b.String())

	bv, ok := streamer.UnmarshalJSON(b.Bytes(), nil).(dgo.BigFloat)
	require.True(t, ok)
	assert.True(t, bf.Cmp(bv.GoBigFloat()) == 0)
}

func TestJSON_CanDoUint64(t *testing.T) {
	ui := uint64(0x8000000000000000)
	b := bytes.Buffer{}
	streamer.New(nil, nil).Stream(vf.Uint64(ui), streamer.JSON(&b))
	assert.Equal(t, `9223372036854775808`, b.String())

	uv, ok := streamer.UnmarshalJSON(b.Bytes(), nil).(dgo.Uint64)
	require.True(t, ok)
	assert.Equal(t, ui, uv.GoUint())
}

func TestJSON_ComplexKeys(t *testing.T) {
	v := vf.Map(vf.BinaryFromString(`AQID`), `value of binary`, `hey`, `value of hey`)
	b := bytes.Buffer{}
	streamer.New(nil, streamer.DefaultOptions()).Stream(v, streamer.JSON(&b))
	assert.Equal(t,
		`{"__type":"map","__value":[{"__type":"binary","__value":"AQID"},"value of binary","hey","value of hey"]}`,
		b.String())
}

func TestUnmarshalJSON_ref(t *testing.T) {
	v := streamer.UnmarshalJSON(
		[]byte(`[{"x":"xxxxxxxxxxxxxxxxxxxxx","y":{"__ref":3}}]`),
		streamer.DgoDialect())
	assert.Equal(t, vf.Values(vf.Map(`x`, `xxxxxxxxxxxxxxxxxxxxx`, `y`, `xxxxxxxxxxxxxxxxxxxxx`)), v)
}

func TestUnmarshalJSON_refBad(t *testing.T) {
	assert.Panic(t, func() {
		streamer.UnmarshalJSON([]byte(`{"x": {"z": ["a","b"]}, "y": {"__ref":"6"}}`), streamer.DgoDialect())
	}, `expected integer after key "__ref"`)
}

func TestUnmarshalJSON_refSuperfluous(t *testing.T) {
	assert.Panic(t, func() {
		streamer.UnmarshalJSON([]byte(`{"x": {"z": ["a","b"]}, "y": {"__ref":6, "x": 1}}`), streamer.DgoDialect())
	}, `expected end of object after "__ref": 6`)
}

func TestUnmarshalJSON_badDelim(t *testing.T) {
	assert.Panic(t, func() {
		streamer.UnmarshalJSON([]byte(``), streamer.DgoDialect())
	}, `unexpected EOF`)
}

func TestUnmarshalJSON_complexKeys(t *testing.T) {
	v := streamer.UnmarshalJSON(
		[]byte(`{"__type":"map","__value":[{"__type":"binary","__value":"AQID"},"value of binary","hey","value of hey"]}`),
		streamer.DgoDialect())
	v2 := vf.Map(vf.BinaryFromString(`AQID`), `value of binary`, `hey`, `value of hey`)
	assert.Equal(t, v, v2)
}

func TestUnmarshalJSON_badInput(t *testing.T) {
	assert.Panic(t, func() { streamer.UnmarshalJSON([]byte(`this is not json`), nil) }, `invalid character`)
}

func ExampleUnmarshalJSON() {
	v := streamer.UnmarshalJSON([]byte(`["hello",true,1,3.14,null,{"a":1}]`), nil)
	fmt.Println(v.Equals(vf.Values(`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1})))
	// Output: true
}

func ExampleMarshalJSON_slice() {
	v := streamer.MarshalJSON(vf.Values(
		`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1}), streamer.DgoDialect())
	fmt.Println(string(v))
	// Output: ["hello",true,1,3.14,null,{"a":1}]
}

func ExampleMarshalJSON_string() {
	v := streamer.MarshalJSON("hello", nil)
	fmt.Println(string(v))
	// Output: "hello"
}
