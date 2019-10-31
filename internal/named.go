package internal

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

type (
	named struct {
		name      string
		ctor      dgo.Constructor
		extractor dgo.InitArgExtractor
		implType  reflect.Type
		ifdType   reflect.Type
	}

	exactNamed struct {
		privNamedType
		value dgo.Value
	}

	privNamedType interface {
		dgo.NamedType

		assignableType() reflect.Type
	}
)

var namedTypes = sync.Map{}

// RemoveNamedType removes a named type from the global type registry. It is primarily intended for
// testing purposes.
func RemoveNamedType(name string) {
	namedTypes.Delete(name)
}

// NewNamedType registers a new named type under the given name with the global type registry. The method panics
// if a type has already been registered with the same name.
func NewNamedType(name string, ctor dgo.Constructor, extractor dgo.InitArgExtractor, implType, ifdType reflect.Type) dgo.NamedType {
	t, loaded := namedTypes.LoadOrStore(name, &named{name: name, ctor: ctor, extractor: extractor, implType: implType, ifdType: ifdType})
	if loaded {
		panic(fmt.Errorf(`attempt to redefine named type '%s'`, name))
	}
	return t.(*named)
}

// ExactNamedType returns the exact NamedType that represents the given value.
func ExactNamedType(namedType dgo.NamedType, value dgo.Value) dgo.NamedType {
	return &exactNamed{privNamedType: namedType.(privNamedType), value: value}
}

// NamedType returns the type with the given name from the global type registry. The function returns
// nil if no type has been registered under the given name.
func NamedType(name string) dgo.NamedType {
	if t, ok := namedTypes.Load(name); ok {
		return t.(*named)
	}
	return nil
}

// NamedTypeFromReflected returns the named type for the reflected implementation type from the global type
// registry. The function returns nil if no such type has been registered.
func NamedTypeFromReflected(rt reflect.Type) dgo.NamedType {
	var t dgo.NamedType
	namedTypes.Range(func(k, v interface{}) bool {
		nt := v.(*named)
		if rt == nt.implType {
			t = nt
			return false
		}
		return true
	})
	return t
}

func (t *named) assignableType() reflect.Type {
	if t.ifdType != nil {
		return t.ifdType
	}
	return t.implType
}

func (t *named) Assignable(other dgo.Type) bool {
	if ot, ok := other.(privNamedType); ok {
		return ot.ReflectType().AssignableTo(t.assignableType())
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *named) New(arg dgo.Value) dgo.Value {
	if t.ctor == nil {
		panic(fmt.Errorf(`creating new instances of %s is not possible`, t.name))
	}
	return t.ctor(arg)
}

func (t *named) Equals(other interface{}) bool {
	if ot, ok := other.(*named); ok {
		return t.name == ot.name
	}
	return false
}

func (t *named) HashCode() int {
	return util.StringHash(t.name)*7 + int(dgo.TiNamed)
}

func (t *named) ExtractInitArg(value dgo.Value) dgo.Value {
	if t.extractor != nil {
		return t.extractor(value)
	}
	panic(fmt.Errorf(`creating new instances of %s is not possible`, t.name))
}

func (t *named) Instance(value interface{}) bool {
	return reflect.TypeOf(value).AssignableTo(t.assignableType())
}

func (t *named) Name() string {
	return t.name
}

func (t *named) ReflectType() reflect.Type {
	return t.implType
}

func (t *named) String() string {
	return TypeString(t)
}

func (t *named) Type() dgo.Type {
	return &metaType{t}
}

func (t *named) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNamed
}

func (t *named) ValueString(v dgo.Value) string {
	b := NewERPIndenter(``)
	util.WriteString(b, t.name)
	if t.extractor != nil {
		switch t.extractor(v).(type) {
		case dgo.Array, dgo.Map, dgo.String:
		default:
			b.AppendRune(' ')
		}
		b.AppendValue(t.extractor(v))
	}
	return b.String()
}

func (t *exactNamed) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*exactNamed); ok {
		return t.value.Equals(ot.value)
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *exactNamed) Equals(other interface{}) bool {
	if ot, ok := other.(*exactNamed); ok {
		return t.value.Equals(ot.value)
	}
	return false
}

func (t *exactNamed) HashCode() int {
	return t.value.HashCode()*7 + int(dgo.TiNamedExact)
}

func (t *exactNamed) Generic() dgo.Type {
	return t.privNamedType
}

func (t *exactNamed) Instance(other interface{}) bool {
	return t.value.Equals(other)
}

func (t *exactNamed) String() string {
	return TypeString(t)
}

func (t *exactNamed) Type() dgo.Type {
	return &metaType{t}
}

func (t *exactNamed) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNamedExact
}

func (t *exactNamed) Value() dgo.Value {
	return t.value
}
