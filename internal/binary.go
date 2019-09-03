package internal

import (
	"bytes"
	"encoding/base64"

	"github.com/lyraproj/got/dgo"
	"gopkg.in/yaml.v3"
)

type (
	binary struct {
		bytes  []byte
		frozen bool
	}

	binaryType int

	exactBinaryType binary
)

const DefaultBinaryType = binaryType(0)

func (t *exactBinaryType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*exactBinaryType); ok {
		return bytes.Equal(t.bytes, ot.bytes)
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *exactBinaryType) Equals(other interface{}) bool {
	if ot, ok := other.(*exactBinaryType); ok {
		return bytes.Equal(t.bytes, ot.bytes)
	}
	return false
}

func (t *exactBinaryType) HashCode() int {
	return bytesHash(t.bytes) * 5
}

func (t *exactBinaryType) Instance(value interface{}) bool {
	if ot, ok := value.(*binary); ok {
		return bytes.Equal(t.bytes, ot.bytes)
	}
	if ot, ok := value.([]byte); ok {
		return bytes.Equal(t.bytes, ot)
	}
	return false
}

func (t *exactBinaryType) IsInstance(bs []byte) bool {
	return bytes.Equal(t.bytes, bs)
}

func (t *exactBinaryType) String() string {
	return TypeString(t)
}

func (t *exactBinaryType) Type() dgo.Type {
	return &metaType{t}
}

func (t *exactBinaryType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.IdBinaryExact
}

func (t *exactBinaryType) Value() dgo.Value {
	v := (*binary)(t)
	return v
}

func (t binaryType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case binaryType, *exactBinaryType:
		return true
	}
	return false
}

func (t binaryType) Equals(other interface{}) bool {
	_, ok := other.(binaryType)
	return ok
}

func (t binaryType) HashCode() int {
	return int(dgo.IdBinary)
}

func (t binaryType) Instance(value interface{}) bool {
	_, ok := value.(*binary)
	return ok
}

func (t binaryType) IsInstance([]byte) bool {
	return true
}

func (t binaryType) String() string {
	return TypeString(t)
}

func (t binaryType) Type() dgo.Type {
	return &metaType{t}
}

func (t binaryType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.IdBinary
}

func Binary(bs []byte) dgo.Binary {
	c := make([]byte, len(bs))
	copy(c, bs)
	return &binary{bytes: c, frozen: true}
}

func BinaryFromString(s string) dgo.Binary {
	bs, err := base64.StdEncoding.Strict().DecodeString(s)
	if err != nil {
		panic(err)
	}
	return &binary{bytes: bs, frozen: true}
}

func (v *binary) CompareTo(other interface{}) (r int, ok bool) {
	var ov binary
	ov, ok = other.(binary)
	if !ok {
		return
	}
	a := v.bytes
	b := ov.bytes
	top := len(a)
	max := len(b)
	r = 0
	if top < max {
		r = -1
		max = top
	} else if top > max {
		r = 1
		top = max
	}
	for i := 0; i < max; i++ {
		c := a[i] - b[i]
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

func (v *binary) HashCode() int {
	return bytesHash(v.bytes)
}

func (v *binary) MarshalYAML() (interface{}, error) {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!!binary`, Value: v.String()}, nil
}

func (v *binary) String() string {
	return base64.StdEncoding.Strict().EncodeToString(v.bytes)
}

func (v *binary) GoBytes() []byte {
	if v.frozen {
		c := make([]byte, len(v.bytes))
		copy(c, v.bytes)
		return c
	}
	return v.bytes
}

func (v *binary) Type() dgo.Type {
	return (*exactBinaryType)(v)
}

func bytesHash(s []byte) int {
	h := 1
	for i := range s {
		h = 31*h + int(s[i])
	}
	return h
}
