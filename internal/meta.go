package internal

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

var reflectTypeType = reflect.TypeOf((*dgo.Type)(nil)).Elem()

// metaType is the Type returned by a Type
type metaType struct {
	tp dgo.Type
}

// DefaultMetaType is the unconstrained meta type
var DefaultMetaType = &metaType{tp: DefaultAnyType}

func (t *metaType) Type() dgo.Type {
	if t.tp == nil {
		return t // type of meta type is meta type
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
		return t.tp.Assignable(ot)
	}
	return false
}

func (t *metaType) New(arg dgo.Value) dgo.Value {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`type`, 1, 1)
		arg = args.Get(0)
	}
	var tv dgo.Type
	switch arg := arg.(type) {
	case dgo.Type:
		tv = arg
	case dgo.String:
		tv = Parse(arg.GoString())
	default:
		panic(illegalArgument(`type`, `type|string`, []interface{}{arg}, 0))
	}
	if !t.Instance(tv) {
		panic(IllegalAssignment(t, tv))
	}
	return tv
}

func (t *metaType) Operator() dgo.TypeOp {
	return dgo.OpMeta
}

func (t *metaType) Operand() dgo.Type {
	return t.tp
}

func (t *metaType) ReflectType() reflect.Type {
	return reflectTypeType
}

func (t *metaType) Resolve(ap dgo.AliasProvider) {
	tp := t.tp
	t.tp = DefaultAnyType
	t.tp = ap.Replace(tp)
}

func (t *metaType) String() string {
	return TypeString(t)
}

func (t *metaType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMeta
}
