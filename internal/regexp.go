package internal

import (
	"regexp"
	"strings"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// regexpType represents an regexp type without constraints
	regexpType int

	exactRegexpType regexp.Regexp

	regexpVal regexp.Regexp
)

// DefaultRegexpType is the unconstrained Regexp type
const DefaultRegexpType = regexpType(0)

func (t regexpType) Assignable(ot dgo.Type) bool {
	switch ot.(type) {
	case regexpType, *exactRegexpType:
		return true
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t regexpType) Equals(v interface{}) bool {
	return t == v
}

func (t regexpType) HashCode() int {
	return int(dgo.TiRegexp)
}

func (t regexpType) Instance(v interface{}) bool {
	_, ok := v.(*regexpVal)
	if !ok {
		_, ok = v.(*regexp.Regexp)
	}
	return ok
}

func (t regexpType) String() string {
	return TypeString(t)
}

func (t regexpType) Type() dgo.Type {
	return &metaType{t}
}

func (t regexpType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiRegexp
}

func (t *exactRegexpType) Assignable(other dgo.Type) bool {
	return t.Equals(other)
}

func (t *exactRegexpType) Equals(other interface{}) bool {
	if ot, ok := other.(*exactRegexpType); ok {
		return (*regexp.Regexp)(t).String() == (*regexp.Regexp)(ot).String()
	}
	return false
}

func (t *exactRegexpType) HashCode() int {
	return (*regexpVal)(t).HashCode()*31 + int(dgo.TiRegexpExact)
}

func (t *exactRegexpType) Instance(value interface{}) bool {
	if ot, ok := value.(*regexpVal); ok {
		return t.IsInstance((*regexp.Regexp)(ot))
	}
	if ot, ok := value.(*regexp.Regexp); ok {
		return t.IsInstance(ot)
	}
	return false
}

func (t *exactRegexpType) IsInstance(v *regexp.Regexp) bool {
	return (*regexp.Regexp)(t).String() == v.String()
}

func (t *exactRegexpType) String() string {
	return TypeString(t)
}

func (t *exactRegexpType) Type() dgo.Type {
	return &metaType{t}
}

func (t *exactRegexpType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiRegexpExact
}

func (t *exactRegexpType) Value() dgo.Value {
	v := (*regexpVal)(t)
	return v
}

func (v *regexpVal) GoRegexp() *regexp.Regexp {
	return (*regexp.Regexp)(v)
}

func (v *regexpVal) Equals(other interface{}) bool {
	if ot, ok := other.(*regexpVal); ok {
		return (*regexp.Regexp)(v).String() == (*regexp.Regexp)(ot).String()
	}
	if ot, ok := other.(*regexp.Regexp); ok {
		return (*regexp.Regexp)(v).String() == (ot).String()
	}
	return false
}

func (v *regexpVal) HashCode() int {
	return stringHash((*regexp.Regexp)(v).String())
}

func (v *regexpVal) String() string {
	return (*regexp.Regexp)(v).String()
}

func (v *regexpVal) Type() dgo.Type {
	return (*exactRegexpType)(v)
}

// RegexpSlashQuote converts the given string into a slash delimited string with internal slashes escaped
// and writes it on the given builder.
func RegexpSlashQuote(sb *strings.Builder, str string) {
	util.WriteByte(sb, '/')
	for _, c := range str {
		switch c {
		case '\t':
			util.WriteString(sb, `\t`)
		case '\n':
			util.WriteString(sb, `\n`)
		case '\r':
			util.WriteString(sb, `\r`)
		case '/':
			util.WriteString(sb, `\/`)
		case '\\':
			util.WriteString(sb, `\\`)
		default:
			if c < 0x20 {
				util.Fprintf(sb, `\u{%X}`, c)
			} else {
				util.WriteRune(sb, c)
			}
		}
	}
	util.WriteByte(sb, '/')
}
