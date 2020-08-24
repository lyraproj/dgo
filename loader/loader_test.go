package loader_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/loader"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func testFinder() dgo.Finder {
	return func(_ dgo.Loader, key string) interface{} {
		return `the ` + key
	}
}

func testNamespace() dgo.NsCreator {
	return func(l dgo.Loader, key string) dgo.Loader {
		return loader.New(l, key, nil, testFinder(), testNamespace())
	}
}

func TestMapLoader_Equals(t *testing.T) {
	l := loader.New(nil, `my`, vf.Map(`a`, `the a`), nil, nil)
	l2 := l.Type().(dgo.NamedType).New(vf.Map(`name`, ``, `entries`, vf.Map(`a`, `the a`)))
	assert.Equal(t, l, l2)
	assert.Equal(t, l.HashCode(), l2.HashCode())
	assert.NotEqual(t, l, vf.Map(`a`, `the a`))
	assert.Equal(t, l.Type(), l2.Type())
	assert.NotSame(t, l.Type(), l2.Type())
	assert.Same(t, typ.Loader, typ.Generic(l.Type()))
	assert.Instance(t, l.Type(), l)
	assert.Equal(t, `the a`, l.Get(`a`))
	assert.Equal(t, `the a`, l.Load(`a`))
	assert.Equal(t, `my`, l.Name())
	assert.Equal(t, `/my`, l.AbsoluteName())
	assert.Equal(t, `mapLoader{"name":"my","entries":{"a":"the a"}}`, l.String())
	assert.Same(t, l, l.Namespace(``))
	assert.Nil(t, l.Namespace(`some`))
}

func TestMapLoader_AbsoluteName(t *testing.T) {
	l := loader.New(nil, `my`, vf.Map(`a`, `the a`), nil, func(l dgo.Loader, key string) dgo.Loader {
		return loader.New(l, key, vf.Map(`a`, `the a`), nil, nil)
	})
	ns := l.Namespace(`ns`)
	assert.NotNil(t, ns)
	assert.Equal(t, `ns`, ns.Name())
	assert.Equal(t, `/my/ns`, ns.AbsoluteName())
}

func TestLoader_Type(t *testing.T) {
	l := loader.New(nil, ``, nil,
		func(l dgo.Loader, key string) interface{} {
			return `the ` + key
		},
		func(l dgo.Loader, name string) dgo.Loader {
			return loader.New(l, name, nil, nil, nil)
		})

	tp := l.Type().(dgo.NamedType)
	l2 := l.Type().(dgo.NamedType).New(tp.ExtractInitArg(l))
	assert.Equal(t, l, l2)
	assert.Equal(t, l.HashCode(), l2.HashCode())
	assert.NotEqual(t, l, vf.Map(`a`, `the a`))
	assert.Equal(t, l.Type(), l2.Type())
	assert.NotSame(t, l.Type(), l2.Type())
	assert.Same(t, typ.DefiningLoader, typ.Generic(tp))
	assert.Instance(t, l.Type(), l)
	assert.Equal(t, l.Get(`a`), `the a`)
	assert.NotNil(t, l.Namespace(`some`))
	assert.Nil(t, l.ParentNamespace())
}

func TestLoader_Get_notString(t *testing.T) {
	l := loader.New(nil, ``, nil, func(l dgo.Loader, key string) interface{} { return key }, nil)
	assert.Nil(t, l.Get(23))
}

func TestLoader_Load(t *testing.T) {
	l := loader.New(nil, ``, nil, testFinder(), nil)
	assert.Equal(t, `the a`, l.Load(`a`))
}

func TestLoader_Load_nil(t *testing.T) {
	first := true
	l := loader.New(nil, ``, nil, func(l dgo.Loader, name string) interface{} {
		if first {
			first = false
		} else {
			t.Error(`finder called more than once`)
		}
		return nil
	}, nil)
	assert.Nil(t, l.Load(`a`))
	assert.Nil(t, l.Load(`a`))
}

func TestLoader_Load_multiple(t *testing.T) {
	l := loader.New(nil, ``, nil, func(l dgo.Loader, name string) interface{} {
		return loader.Multiple(vf.Map(
			`a`, `the a`,
			`b`, `the b`))
	}, nil)
	assert.Equal(t, `the a`, l.Load(`a`))
	assert.Equal(t, `the b`, l.Load(`b`))
}

