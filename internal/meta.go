package internal

import (
	"github.com/lyraproj/dgo/dgo"
)

// metaType is the Type returned by a Type
type metaType struct {
	tp dgo.Type
}

func (t *metaType) Type() dgo.Type {
	if t.tp == nil {
		return t
	}
	return &metaType{nil} // Short circuit meta chain
}

func (t *metaType) Assignable(ot dgo.Type) bool {
	if mt, ok := ot.(*metaType); ok {
		if t.tp == nil {
			// Only MetaTypeType is assignable to MetaTypeType
			return mt.tp == nil
		}
		return t.tp.Equals(mt.tp)
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t *metaType) Equals(v interface{}) bool {
	if mt, ok := v.(*metaType); ok {
		if t.tp == nil {
			return mt.tp == nil
		}
		return t.tp.Equals(mt.tp)
	}
	return false
}

func (t *metaType) HashCode() int {
	return int(dgo.TiMeta)*1321 + t.tp.HashCode()
}

func (t *metaType) Instance(v interface{}) bool {
	if ot, ok := v.(dgo.Type); ok {
		if t.tp == nil {
			// MetaTypeType
			_, ok = ot.(*metaType)
			return ok
		}
		return t.tp.Equals(ot)
	}
	return false
}

func (t *metaType) Operator() dgo.TypeOp {
	return dgo.OpMeta
}

func (t *metaType) Operand() dgo.Type {
	return t.tp
}

func (t *metaType) Resolve(ap dgo.AliasProvider) {
	t.tp = ap.Replace(t.tp)
}

func (t *metaType) String() string {
	return TypeString(t)
}

func (t *metaType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMeta
}
