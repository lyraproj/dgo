package internal

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
)

const elemCount = 500000

func x(_ dgo.Value) {}

// BenchmarkRangeIndex `for _, e := range s { x(e) }`
func BenchmarkRange(b *testing.B) {
	s := make([]dgo.Value, elemCount)
	for i := 0; i < elemCount; i++ {
		s[i] = Integer(i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, e := range s {
			x(e)
		}
	}
}

// BenchmarkRangeIndex `for i := range s { x(s[i]) }`
//
// about 2.5 times faster than BenchmarkRange
func BenchmarkRangeIndex(b *testing.B) {
	s := make([]dgo.Value, elemCount)
	for i := 0; i < elemCount; i++ {
		s[i] = Integer(i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := range s {
			x(s[i])
		}
	}
}

// BenchmarkLenIndex `for i := 0; i < l; i++ { x(s[i]) }`
//
// on par with the BenchmarkRangeIndex, so obviously not worth doing
func BenchmarkLenIndex(b *testing.B) {
	s := make([]dgo.Value, elemCount)
	for i := 0; i < elemCount; i++ {
		s[i] = Integer(i)
	}

	b.ResetTimer()
	l := len(s)
	for n := 0; n < b.N; n++ {
		for i := 0; i < l; i++ {
			x(s[i])
		}
	}
}

type hsDp struct {
	hstring
}

func (h *hsDp) deepCompare(seen []dgo.Value, other deepCompare) (int, bool) {
	return 1, true
}

type dummyDp int

func (d dummyDp) deepEqual(seen []dgo.Value, other deepEqual) bool {
	return true
}

func (d dummyDp) deepHashCode(seen []dgo.Value) int {
	return 1
}

func (d dummyDp) String() string {
	return ``
}

func (d dummyDp) Type() dgo.Type {
	return nil
}

func (d dummyDp) Equals(other interface{}) bool {
	return true
}

func (d dummyDp) HashCode() int {
	return 0
}

func (d dummyDp) deepCompare(seen []dgo.Value, other deepCompare) (int, bool) {
	return 1, true
}

func Test_deepCompareWrongType(t *testing.T) {
	_, ok := (&array{}).deepCompare(nil, &hsDp{*makeHString(``)})
	if ok {
		t.FailNow()
	}
}

func Test_equalsWrongKind(t *testing.T) {
	ok := equals(nil, &array{}, dummyDp(0))
	if ok {
		t.FailNow()
	}
}

func Test_compareWrongKind(t *testing.T) {
	c, ok := compare(nil, &array{}, dummyDp(0))
	if !ok && c == -1 {
		t.FailNow()
	}
}

func isSame(a, b interface{}) bool {
	return a == b
}

func Test_ArrayEquality(t *testing.T) {
	a := &array{slice: []dgo.Value{Nil}}
	b := &array{slice: []dgo.Value{Nil}}
	if isSame(a, b) {
		t.Error(`== on array goes deep`)
	}
	if !isSame(a, a) {
		t.Error(`== on array isn't true for same object`)
	}
}
