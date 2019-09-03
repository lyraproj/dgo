package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// regexpType represents an regexp type without constraints
	regexpType int

	exactRegexpType regexp.Regexp

	Regexp regexp.Regexp
)

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
	return int(dgo.IdRegexp)
}

func (t regexpType) Instance(v interface{}) bool {
	_, ok := v.(*Regexp)
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
	return dgo.IdRegexp
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
	return (*Regexp)(t).HashCode()*31 + int(dgo.IdRegexpExact)
}

func (t *exactRegexpType) Instance(value interface{}) bool {
	if ot, ok := value.(*Regexp); ok {
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
	return dgo.IdRegexpExact
}

func (t *exactRegexpType) Value() dgo.Value {
	v := (*Regexp)(t)
	return v
}

func (v *Regexp) GoRegexp() *regexp.Regexp {
	return (*regexp.Regexp)(v)
}

func (v *Regexp) Equals(other interface{}) bool {
	if ot, ok := other.(*Regexp); ok {
		return (*regexp.Regexp)(v).String() == (*regexp.Regexp)(ot).String()
	}
	if ot, ok := other.(*regexp.Regexp); ok {
		return (*regexp.Regexp)(v).String() == (ot).String()
	}
	return false
}

func (v *Regexp) HashCode() int {
	return stringHash((*regexp.Regexp)(v).String())
}

func (v *Regexp) String() string {
	return (*regexp.Regexp)(v).String()
}

func (v *Regexp) Type() dgo.Type {
	return (*exactRegexpType)(v)
}

func RegexpSlashQuote(sb *strings.Builder, str string) {
	sb.WriteByte('/')
	for _, c := range str {
		switch c {
		case '\t':
			sb.WriteString(`\t`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '/':
			sb.WriteString(`\/`)
		case '\\':
			sb.WriteString(`\\`)
		default:
			if c < 0x20 {
				_, _ = fmt.Fprintf(sb, `\u{%X}`, c)
			} else {
				sb.WriteRune(c)
			}
		}
	}
	sb.WriteByte('/')
}