func TestLoader_Load_multipleNotRequested(t *testing.T) {
	l := loader.New(nil, ``, nil, func(l dgo.Loader, name string) interface{} {
		return loader.Multiple(vf.Map(
			`c`, `the c`,
			`b`, `the b`))
	}, nil)
	assert.Panic(t, func() { l.Load(`a`) }, `map returned from finder doesn't contain original key "a"`)
}

func TestLoader_loadRedefine(t *testing.T) {
	first := true
	l := loader.New(nil, ``, nil, func(l dgo.Loader, key string) interface{} {
		if first {
			first = false
			return l.Load(key).String() + ` diff`
		}
		return "the " + key
	}, nil)
	assert.Panic(t, func() { l.Get(`a`) }, `attempt to override entry "a"`)
}

func TestLoader_NewChild(t *testing.T) {
	theA := vf.String(`the a`)
	l := loader.New(nil, ``, vf.Map(`a`, theA), nil, nil)
	c := l.NewChild(testFinder(), nil)
	assert.NotEqual(t, l.HashCode(), c.HashCode())

	assert.Equal(t, c.Get(`a`), `the a`)
	assert.Equal(t, c.Get(`b`), `the b`)
	assert.True(t, l.Get(`b`) == nil)

	assert.Equal(t,
		`childLoader{"loader":loader{"name":"","entries":{"b":"the b"},"namespaces":{}},"parent":mapLoader{"name":"","entries":{"a":"the a"}}}`,
		c.String())
	assert.Equal(t,
		`childLoader{"loader":loader{"name":"","entries":{"b":"the b"},"namespaces":{}},"parent":mapLoader{"name":"","entries":{"a":"the a"}}}`,
		c.Type().String())

	assert.Same(t, c.Get(`a`), theA)
	assert.Nil(t, c.Namespace(`ns`))
}

func TestLoader_Namespace(t *testing.T) {
	l := loader.New(nil, ``, nil, testFinder(), testNamespace())
	assert.Same(t, l, l.Namespace(``))

	b := l.NewChild(testFinder(), testNamespace())
	c := l.NewChild(testFinder(), nil)
	d := c.NewChild(testFinder(), testNamespace())

	assert.Equal(t, b, c)
	assert.NotEqual(t, b, l)

	assert.Same(t, b, b.Namespace(``))
	assert.Equal(t, b.Load(`b/x`), `the x`)
	assert.Equal(t, c.Load(`a/x`), `the x`)
	assert.Equal(t, d.Load(`a/x`), `the x`)
	assert.Equal(t, d.Load(`/a/x`), `the x`)

	tp := c.Type().(dgo.NamedType)
	cp := tp.New(tp.ExtractInitArg(c))
	assert.NotSame(t, c, cp)
	assert.Equal(t, c, cp)
}

func TestLoader_Namespace_noParentNs(t *testing.T) {
	l := loader.New(nil, ``, nil, testFinder(), nil)
	b := l.NewChild(testFinder(), testNamespace())
	assert.Nil(t, l.Load(`b/x`))
	assert.Equal(t, b.Load(`b/x`), `the x`)
}

func TestLoader_Namespace_redefinedOk(t *testing.T) {
	first := true
	l := loader.New(nil, ``, nil, nil, func(ld dgo.Loader, name string) dgo.Loader {
		if first {
			first = false
			return ld.Namespace(name)
		}
		return loader.New(nil, ``, vf.Map(`a`, `the a`), nil, nil)
	})
	ns := l.Namespace(`ns`)
	assert.Equal(t, `the a`, ns.Load(`a`))
}

func TestLoader_Namespace_redefinedBad(t *testing.T) {
	first := true
	l := loader.New(nil, ``, nil, nil, func(ld dgo.Loader, name string) dgo.Loader {
		if first {
			first = false
			ld.Namespace(name)
			return loader.New(nil, ``, vf.Map(`a`, `the b`), nil, nil)
		}
		return loader.New(nil, ``, vf.Map(`a`, `the a`), nil, nil)
	})
	assert.Panic(t, func() { l.Namespace(`ns`) }, `namespace "ns" is already defined`)
}

func TestLoader_parsed(t *testing.T) {
	tp := tf.ParseType(
		`childLoader{"loader":loader{"name":"", "entries":{"b":"the b"},"namespaces":{}},"parent":mapLoader{"name":"","entries":{"a":"the a"}}}`)
	et, ok := tp.(dgo.ExactType)
	require.True(t, ok)
	ld, ok := et.ExactValue().(dgo.Loader)
	require.True(t, ok)
	assert.Equal(t, `the b`, ld.Get(`b`))
	assert.Equal(t, `the a`, ld.Get(`a`))
}
