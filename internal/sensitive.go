package internal

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	sensitive struct {
		value dgo.Value
	}

	sensitiveType struct {
		wrapped dgo.Type
	}
)

// DefaultSensitiveType is the unconstrained Sensitive type
var DefaultSensitiveType = &sensitiveType{wrapped: DefaultAnyType}

func SensitiveType(wrappedType dgo.Type) dgo.Type {
	return &sensitiveType{wrapped: wrappedType}
}

func (t *sensitiveType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*sensitiveType); ok {
		return t.wrapped.Assignable(ot.wrapped)
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *sensitiveType) Equals(other interface{}) bool {
	if ot, ok := other.(*sensitiveType); ok {
		return t.wrapped.Equals(ot.wrapped)
	}
	return false
}

func (t *sensitiveType) HashCode() int {
	return int(dgo.TiSensitive)*31 + t.wrapped.HashCode()
}

func (t *sensitiveType) Instance(value interface{}) bool {
	if ov, ok := value.(*sensitive); ok {
		return t.wrapped.Instance(ov.value)
	}
	return false
}

var reflectSensitiveType = reflect.TypeOf((*dgo.Sensitive)(nil)).Elem()

func (t *sensitiveType) ReflectType() reflect.Type {
	return reflectSensitiveType
}

func (t *sensitiveType) Operand() dgo.Type {
	return t.wrapped
}

func (t *sensitiveType) Operator() dgo.TypeOp {
	return dgo.OpSensitive
}

func (t *sensitiveType) String() string {
	return TypeString(t)
}

func (t *sensitiveType) Type() dgo.Type {
	return &metaType{t}
}

func (t *sensitiveType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiSensitive
}

// Sensitive creates a new Sensitive that wraps the given value
func Sensitive(v dgo.Value) dgo.Sensitive {
	return &sensitive{v}
}

func (v *sensitive) Equals(other interface{}) bool {
	if ov, ok := other.(*sensitive); ok {
		return v.value.Equals(ov.value)
	}
	return false
}

func (v *sensitive) Freeze() {
	if f, ok := v.value.(dgo.Freezable); ok {
		f.Freeze()
	}
}

func (v *sensitive) Frozen() bool {
	if f, ok := v.value.(dgo.Freezable); ok {
		return f.Frozen()
	}
	return true
}

func (v *sensitive) FrozenCopy() dgo.Value {
	if f, ok := v.value.(dgo.Freezable); ok && !f.Frozen() {
		return &sensitive{f.FrozenCopy()}
	}
	return v
}

func (v *sensitive) HashCode() int {
	return v.value.HashCode() * 7
}

func (v *sensitive) String() string {
	return `sensitive [value redacted]`
}

func (v *sensitive) Type() dgo.Type {
	return &sensitiveType{wrapped: Generic(v.value.Type())}
}

func (v *sensitive) Unwrap() dgo.Value {
	return v.value
}
