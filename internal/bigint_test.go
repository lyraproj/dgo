package internal_test

import (
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/test/require"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func newBigInt(t *testing.T, s string) dgo.BigInt {
	t.Helper()
	bi, ok := new(big.Int).SetString(s, 0)
	require.True(t, ok)
	return vf.BigInt(bi)
}

func newBigIntRange(t *testing.T, min string, max string, inclusive bool) dgo.IntegerType {
	t.Helper()
	var minBi, maxBi dgo.BigInt
	if min != "" {
		minBi = newBigInt(t, min)
	}
	if max != "" {
		maxBi = newBigInt(t, max)
	}
	return tf.Integer(minBi, maxBi, inclusive)
}

func TestBigIntType_New(t *testing.T) {
	ex, _ := new(big.Int).SetString(`42`, 10)
	assert.True(t, ex.Cmp(newBigInt(t, `42`).GoBigInt()) == 0)
	assert.True(t, ex.Cmp(
		vf.New(newBigIntRange(t, `20`, `45`, true), vf.String(`42`)).(dgo.BigInt).GoBigInt()) == 0)
	assert.True(t, ex.Cmp(
		vf.New(newBigInt(t, `42`).Type(), vf.String(`42`)).(dgo.BigInt).GoBigInt()) == 0)

	assert.Equal(t, new(big.Int).SetInt64(123), vf.New(typ.BigInt, vf.Float(123)))
	assert.Equal(t, big.NewInt(123), vf.New(typ.BigInt, vf.Integer(123)))
	assert.Equal(t, big.NewInt(0), vf.New(typ.BigInt, vf.Boolean(false)))
	assert.Equal(t, big.NewInt(1), vf.New(typ.BigInt, vf.Boolean(true)))

	assert.Equal(t, ex, vf.New(typ.BigInt, vf.Arguments(`2a`, 16)))
	assert.Equal(t, ex, vf.New(typ.BigInt, vf.Arguments(`0x2a`)))

	assert.Panic(t, func() {
		vf.New(typ.BigInt, vf.Binary([]byte{1, 2}, false))
	}, `cannot be converted`)

	assert.Panic(t, func() {
		vf.New(newBigIntRange(t, `4`, `10`, true), vf.Integer(3))
	}, `cannot be assigned`)
}

func TestBigIntType_Range(t *testing.T) {
	tp := newBigIntRange(t, `123`, `678`, true)
	assert.Assignable(t, tp, newBigIntRange(t, `123`, `678`, true))
	assert.Assignable(t, tp, newBigIntRange(t, `234`, `678`, true))
	assert.Assignable(t, tp, newBigIntRange(t, `123`, `567`, true))
	assert.Assignable(t, tp, newBigIntRange(t, `234`, `567`, true))
	assert.Assignable(t, tp, newBigIntRange(t, `123`, `679`, false))
	assert.NotAssignable(t, newBigIntRange(t, `123`, `678`, false), tp)

	assert.Assignable(t, newBigInt(t, `345`).Type(), newBigInt(t, `345`).Type())
	assert.Assignable(t, typ.Generic(newBigInt(t, `345`).Type()), tp)
	assert.NotAssignable(t, newBigInt(t, `345`).Type(), tp)

	assert.Instance(t, tp, newBigInt(t, `345`))
	assert.Instance(t, newBigInt(t, `345`).Type(), newBigInt(t, `345`))
	assert.NotInstance(t, tp, newBigInt(t, `120`))

	assert.Equal(t, `123..678`, tp.String())
}

func TestBigIntType_ReflectType(t *testing.T) {
	tp := typ.BigInt
	assert.Equal(t, tp.ReflectType(), reflect.TypeOf(new(big.Int)))
	tp = newBigIntRange(t, `123`, `678`, true)
	assert.Equal(t, tp.ReflectType(), reflect.TypeOf(new(big.Int)))
	tp = newBigInt(t, `123`).(dgo.IntegerType)
	assert.Equal(t, tp.ReflectType(), reflect.TypeOf(new(big.Int)))
}

func TestBigIntType_TypeIdentifier(t *testing.T) {
	tp := typ.BigInt
	assert.Equal(t, tp.TypeIdentifier(), dgo.TiBigInt)
	tp = newBigIntRange(t, `123`, `678`, true)
	assert.Equal(t, tp.TypeIdentifier(), dgo.TiBigIntRange)
	tp = newBigInt(t, `123`).(dgo.IntegerType)
	assert.Equal(t, tp.TypeIdentifier(), dgo.TiBigIntExact)
}

