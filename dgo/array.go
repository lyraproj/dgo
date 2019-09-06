package dgo

import (
	"encoding/json"

	"github.com/lyraproj/dgo/util"
	"gopkg.in/yaml.v3"
)

type (
	// Doer is performs some task on behalf of a caller
	Doer func(value Value)

	// Mapper maps a value to another value
	Mapper func(value Value) interface{}

	// DoWithIndex performs some task on behalf of an indexed caller
	DoWithIndex func(value Value, index int)

	// Predicate returns true of false based on the given value
	Predicate func(value Value) bool

	// Freezable is implemented by objects that might be mutable but can present themselves in an immutable form
	Freezable interface {

		// Freeze makes a mutable object immutable. It does nothing when called on an immutable object.
		//
		// The freeze operation is recursive so collections containing objects that implement this
		// interface will receive a call to this method.
		//
		// All Arrays returned by methods on a frozen Array will also be frozen with the exception of
		// the method Copy when called with frozen = false.
		Freeze()

		// Frozen returns true if the object is frozen, false otherwise
		Frozen() bool

		// FrozenCopy checks if the receiver is frozen. If it is, it returned. If not, a frozen copy
		// of the receiver is returned.
		FrozenCopy() Value
	}

	// Iterable enables the implementor to express how iteration is performed over contained elements
	Iterable interface {
		Value

		// Each calls the given function once for each value of this Iterable.
		Each(doer Doer)
	}

	// Array represents an immutable list of values. It is ensured that an Array never contains any
	// unset positions. Uninitialized positions are always converted to the Nil value.
	Array interface {
		Iterable
		Freezable
		Comparable
		util.Indentable
		json.Marshaler
		json.Unmarshaler
		yaml.Marshaler
		yaml.Unmarshaler

		// Add adds the given value to the end of this array. It panics if the receiver is frozen.
		Add(val interface{})

		// AddAll adds the elements of the given Array to the end of this array. It panics if the receiver is frozen.
		AddAll(values Array)

		// AddValues adds the given values to the end of this array. It panics if the receiver is frozen.
		AddValues(values ...interface{})

		// All returns true if the predicate returns true for all values of this Array.
		All(predicate Predicate) bool

		// Any returns true if the predicate returns true for any value of this Array.
		Any(predicate Predicate) bool

		// AppendToSlice appends all values of this array to the given slice and returns the
		// result of the append.
		AppendToSlice([]Value) []Value

		// ContainsAll returns true if this Array contains all elements of the other Array
		ContainsAll(other Array) bool

		// Copy returns a copy of the Array. The copy is frozen or mutable depending on
		// the given argument. A request to create a frozen copy of an already frozen Array
		// is a no-op that returns the receiver.
		//
		// If a frozen copy is requested from a non-frozen Array, then all non-frozen elements
		// will be copied and frozen recursively.
		//
		// A Copy of an array that contains back references to itself will result in a stack
		// overflow panic.
		Copy(frozen bool) Array

		// EachWithIndex calls the given function once for each value of this Array. The index of
		// the current value is provided in the call.
		EachWithIndex(doer DoWithIndex)

		// Get returns the value at the given position. A negative position or a position
		// that is greater or equal to the length of the array will result in a panic.
		Get(position int) Value

		// GoSlice returns the internal slice or, in case the Array is frozen, a copy
		// of the internal slice.
		GoSlice() []Value

		// IndexOf returns the index of the given value in this Array. The index is determined
		// by calling the Equals method on each element until a matching element is found. The
		// method returns -1 to indicate not found.
		IndexOf(value interface{}) int

		// Insert inserts the given value at the given position and moves all values after that position
		// one step forward. The method panics if the receiver is frozen.
		Insert(pos int, val interface{})

		// Len returns the number of values in this Array.
		Len() int

		// Map returns a new equally sized Array where each value has been replaced using the
		// given mapper function.
		Map(mapper Mapper) Array

		// MapTo is like Map but with the added ability to constrain the created Array with a given
		// ArrayType.
		MapTo(t ArrayType, mapper Mapper) Array

		// One returns true if the predicate returns true for exactly one value of this Array.
		One(predicate Predicate) bool

		// Pop removes and returns the last element of the array together with a boolean indicating
		// if the pop was possible (i.e. if the array had any elements)
		Pop() (Value, bool)

		// Reduce calls the given reducer function once for each value in the Array. The first argument,
		// the memo, is the result of the previous call to the reducer function, or when the iteration
		// starts, the memo given to this method. The second argument is the current value.
		//
		// Reduce returns the last computed memo. For an empty array, this will be the initial memo.
		Reduce(memo interface{}, reductor func(memo Value, elem Value) interface{}) Value

		// Reject returns a new Array where all values for which the predicate returned true
		// has been removed.
		Reject(predicate Predicate) Array

		// Remove removes the value at the given position and moves all values after that position one
		// step back. The removed value fis returned. The method panics if the receiver is frozen
		Remove(pos int) Value

		// RemoveValue removes the first found occurrence of the given value and moves all values after its position one
		// step back. The method returns true if the removal was performed and false when the value wasn't found.  The
		// method panics if the receiver is frozen.
		RemoveValue(value interface{}) bool

		// SameValues returns true if this Array is the same size as the other Array and contains all of its values
		SameValues(other Array) bool

		// Select returns a new Array where only values for which the predicate returned true
		// are included.
		Select(predicate Predicate) Array

		// Set replaces the given value at the given position and returns the old value for the position.
		// The method panics if the receiver is frozen
		Set(pos int, val interface{}) Value

		// SetType sets the type for this Array to the given ArrayType. The Array must be mutable and an instance of the given type
		SetType(t ArrayType)

		// Sort returns a new Array with all elements sorted using their natural order. The method
		// will panic unless all elements implement the Comparable interface
		Sort() Array

		// ToMap returns this Array as a Map. The first and second elements of the array becomes the first key and
		// value association of the Map, the third and fourth element becomes the second association, and so on. The
		// association will have a Nil value if the Array has an uneven number of elements. The frozen status of this
		// array is inherited by the new Map.
		ToMap() Map

		// ToMapFromEntries assumes that all elements of this Array are either Arrays with two elements or MapEntries.
		// If it is the former, MapEntries are created using the two elements of the Arrays as the key and value. A
		// new Map is created that will contained all the MapEntries. The frozen status of this array is inherited
		// by the new Map.
		ToMapFromEntries() (Map, bool)

		// Unique returns a new Array where all duplicate values have been removed
		Unique() Array

		// With appends the given value to a copy of this Array and returns the result.
		With(value interface{}) Array

		// WithAll appends the elements of the given Array to a copy of this array and returns the resulting Array
		WithAll(values Array) Array

		// WithValues appends the given values to a copy of this array and returns the resulting Array
		WithValues(values ...interface{}) Array
	}

	// ArrayType is implemented by types representing implementations of the Array value
	ArrayType interface {
		SizedType
		// ElementType returns the type of the elements for instances of this type
		ElementType() Type
	}

	// TupleType describes an array with a fixed set of elements where each element must conform to a specific type.
	TupleType interface {
		ArrayType

		ElementTypes() Array
	}
)
