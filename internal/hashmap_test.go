package internal

import (
	"bytes"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

const rndFactor = 1234
const rndLen = 11
const lookupsPerOp = 50

func buildStringKeys(count int) []string {
	k := make([]string, count)
	b := bytes.NewBufferString(``)
	for i := 0; i < count; i++ {
		s := strconv.Itoa(rand.Intn(count * rndFactor))
		b.Reset()
		l := len(s)
		p := 0
		for r := 0; r < rand.Intn(rndLen)+l; r++ {
			util.WriteByte(b, s[p])
			p++
			if p >= l {
				p = 0
			}
		}
		k[i] = b.String()
	}
	return k
}

func BenchmarkHashMapStrings(b *testing.B) {
	sz := b.N + lookupsPerOp
	ks := buildStringKeys(sz)
	k := make([]dgo.Value, sz)
	c := float64(sz) / loadFactor
	m := MutableMap(int(c), nil)
	for i := 0; i < sz; i++ {
		key := makeHString(ks[i])
		m.Put(key, intVal(i))
		k[i] = key
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := 0; v < lookupsPerOp; v++ {
			m.Get(k[i+v])
		}
	}
}

// BenchmarkHashMapStringsNoHashCache measures lookup using strings that have
// no precomputed hash code by resetting the cache prior to each lookup
func BenchmarkHashMapStringsNoHashCache(b *testing.B) {
	sz := b.N + lookupsPerOp
	ks := buildStringKeys(sz)
	k := make([]*hstring, sz)
	c := float64(sz) / loadFactor
	m := MutableMap(int(c), nil)
	for i := 0; i < sz; i++ {
		key := makeHString(ks[i])
		m.Put(key, intVal(i))
		k[i] = key
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := 0; v < lookupsPerOp; v++ {
			key := k[i+v]
			key.h = 0
			m.Get(key)
		}
	}
}

func BenchmarkHashMapIntegers(b *testing.B) {
	sz := b.N + lookupsPerOp
	k := make([]dgo.Value, sz)
	c := float64(sz) / loadFactor
	m := MutableMap(int(c), nil)
	for i := 0; i < sz; i++ {
		key := intVal(rand.Intn(sz * rndFactor))
		m.Put(key, intVal(i))
		k[i] = key
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := 0; v < lookupsPerOp; v++ {
			m.Get(k[i+v])
		}
	}
}

func BenchmarkMapStringKeys(b *testing.B) {
	sz := b.N + lookupsPerOp
	c := float64(sz) / loadFactor
	m := make(map[string]dgo.Value, int(c))

	k := buildStringKeys(sz)
	for i := 0; i < sz; i++ {
		m[k[i]] = intVal(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := 0; v < lookupsPerOp; v++ {
			_ = m[k[i+v]]
		}
	}
}

func BenchmarkMapIntKeys(b *testing.B) {
	sz := b.N + lookupsPerOp
	k := make([]int, sz)
	c := float64(sz) / loadFactor
	m := make(map[int]int, int(c))

	for i := 0; i < sz; i++ {
		key := rand.Intn(sz * rndFactor)
		m[key] = i
		k[i] = key
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := 0; v < lookupsPerOp; v++ {
			_ = m[k[i+v]]
		}
	}
}

func BenchmarkMapGenericKeysString(b *testing.B) {
	sz := b.N + lookupsPerOp
	k := buildStringKeys(sz)
	c := float64(sz) / loadFactor
	m := make(map[interface{}]int, int(c))

	for i := 0; i < sz; i++ {
		m[k[i]] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := 0; v < lookupsPerOp; v++ {
			_ = m[k[i+v]]
		}
	}
}

func BenchmarkMapGenericKeysInteger(b *testing.B) {
	sz := b.N + lookupsPerOp
	k := make([]int, sz)
	c := float64(sz) / loadFactor
	m := make(map[interface{}]int, int(c))

	for i := 0; i < sz; i++ {
		key := rand.Intn(sz * rndFactor)
		m[key] = i
		k[i] = key
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := 0; v < lookupsPerOp; v++ {
			_ = m[k[i+v]]
		}
	}
}

func Test_tableSize(t *testing.T) {
	if tableSizeFor(0) != 1 {
		t.Fatal(`table size 0 is not 1`)
	}
	if tableSizeFor(1) != 1 {
		t.Fatal(`table size 1 is not 1`)
	}
	if tableSizeFor(math.MaxInt64) != maximumCapacity {
		t.Fatal(`table size is not capped to maximum capacity`)
	}
}
