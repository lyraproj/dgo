package streamer

import (
	"fmt"
	"time"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
)

type dataDecoder struct {
	BasicCollector
	dialect  Dialect
	aliasMap dgo.AliasMap
}

// DataDecoder returns a decoder capable of decoding a stream of rich data representations into the corresponding values
func DataDecoder(aliasMap dgo.AliasMap, d Dialect) Collector {
	c := &dataDecoder{aliasMap: aliasMap, dialect: d}
	c.Init()
	return c
}

func (d *dataDecoder) Init() {
	d.BasicCollector.Init()
	if d.dialect == nil {
		d.dialect = DgoDialect()
	}
}

// AddMap initializes and adds a new map and then calls the function with is supposed to
// add an even number of elements as a sequence of key, value, [key, value, ...]
func (d *dataDecoder) AddMap(cap int, doer dgo.Doer) {
	d.BasicCollector.AddMap(cap, doer)
	m := d.PeekLast().(dgo.Map)
	dl := d.dialect
	if ts := m.Get(dl.TypeKey()); ts != nil {
		d.ReplaceLast(d.decode(m))
	}
}

func (d *dataDecoder) decode(v dgo.Value) dgo.Value {
	if m, ok := v.(dgo.Map); ok {
		dl := d.dialect
		if ts := m.Get(dl.TypeKey()); ts != nil {
			mv := m.Get(dl.ValueKey())
			switch {
			case mv == nil:
				v = dl.ParseType(ts.(dgo.String))
			case ts.Equals(dl.MapTypeName()):
				nm := mv.(dgo.Array).ToMap()
				// Replace all occurrences of v in the new map recursively with the new map as it
				// might contain references to itself
				replaceInstance(v, nm, nm)
				v = nm
			case ts.Equals(dl.SensitiveTypeName()):
				v = vf.Sensitive(mv)
			case ts.Equals(dl.BinaryTypeName()):
				v = vf.BinaryFromString(mv.(dgo.String).GoString())
			case ts.Equals(dl.TimeTypeName()):
				ts, err := time.Parse(time.RFC3339Nano, mv.(dgo.String).GoString())
				if err != nil {
					panic(err)
				}
				v = vf.Time(ts)
			case ts.Equals(dl.AliasTypeName()):
				ad := mv.(dgo.Array)
				v = dl.ParseType(ad.Get(1).(dgo.String))
				if d.aliasMap != nil {
					d.aliasMap.Add(v.(dgo.Type), ad.Get(0).(dgo.String))
				}
			default:
				panic(fmt.Errorf(`unable to decode %s: %s`, dl.TypeKey(), ts))
			}
		}
	}
	return v
}

func replaceInstance(orig, repl, in dgo.Value) (dgo.Value, bool) {
	if in == orig {
		return repl, true
	}

	replaceHappened := false
	switch iv := in.(type) {
	case dgo.Map:
		iv.EachEntry(func(v dgo.MapEntry) {
			if re, rh := replaceInstance(orig, repl, v.Value()); rh {
				replaceHappened = true
				iv.Put(v.Key(), re)
			}
		})
	case dgo.Array:
		iv.EachWithIndex(func(v dgo.Value, i int) {
			if re, rh := replaceInstance(orig, repl, v); rh {
				replaceHappened = true
				iv.Set(i, re)
			}
		})
	case dgo.Sensitive:
		if rw, rh := replaceInstance(orig, repl, iv.Unwrap()); rh {
			replaceHappened = true
			in = vf.Sensitive(rw)
		}
	}
	return in, replaceHappened
}
