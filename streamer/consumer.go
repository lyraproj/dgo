package streamer

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
)

// A Consumer is used by a Data streaming mechanism that maintains a reference index
// which is increased by one for each value that it streams. The reference index
// originates from zero.
type Consumer interface {
	// CanDoBinary returns true if the value can handle binary efficiently. This tells
	// the Serializer to pass dgo.Binary verbatim to Add
	CanDoBinary() bool

	// CanDoTime returns true if the value can handle timestamp efficiently. This tells
	// the Serializer to pass dgo.Time verbatim to Add
	CanDoTime() bool

	// CanDoComplexKeys returns true if complex values can be used as keys. If this
	// method returns false, all keys must be strings
	CanDoComplexKeys() bool

	// StringDedupThreshold returns the preferred threshold for dedup of strings. Strings
	// shorter than this threshold will not be subjected to de-duplication.
	StringDedupThreshold() int

	// AddArray starts a new array, calls the doer function, and then ends the Array.
	//
	// The callers reference index is increased by one.
	AddArray(len int, doer dgo.Doer)

	// AddMap starts a new map, calls the doer function, and then ends the map.
	//
	// The callers reference index is increased by one.
	AddMap(len int, doer dgo.Doer)

	// Add adds the next value.
	//
	// Calls following a StartArray will add elements to the Array
	//
	// Calls following a StartHash will first add a key, then a value. This
	// repeats until End or StartArray is called.
	//
	// The callers reference index is increased by one.
	Add(element dgo.Value)

	// AddRef adds a reference to a previously added afterElement, hash, or array.
	AddRef(ref int)
}

// A Collector receives streaming events and produces an Value
type Collector interface {
	Consumer

	Value() dgo.Value
}

// A BasicCollector is an extendable basic implementation of the Consumer interface
type BasicCollector struct {
	// Values is an array of all values that are added to the BasicCollector. When adding
	// a reference, the reference is considered to be an index in this array.
	Values dgo.Array

	// The Stack of values that is used when adding nested constructs (arrays and maps)
	Stack []dgo.Array
}

// NewCollector returns a new BasicCollector instance
func NewCollector() Collector {
	hm := &BasicCollector{}
	hm.Init()
	return hm
}

// Init initializes the internal stack and reference storage
func (hm *BasicCollector) Init() {
	hm.Values = vf.MutableValues(nil)
	hm.Stack = make([]dgo.Array, 1, 8)
	hm.Stack[0] = vf.MutableValues(nil)
}

// AddArray initializes and adds a new array and then calls the function with is supposed to
// add the elements.
func (hm *BasicCollector) AddArray(cap int, doer dgo.Doer) {
	a := vf.ArrayWithCapacity(nil, cap)
	hm.Add(a)
	top := len(hm.Stack)
	hm.Stack = append(hm.Stack, a)
	doer()
	hm.Stack = hm.Stack[0:top]
}

// AddMap initializes and adds a new map and then calls the function with is supposed to
// add an even number of elements as a sequence of key, value, [key, value, ...]
func (hm *BasicCollector) AddMap(cap int, doer dgo.Doer) {
	h := vf.MutableMap(nil)
	hm.Add(h)
	a := vf.ArrayWithCapacity(nil, cap*2)
	top := len(hm.Stack)
	hm.Stack = append(hm.Stack, a)
	doer()
	hm.Stack = hm.Stack[0:top]
	h.PutAll(a.ToMap())
}

// Add adds a new value
func (hm *BasicCollector) Add(element dgo.Value) {
	hm.StackTop().Add(element)
	hm.Values.Add(element)
}

// AddRef adds the nth value of the values that has been added once again.
func (hm *BasicCollector) AddRef(ref int) {
	hm.StackTop().Add(hm.Values.Get(ref))
}

// CanDoBinary returns true
func (hm *BasicCollector) CanDoBinary() bool {
	return true
}

// CanDoTime returns true
func (hm *BasicCollector) CanDoTime() bool {
	return true
}

// CanDoComplexKeys returns true
func (hm *BasicCollector) CanDoComplexKeys() bool {
	return true
}

// StringDedupThreshold returns 0
func (hm *BasicCollector) StringDedupThreshold() int {
	return 0
}

// Value returns the last value added to this collector
func (hm *BasicCollector) Value() dgo.Value {
	return hm.Stack[0].Get(0)
}

// StackTop returns the Array at the top of the collector stack.
func (hm *BasicCollector) StackTop() dgo.Array {
	return hm.Stack[len(hm.Stack)-1]
}