func TestBigInt_CompareTo(t *testing.T) {
	b := newBigInt(t, `123`)
	c, ok := b.CompareTo(b)
	require.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = b.CompareTo(big.NewInt(140))
	require.True(t, ok)
	assert.Equal(t, -1, c)

	_, ok = b.CompareTo(`123`)
	assert.False(t, ok)

	c, ok = b.CompareTo(nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(120.0)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(vf.Float(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(120)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(vf.Integer(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(float32(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(float64(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(int16(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(int8(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(int32(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(int64(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(uint8(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(uint16(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(uint32(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(uint64(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(uint(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(big.NewFloat(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(big.NewFloat(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(big.NewFloat(123))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = b.CompareTo(big.NewFloat(123.1))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	b = vf.BigInt(new(big.Int).SetUint64(uint64(math.MaxUint64)))
	c, ok = b.CompareTo(uint(math.MaxUint64))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = b.CompareTo(uint64(math.MaxUint64))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = b.CompareTo(math.MaxInt64)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	i := vf.Integer(123)
	c, ok = i.CompareTo(big.NewInt(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = i.CompareTo(vf.BigInt(big.NewInt(120)))
	assert.True(t, ok)
	assert.Equal(t, 1, c)
}

func TestBigInt_Equals(t *testing.T) {
	b := vf.BigInt(big.NewInt(140))
	assert.True(t, b.Equals(b))
	assert.True(t, b.Equals(big.NewInt(140)))
	assert.True(t, b.Equals(vf.BigInt(big.NewInt(140))))
	assert.True(t, b.Equals(140))
	assert.True(t, b.Equals(uint(140)))
	assert.True(t, b.Equals(vf.Integer(140)))
	assert.False(t, b.Equals(float64(140)))

	b = vf.BigInt(new(big.Int).SetUint64(uint64(math.MaxUint64)))
	assert.True(t, b.Equals(uint64(math.MaxUint64)))
	assert.False(t, b.Equals(float64(math.MaxUint64)))
}

func TestBigInt_Int(t *testing.T) {
	b := vf.BigInt(big.NewInt(140))
	assert.Same(t, b, b.Integer())
}

func TestBigInt_GoInt(t *testing.T) {
	assert.Equal(t, 123, vf.BigInt(big.NewInt(123)).GoInt())
}

func TestBigInt_Dec(t *testing.T) {
	assert.Equal(t, 122, vf.BigInt(big.NewInt(123)).Dec())
}

func TestBigInt_Inc(t *testing.T) {
	assert.Equal(t, 124, vf.BigInt(big.NewInt(123)).Inc())
}

func TestBigInt_HashCode(t *testing.T) {
	assert.Equal(t, vf.BigInt(big.NewInt(123)).HashCode(), vf.BigInt(big.NewInt(123)).HashCode())
	assert.NotEqual(t, vf.BigInt(big.NewInt(123)).HashCode(), vf.BigInt(big.NewInt(124)).HashCode())
}

func TestBigInt_GoInt_out_of_bounds(t *testing.T) {
	assert.Panic(t, func() {
		newBigInt(t, `0x10000000000000000`).GoInt()
	}, `cannot fit`)
}

func TestBigInt_Float(t *testing.T) {
	bi := vf.BigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(1000), nil))
	bf := vf.BigFloat(new(big.Float).SetInt(bi.GoBigInt()))
	iv := bi.Float()
	assert.Equal(t, iv, bf)

	_, ok := bi.ToFloat()
	assert.False(t, ok)

	bi = vf.BigInt(big.NewInt(123))
	f, ok := bi.ToFloat()
	assert.True(t, ok)
	assert.Equal(t, f, float64(123))
}

func TestBigInt_ReflectTo(t *testing.T) {
	bi := newBigInt(t, `1234`)
	var mi interface{}
	bi.ReflectTo(reflect.ValueOf(&mi).Elem())
	fc, ok := mi.(*big.Int)
	require.True(t, ok)
	assert.True(t, fc.Cmp(bi.GoBigInt()) == 0)

	var mf *big.Int
	bi.ReflectTo(reflect.ValueOf(&mf).Elem())
	assert.True(t, mf.Cmp(bi.GoBigInt()) == 0)

	var m big.Int
	bi.ReflectTo(reflect.ValueOf(&m).Elem())
	assert.True(t, m.Cmp(bi.GoBigInt()) == 0)
}

func TestBigInt_String(t *testing.T) {
	// Number that requires more than 53 bits precision
	b := vf.BigInt(big.NewInt(123))
	assert.Equal(t, `123`, b.String())
}

func TestBigInt_ToBigInt(t *testing.T) {
	// Number that requires more than 53 bits precision
	b := big.NewInt(123)
	assert.True(t, b == vf.BigInt(b).ToBigInt())
}

func TestBigInt_ToBigFloat(t *testing.T) {
	// Number that requires more than 53 bits precision
	b := big.NewInt(123)
	assert.Equal(t, big.NewFloat(123), vf.BigInt(b).ToBigFloat())
}

func TestBigInt_ToInt(t *testing.T) {
	// Number that requires more than 53 bits precision
	i := vf.Integer(123)
	bf := vf.BigInt(big.NewInt(123))
	f2, ok := bf.ToInt()
	assert.True(t, ok)
	assert.Equal(t, i, f2)
}

func TestBigInt_ToInt_out_of_bounds(t *testing.T) {
	// Number that requires more than 53 bits precision
	_, ok := newBigInt(t, `0x10000000000000000`).ToInt()
	assert.False(t, ok)
	_, ok = newBigInt(t, `-0x10000000000000000`).ToInt()
	assert.False(t, ok)
}

func TestBigInt_range_specific(t *testing.T) {
	it := newBigInt(t, `0x10000000000000000`).(dgo.IntegerType)
	assert.True(t, it.Inclusive())
	assert.Same(t, it, it.Min())
	assert.Same(t, it, it.Max())
}
