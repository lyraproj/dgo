package internal_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/tada/dgo/vf"

	"github.com/tada/dgo/dgo"

	"github.com/tada/dgo/typ"

	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/tf"
)

type testNamed int

func (a testNamed) String() string {
	return a.Type().(dgo.NamedType).ValueString(a)
}

func (a testNamed) Type() dgo.Type {
	return tf.ExactNamed(tf.Named(`testNamed`), a)
}

func (a testNamed) Equals(other interface{}) bool {
	return a == other
}

func (a testNamed) HashCode() int {
	return int(a)
}

type testNamedB struct {
	testNamed
}

type testNamedC int

type testNamedDummy interface {
	Dummy()
}

func (testNamed) Dummy() {
}

func (t *testNamedB) Format(s fmt.State, _ rune) {
	_, _ = s.Write([]byte(`hello from testNamedB Format`))
}

func TestNamedType(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil, nil)
	require.Equal(t, tp, tp)
	require.Equal(t, tp.Name(), `testNamed`)
	require.Equal(t, tp.String(), `testNamed`)
	require.NotEqual(t, tp, `testNamed`)

	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.Any)
	require.Instance(t, tp.Type(), tp)
	require.Instance(t, tp, testNamed(0))

	require.NotEqual(t, 0, tp.HashCode())
	require.Equal(t, tp.HashCode(), tp.HashCode())
}

func TestNamedType_redefined(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil, nil)
	require.Panic(t, func() {
		tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(&testNamedB{0}), nil, nil)
	}, `attempt to redefine named type 'testNamed'`)
}

func TestNamedTypeFromReflected(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil, nil)
	require.Same(t, tp, tf.NamedFromReflected(reflect.TypeOf(testNamed(0))))
	require.Nil(t, tf.NamedFromReflected(reflect.TypeOf(testNamedC(0))))
}

func TestNamedType_Assignable(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	defer tf.RemoveNamed(`testNamedB`)
	defer tf.RemoveNamed(`testNamedC`)
	tp := tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), reflect.TypeOf((*testNamedDummy)(nil)).Elem(), nil)
	require.Assignable(t, tp, tf.NewNamed(`testNamedB`, nil, nil, reflect.TypeOf(&testNamedB{0}), nil, nil))
	require.NotAssignable(t, tp, tf.NewNamed(`testNamedC`, nil, nil, reflect.TypeOf(testNamedC(0)), nil, nil))
	require.NotAssignable(t, tf.Named(`testNamedB`), tf.Named(`testNamedC`))
}

func TestNamedType_Format(t *testing.T) {
	defer tf.RemoveNamed(`testNamedB`)
	_ = tf.NewNamed(`testNamedB`, nil, nil, reflect.TypeOf(&testNamedB{0}), nil, nil)
	v := vf.Value(&testNamedB{0})
	require.Equal(t, `hello from testNamedB Format`, fmt.Sprintf("%v", v))
}

func TestNamedType_New(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)

	v := tp.New(vf.Integer(3))
	require.Equal(t, v, testNamed(3))
	require.Equal(t, 3, tp.ExtractInitArg(v))
}

func TestNamedType_New_notApplicable(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil, nil)

	require.Panic(t, func() { tp.New(vf.Integer(3)) }, `creating new instances of testNamed is not possible`)
	require.Panic(t, func() { tp.ExtractInitArg(testNamed(0)) }, `creating new instances of testNamed is not possible`)
}

func TestNamedType_ValueString(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)

	v := tp.New(vf.Integer(3))
	require.Equal(t, `testNamed 3`, tp.ValueString(v))
}

func TestNamedType_parse(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)
	require.Same(t, tf.ParseType(`testNamed`), tp)
}

func TestNamedType_exact(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)

	v := tp.New(vf.Integer(3))
	et := v.Type()
	require.Same(t, tp, typ.Generic(et))
	require.Assignable(t, tp, et)
	require.NotAssignable(t, et, tp)
	require.NotAssignable(t, et, tp.New(vf.Integer(4)).Type())
	require.Instance(t, et, v)
	require.NotInstance(t, et, tp.New(vf.Integer(4)))
	require.Equal(t, `testNamed 3`, et.String())
	require.Equal(t, et, tf.ParseType(`testNamed 3`))

	require.Instance(t, et.Type(), et)
	require.Instance(t, tp.Type(), et)
	require.NotInstance(t, et.Type(), tp)

	require.NotEqual(t, tp, et)
	require.Equal(t, et, tp.New(vf.Integer(3)).Type())
	require.NotEqual(t, et, tp.New(vf.Integer(3)))
	require.NotEqual(t, et, tp.New(vf.Integer(4)).Type())
	require.NotEqual(t, tp.HashCode(), et.HashCode())
}

func TestNamedType_parameterized(t *testing.T) {
	minMax := func(a dgo.Array) (int, int) {
		switch a.Len() {
		case 0:
			return 0, 0
		case 1:
			return int(a.Get(0).(dgo.Integer).GoInt()), math.MaxInt64
		default:
			return int(a.Get(0).(dgo.Integer).GoInt()), int(a.Get(1).(dgo.Integer).GoInt())
		}
	}

	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil,
		func(self dgo.NamedType, typ dgo.Type) bool {
			if ot, ok := typ.(dgo.NamedType); ok && self.Name() == ot.Name() {
				var oMin, oMax int
				if et, ok := ot.(dgo.ExactType); ok {
					oMin = int(et.ExactValue().(testNamed))
					oMax = oMin
				} else {
					oMin, oMax = minMax(ot.Parameters())
				}
				sMin, sMax := minMax(self.Parameters())
				if sMin <= oMin && oMax <= sMax {
					return true
				}
			}
			return false
		})

	tpp := tf.Parameterized(tp, vf.Values(0, 10))
	tpp2 := tf.Parameterized(tp, vf.Values(0, 10))
	require.Equal(t, `testNamed[0,10]`, tpp.String())
	require.Assignable(t, tpp, tpp2)
	require.NotEqual(t, tpp, tp)
	require.Equal(t, tpp, tpp2)
	require.Equal(t, tpp.Type(), tpp2.Type())
	require.Equal(t, tpp.HashCode(), tpp2.HashCode())
	require.Same(t, tp, typ.Generic(tpp))
	require.Same(t, tp, typ.Generic(tpp2))

	require.Instance(t, tpp, testNamed(3))
	require.NotInstance(t, tpp, testNamed(11))

	require.Panic(t, func() { vf.New(tpp, vf.Integer(11)) },
		`the value testNamed 11 cannot be assigned to a variable of type testNamed\[0,10\]`)
}

func TestNamedType_parameterized_noAsgChecker(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)

	tpp := tf.Parameterized(tp, vf.Values(0, 10))
	require.Equal(t, `testNamed[0,10]`, tpp.String())
	require.Assignable(t, tpp, tp)
	require.NotEqual(t, tpp, tp)
	require.Instance(t, tpp, testNamed(3))
	require.Instance(t, tpp, testNamed(11))
}
