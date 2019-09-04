package internal

import "github.com/lyraproj/dgo/dgo"

type notType struct {
	Negated dgo.Type
}

// DefaultNotType is the unconstrained Not type
var DefaultNotType = &notType{DefaultAnyType}

// NotType returns a type that represents all values that are not represented by the given type
func NotType(t dgo.Type) dgo.Type {
	// Avoid double negation
	if nt, ok := t.(*notType); ok {
		return nt.Negated
	}
	return &notType{Negated: t}
}

func (t *notType) Equals(other interface{}) bool {
	if ot, ok := other.(*notType); ok {
		return t.Negated.Equals(ot.Negated)
	}
	return false
}

func (t *notType) HashCode() int {
	return 1579 + t.Negated.HashCode()
}

func (t *notType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case *notType:
		// Reverse order of Negated test
		return ot.Negated.Assignable(t.Negated)
	case *anyOfType:
		ts := ot.slice
		for i := range ts {
			if t.Assignable(ts[i].(dgo.Type)) {
				return true
			}
		}
		return false
	case *allOfType:
		ts := ot.slice
		for i := range ts {
			if !t.Assignable(ts[i].(dgo.Type)) {
				return false
			}
		}
		return true
	case *oneOfType:
		f := false
		ts := ot.slice
		for i := range ts {
			if t.Assignable(ts[i].(dgo.Type)) {
				if f {
					return false
				}
				f = true
			}
		}
		return f
	default:
		return !t.Negated.Assignable(other)
	}
}

func (t *notType) Instance(value interface{}) bool {
	return !t.Negated.Instance(value)
}

func (t *notType) Operand() dgo.Type {
	return t.Negated
}

func (t *notType) Operator() dgo.TypeOp {
	return dgo.OpNot
}

func (t *notType) String() string {
	return TypeString(t)
}

func (t *notType) Type() dgo.Type {
	return &metaType{t}
}

func (t *notType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNot
}
