package internal

import (
	"fmt"

	"github.com/lyraproj/got/dgo"
)

type (
	typeError struct {
		expected dgo.Type
		actual   dgo.Type
	}

	sizeError struct {
		sizedType     dgo.Type
		attemptedSize int
	}
)

func (v *typeError) Error() string {
	return fmt.Sprintf("a value of type %s cannot be assigned to type %s", TypeString(v.actual), TypeString(v.expected))
}

func (v *typeError) Equals(other interface{}) bool {
	if ov, ok := other.(*typeError); ok {
		return v.expected.Equals(ov.expected) && v.actual.Equals(ov.actual)
	}
	return false
}

func (v *typeError) HashCode() int {
	return v.expected.HashCode()*31 + v.actual.HashCode()
}

func (v *typeError) String() string {
	return v.Error()
}

func (v *typeError) Type() dgo.Type {
	return DefaultErrorType
}

func (v *sizeError) Equals(other interface{}) bool {
	if ov, ok := other.(*sizeError); ok {
		return v.sizedType.Equals(ov.sizedType) && v.attemptedSize == ov.attemptedSize
	}
	return false
}

func (v *sizeError) HashCode() int {
	return v.sizedType.HashCode()*7 + v.attemptedSize
}

func (v *sizeError) Error() string {
	return fmt.Sprintf("size constraint violation on type %s when attempting resize to %d", TypeString(v.sizedType), v.attemptedSize)
}

func (v *sizeError) String() string {
	return v.Error()
}

func (v *sizeError) Type() dgo.Type {
	return DefaultErrorType
}

func IllegalAssignment(t dgo.Type, v dgo.Value) dgo.Value {
	return &typeError{t, v.Type()}
}

func IllegalSize(t dgo.Type, sz int) dgo.Value {
	return &sizeError{t, sz}
}
