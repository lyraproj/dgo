package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"gopkg.in/yaml.v3"
)

type (
	binary struct {
		bytes  []byte
		frozen bool
	}

	binaryType struct {
		min int
		max int
	}
)

// DefaultBinaryType is the unconstrained Binary type
var DefaultBinaryType = &binaryType{0, math.MaxInt64}

// BinaryType returns a new dgo.BinaryType. It can be called with two optional integer arguments denoting
// the min and max length of the binary. If only one integer is given, it represents the min length.
func BinaryType(args ...interface{}) dgo.BinaryType {
	switch len(args) {
	case 0:
		return DefaultBinaryType
	case 1:
		if a0, ok := Value(args[0]).(dgo.Integer); ok {
			return SizedBinaryType(int(a0.GoInt()), int(a0.GoInt()))
		}
		panic(illegalArgument(`BinaryType`, `Integer`, args, 0))
	case 2:
		if a0, ok := Value(args[0]).(dgo.Integer); ok {
			var a1 dgo.Integer
			if a1, ok = Value(args[1]).(dgo.Integer); ok {
				return SizedBinaryType(int(a0.GoInt()), int(a1.GoInt()))
			}
			panic(illegalArgument(`BinaryType`, `Integer`, args, 1))
		}
		panic(illegalArgument(`BinaryType`, `Integer`, args, 0))
	}
	panic(fmt.Errorf(`illegal number of arguments for BinaryType. Expected 0 - 2, got %d`, len(args)))
}

// SizedBinaryType returns a BinaryType that is constrained to binaries whose size is within the
// inclusive range given by min and max.
func SizedBinaryType(min, max int) dgo.BinaryType {
	if min < 0 {
		min = 0
	}
	if max < min {
		tmp := max
		max = min
		min = tmp
	}
	if min == 0 && max == math.MaxInt64 {
		return DefaultBinaryType
	}
	return &binaryType{min: min, max: max}
}

func (t *binaryType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*binaryType); ok {
		return t.min <= ot.min && t.max >= ot.max
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *binaryType) Equals(other interface{}) bool {
	if ob, ok := other.(*binaryType); ok {
		return *t == *ob
	}
	return false
}

func (t *binaryType) HashCode() int {
	h := int(dgo.TiBinary)
	if t.min > 0 {
		h = h*31 + t.min
	}
	if t.max < math.MaxInt64 {
		h = h*31 + t.max
	}
	return h
}

func (t *binaryType) Instance(value interface{}) bool {
	if ov, ok := value.(*binary); ok {
		return t.IsInstance(ov.bytes)
	}
	if ov, ok := value.([]byte); ok {
		return t.IsInstance(ov)
	}
	return false
}

func (t *binaryType) IsInstance(v []byte) bool {
	l := len(v)
	return t.min <= l && l <= t.max
}

func (t *binaryType) Max() int {
	return t.max
}

func (t *binaryType) Min() int {
	return t.min
}

var reflectBinaryType = reflect.TypeOf([]byte{})

func (t *binaryType) ReflectType() reflect.Type {
	return reflectBinaryType
}

func (t *binaryType) String() string {
	return TypeString(t)
}

func (t *binaryType) Type() dgo.Type {
	return &metaType{t}
}

func (t *binaryType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBinary
}

func (t *binaryType) Unbounded() bool {
	return t.min == 0 && t.max == math.MaxInt64
}

// Binary creates a new Binary based on the given slice. If frozen is true, the
// binary will be immutable and contain a copy of the slice, otherwise the slice
// is simply wrapped and modifications to its elements will also modify the binary.
func Binary(bs []byte, frozen bool) dgo.Binary {
	if frozen {
		c := make([]byte, len(bs))
		copy(c, bs)
		bs = c
	}
	return &binary{bytes: bs, frozen: frozen}
}

// BinaryFromString creates a new Binary from the base64 encoded string
func BinaryFromString(s string) dgo.Binary {
	bs, err := base64.StdEncoding.Strict().DecodeString(s)
	if err != nil {
		panic(err)
	}
	return &binary{bytes: bs, frozen: true}
}

func (v *binary) Copy(frozen bool) dgo.Binary {
	if frozen && v.frozen {
		return v
	}
	cp := make([]byte, len(v.bytes))
	copy(cp, v.bytes)
	return &binary{bytes: cp, frozen: frozen}
}

func (v *binary) CompareTo(other interface{}) (r int, ok bool) {
	var b []byte
	var ob *binary
	if ob, ok = other.(*binary); ok {
		if v == ob {
			return 0, true
		}
		b = ob.bytes
	} else {
		b, ok = other.([]byte)
		if !ok {
			if other == nil || other == Nil {
				return 1, true
			}
			return 0, false
		}
	}
	a := v.bytes
	top := len(a)
	max := len(b)
	r = 0
	if top < max {
		r = -1
		max = top
	} else if top > max {
		r = 1
	}
	for i := 0; i < max; i++ {
		c := int(a[i]) - int(b[i])
		if c != 0 {
			if c > 0 {
				r = 1
			} else {
				r = -1
			}
			break
		}
	}
	return
}

func (v *binary) Equals(other interface{}) bool {
	if ot, ok := other.(*binary); ok {
		return bytes.Equal(v.bytes, ot.bytes)
	}
	if ot, ok := other.([]byte); ok {
		return bytes.Equal(v.bytes, ot)
	}
	return false
}

func (v *binary) Freeze() {
	if !v.frozen {
		bs := v.bytes
		v.bytes = make([]byte, len(bs))
		copy(v.bytes, bs)
		v.frozen = true
	}
}

func (v *binary) Frozen() bool {
	return v.frozen
}

func (v *binary) FrozenCopy() dgo.Value {
	if !v.frozen {
		cs := make([]byte, len(v.bytes))
		copy(cs, v.bytes)
		return &binary{bytes: cs, frozen: true}
	}
	return v
}

func (v *binary) GoBytes() []byte {
	if v.frozen {
		c := make([]byte, len(v.bytes))
		copy(c, v.bytes)
		return c
	}
	return v.bytes
}

func (v *binary) HashCode() int {
	return bytesHash(v.bytes)
}

func (v *binary) MarshalYAML() (interface{}, error) {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!!binary`, Value: v.String()}, nil
}

func (v *binary) ReflectTo(value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		x := reflect.New(reflectBinaryType)
		x.Elem().SetBytes(v.GoBytes())
		value.Set(x)
	case reflect.Slice:
		value.SetBytes(v.GoBytes())
	default:
		value.Set(reflect.ValueOf(v.GoBytes()))
	}
}

func (v *binary) String() string {
	return base64.StdEncoding.Strict().EncodeToString(v.bytes)
}

func (v *binary) Type() dgo.Type {
	l := len(v.bytes)
	return &binaryType{l, l}
}

func bytesHash(s []byte) int {
	h := 1
	for i := range s {
		h = 31*h + int(s[i])
	}
	return h
}
