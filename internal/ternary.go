package internal

import (
	"github.com/lyraproj/dgo/dgo"
)

type (
	allOfType array

	allOfValueType array

	anyOfType array

	oneOfType array
)

// DefaultAllOfType is the unconstrained AllOf type
var DefaultAllOfType = &allOfType{}

// DefaultAnyOfType is the unconstrained AnyOf type
var DefaultAnyOfType = &anyOfType{}

// DefaultOneOfType is the unconstrained OneOf type
var DefaultOneOfType = &oneOfType{}

// AllOfType returns a type that represents all values that matches all of the included types
func AllOfType(types []dgo.Type) dgo.Type {
	l := len(types)
	switch l {
	case 0:
		// And of no types is an unconstrained type
		return DefaultAnyType
	case 1:
		return types[0]
	}
	ts := make([]dgo.Value, l)
	for i := range types {
		ts[i] = types[i]
	}
	return &allOfType{slice: ts, frozen: true}
}

func (t *allOfType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *allOfType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	ts := t.slice
	if ot, ok := other.(*allOfType); ok {
		return lessRestrictive(guard, typeSlice(ts, typeAsType), typeSlice(ot.slice, typeAsType))
	}
	if ot, ok := other.(*allOfValueType); ok {
		return lessRestrictive(guard, typeSlice(ts, typeAsType), typeSlice(ot.slice, valueAsType))
	}
	for i := range ts {
		if !Assignable(guard, ts[i].(dgo.Type), other) {
			return CheckAssignableTo(guard, other, t)
		}
	}
	return true
}

// AssignableTo returns true if the other type is assignable from all of the contained types.
func (t *allOfType) AssignableTo(guard dgo.RecursionGuard, other dgo.Type) bool {
	ts := t.slice
	for i := range ts {
		if !Assignable(guard, other, ts[i].(dgo.Type)) {
			return false
		}
	}
	return true
}

func (t *allOfType) Equals(other interface{}) bool {
	if ot, ok := other.(*allOfType); ok {
		return (*array)(t).SameValues((*array)(ot))
	}
	return false
}

func (t *allOfType) HashCode() int {
	return (*array)(t).HashCode()*7 + int(dgo.TiAllOf)
}

func (t *allOfType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *allOfType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	ts := t.slice
	for i := range ts {
		if !Instance(guard, ts[i].(dgo.Type), value) {
			return false
		}
	}
	return true
}

func (t *allOfType) Operands() dgo.Array {
	return (*array)(t)
}

func (t *allOfType) Operator() dgo.TypeOp {
	return dgo.OpAnd
}

func (t *allOfType) String() string {
	return TypeString(t)
}

func (t *allOfType) Type() dgo.Type {
	return &metaType{t}
}

func (t *allOfType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiAllOf
}

func (t *allOfValueType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *allOfValueType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	ts := t.slice
	if ot, ok := other.(*allOfValueType); ok {
		// Safe since explicit values only are assignable to themselves and excess elements
		// in other only adds to the constraint
		return (*array)(t).ContainsAll((*array)(ot))
	}
	if ot, ok := other.(*allOfType); ok {
		return lessRestrictive(guard, typeSlice(ts, valueAsType), typeSlice(ot.slice, typeAsType))
	}
	for i := range ts {
		if !Assignable(guard, ts[i].Type(), other) {
			return false
		}
	}
	return true
}

// AssignableTo returns true if the other type is assignable from all of the contained types.
func (t *allOfValueType) AssignableTo(guard dgo.RecursionGuard, other dgo.Type) bool {
	ts := t.slice
	for i := range ts {
		if !Assignable(guard, other, ts[i].Type()) {
			return false
		}
	}
	return true
}

func (t *allOfValueType) Equals(other interface{}) bool {
	if ot, ok := other.(*allOfValueType); ok {
		return (*array)(t).SameValues((*array)(ot))
	}
	return false
}

func (t *allOfValueType) HashCode() int {
	return (*array)(t).HashCode()*7 + int(dgo.TiAllOfValue)
}

func (t *allOfValueType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *allOfValueType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	ts := t.slice
	for i := range ts {
		if !Instance(guard, ts[i].Type(), value) {
			return false
		}
	}
	return true
}

func (t *allOfValueType) Operands() dgo.Array {
	return (*array)(t)
}

func (t *allOfValueType) Operator() dgo.TypeOp {
	return dgo.OpAnd
}

func (t *allOfValueType) String() string {
	return TypeString(t)
}

func (t *allOfValueType) Type() dgo.Type {
	return &metaType{t}
}

func (t *allOfValueType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiAllOfValue
}

var notAnyType = &notType{DefaultAnyType}

// AnyOfType returns a type that represents all values that matches at least one of the included types
func AnyOfType(types []dgo.Type) dgo.Type {
	l := len(types)
	switch l {
	case 0:
		// Or of no types doesn't represent any values at all
		return notAnyType
	case 1:
		return types[0]
	}
	ts := make([]dgo.Value, l)
	for i := range types {
		ts[i] = types[i]
	}
	return &anyOfType{slice: ts, frozen: true}
}

