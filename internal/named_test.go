package internal_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
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

func (a testNamed) HashCode() dgo.Hash {
	return dgo.Hash(a)
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
	assert.Equal(t, tp, tp)
	assert.Equal(t, tp.Name(), `testNamed`)
	assert.Equal(t, tp.String(), `testNamed`)
	assert.NotEqual(t, tp, `testNamed`)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.Any)
	assert.Instance(t, tp.Type(), tp)
	assert.Instance(t, tp, testNamed(0))
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Equal(t, tp.HashCode(), tp.HashCode())
}

func TestNamedType_redefined(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil, nil)
	assert.Panic(t, func() {
		tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(&testNamedB{0}), nil, nil)
	}, `attempt to redefine named type 'testNamed'`)
}

func TestNamedTypeFromReflected(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil, nil)
	assert.Same(t, tp, tf.NamedFromReflected(reflect.TypeOf(testNamed(0))))
	assert.Nil(t, tf.NamedFromReflected(reflect.TypeOf(testNamedC(0))))
}

func TestNamedType_Assignable(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	defer tf.RemoveNamed(`testNamedB`)
	defer tf.RemoveNamed(`testNamedC`)
	tp := tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), reflect.TypeOf((*testNamedDummy)(nil)).Elem(), nil)
	assert.Assignable(t, tp, tf.NewNamed(`testNamedB`, nil, nil, reflect.TypeOf(&testNamedB{0}), nil, nil))
	assert.NotAssignable(t, tp, tf.NewNamed(`testNamedC`, nil, nil, reflect.TypeOf(testNamedC(0)), nil, nil))
	assert.NotAssignable(t, tf.Named(`testNamedB`), tf.Named(`testNamedC`))
}

func TestNamedType_Format(t *testing.T) {
	defer tf.RemoveNamed(`testNamedB`)
	_ = tf.NewNamed(`testNamedB`, nil, nil, reflect.TypeOf(&testNamedB{0}), nil, nil)
	v := vf.Value(&testNamedB{0})
	assert.Equal(t, `hello from testNamedB Format`, fmt.Sprintf("%v", v))
}

func TestNamedType_New(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Int64(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)

	v := tp.New(vf.Int64(3))
	assert.Equal(t, v, testNamed(3))
	assert.Equal(t, 3, tp.ExtractInitArg(v))
}

func TestNamedType_New_notApplicable(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil, nil)
	assert.Panic(t, func() { tp.New(vf.Int64(3)) }, `creating new instances of testNamed is not possible`)
	assert.Panic(t, func() { tp.ExtractInitArg(testNamed(0)) }, `creating new instances of testNamed is not possible`)
}

func TestNamedType_ValueString(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Int64(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)

	v := tp.New(vf.Int64(3))
	assert.Equal(t, `testNamed 3`, tp.ValueString(v))
}

func TestNamedType_parse(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Int64(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)
	assert.Same(t, tf.ParseType(`testNamed`), tp)
}

func TestNamedType_exact(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Int64(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)

	v := tp.New(vf.Int64(3))
	et := v.Type()
	assert.Same(t, tp, typ.Generic(et))
	assert.Assignable(t, tp, et)
	assert.NotAssignable(t, et, tp)
	assert.NotAssignable(t, et, tp.New(vf.Int64(4)).Type())
	assert.Instance(t, et, v)
	assert.NotInstance(t, et, tp.New(vf.Int64(4)))
	assert.Equal(t, `testNamed 3`, et.String())
	assert.Equal(t, et, tf.ParseType(`testNamed 3`))
	assert.Instance(t, et.Type(), et)
	assert.Instance(t, tp.Type(), et)
	assert.NotInstance(t, et.Type(), tp)
	assert.NotEqual(t, tp, et)
	assert.Equal(t, et, tp.New(vf.Int64(3)).Type())
	assert.NotEqual(t, et, tp.New(vf.Int64(3)))
	assert.NotEqual(t, et, tp.New(vf.Int64(4)).Type())
	assert.NotEqual(t, tp.HashCode(), et.HashCode())
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
		return vf.Int64(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil,
		func(self dgo.NamedType, typ dgo.Type) bool {
			if ot, ok := typ.(dgo.NamedType); ok && self.Name() == ot.Name() {
				var oMin, oMax int
				var et dgo.ExactType
				if et, ok = ot.(dgo.ExactType); ok {
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
	assert.Equal(t, `testNamed[0,10]`, tpp.String())
	assert.Assignable(t, tpp, tpp2)
	assert.NotEqual(t, tpp, tp)
	assert.Equal(t, tpp, tpp2)
	assert.Equal(t, tpp.Type(), tpp2.Type())
	assert.Equal(t, tpp.HashCode(), tpp2.HashCode())
	assert.Same(t, tp, typ.Generic(tpp))
	assert.Same(t, tp, typ.Generic(tpp2))
	assert.Instance(t, tpp, testNamed(3))
	assert.NotInstance(t, tpp, testNamed(11))
	assert.Panic(t, func() { vf.New(tpp, vf.Int64(11)) },
		`the value testNamed 11 cannot be assigned to a variable of type testNamed\[0,10\]`)
}

func TestNamedType_parameterized_noAsgChecker(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Int64(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil, nil)

	tpp := tf.Parameterized(tp, vf.Values(0, 10))
	assert.Equal(t, `testNamed[0,10]`, tpp.String())
	assert.Assignable(t, tpp, tp)
	assert.NotEqual(t, tpp, tp)
	assert.Instance(t, tpp, testNamed(3))
	assert.Instance(t, tpp, testNamed(11))
}