func (t *anyOfType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *anyOfType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	ts := t.slice
	if len(ts) == 0 {
		_, ok := other.(*anyOfType)
		return ok
	}
	for i := range ts {
		if Assignable(guard, ts[i].(dgo.Type), other) {
			return true
		}
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *anyOfType) AssignableTo(guard dgo.RecursionGuard, other dgo.Type) bool {
	ts := t.slice
	for i := range ts {
		if !Assignable(guard, other, ts[i].(dgo.Type)) {
			return false
		}
	}
	return len(ts) > 0
}

func (t *anyOfType) Equals(other interface{}) bool {
	if ot, ok := other.(*anyOfType); ok {
		return (*array)(t).SameValues((*array)(ot))
	}
	return false
}

func (t *anyOfType) HashCode() int {
	return (*array)(t).HashCode()*7 + int(dgo.TiAnyOf)
}

func (t *anyOfType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *anyOfType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	ts := t.slice
	for i := range ts {
		if Instance(guard, ts[i].(dgo.Type), value) {
			return true
		}
	}
	return false
}

func (t *anyOfType) Operands() dgo.Array {
	return (*array)(t)
}

func (t *anyOfType) Operator() dgo.TypeOp {
	return dgo.OpOr
}

func (t *anyOfType) String() string {
	return TypeString(t)
}

func (t *anyOfType) Type() dgo.Type {
	return &metaType{t}
}

func (t *anyOfType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiAnyOf
}

// OneOfType returns a type that represents all values that matches exactly one of the included types
func OneOfType(types []dgo.Type) dgo.Type {
	l := len(types)
	switch l {
	case 0:
		// One of no types doesn't represent any values at all
		return notAnyType
	case 1:
		return types[0]
	}
	ts := make([]dgo.Value, l)
	for i := range types {
		ts[i] = types[i]
	}
	return &oneOfType{slice: ts, frozen: true}
}

func (t *oneOfType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *oneOfType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	ts := t.slice
	if len(ts) == 0 {
		_, ok := other.(*oneOfType)
		return ok
	}
	found := false
	for i := range ts {
		if Assignable(guard, ts[i].(dgo.Type), other) {
			if found {
				found = false
				break
			}
			found = true
		}
	}
	if !found {
		found = CheckAssignableTo(guard, other, t)
	}
	return found
}

// AssignableTo returns true if the other type is assignable from all of the contained types
func (t *oneOfType) AssignableTo(guard dgo.RecursionGuard, other dgo.Type) bool {
	ts := t.slice
	for i := range ts {
		if !Assignable(guard, other, ts[i].(dgo.Type)) {
			return false
		}
	}
	return len(ts) > 0
}

func (t *oneOfType) Equals(other interface{}) bool {
	if ot, ok := other.(*oneOfType); ok {
		return (*array)(t).SameValues((*array)(ot))
	}
	return false
}

func (t *oneOfType) HashCode() int {
	return (*array)(t).HashCode()
}

func (t *oneOfType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *oneOfType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	ts := t.slice
	found := false
	for i := range ts {
		if Instance(guard, ts[i].(dgo.Type), value) {
			if found {
				// Found twice
				return false
			}
			found = true
		}
	}
	return found
}

func (t *oneOfType) Operands() dgo.Array {
	return (*array)(t)
}

func (t *oneOfType) Operator() dgo.TypeOp {
	return dgo.OpOne
}

func (t *oneOfType) String() string {
	return TypeString(t)
}

func (t *oneOfType) Type() dgo.Type {
	return &metaType{t}
}

func (t *oneOfType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiOneOf
}

func typeSlice(s []dgo.Value, fc func(dgo.Value) dgo.Type) []dgo.Type {
	ts := make([]dgo.Type, len(s))
	for i := range s {
		ts[i] = fc(s[i])
	}
	return ts
}

// lessRestrictive is true if slice 'a' is less restrictive than slice 'b'
//
// This is true if:
// 1. All types in slice a finds at least one assignable type in slice b
// 2. All types in slice b is assignable to at least one type in slice a or not less restrictive than any type in a
func lessRestrictive(guard dgo.RecursionGuard, ats, bts []dgo.Type) bool {
	matched := make([]bool, len(bts))

	// Assume that each type in 'a' has an assignable match in 'b'
nextA:
	for ia := range ats {
		ta := ats[ia]
		for ib := range bts {
			if !matched[ib] && Assignable(guard, ta, bts[ib]) {
				matched[ib] = true
				continue nextA
			}
		}
		// Check with already matched. Some types in 'b' might be matched by the same
		// type in 'a'
		for ib := range bts {
			if matched[ib] && Assignable(guard, ta, bts[ib]) {
				continue nextA
			}
		}
		// Not found in 'b', so 'b' is less restrictive
		return false
	}

	// Extra entries in 'b' are allowed provided they don't lessen the constraint
	var swg dgo.RecursionGuard
	for ib := range bts {
		if !matched[ib] {
			tb := bts[ib]
			if nil != guard && nil == swg {
				swg = guard.Swap()
			}
			for ia := range ats {
				if Assignable(swg, tb, ats[ia]) {
					return false
				}
			}
		}
	}
	return true
}
