package internal

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

// States:
const (
	exListComma   = iota // Expect comma or end of array
	exParamsComma        // Expect comma or end of parameter list
	exLeftBracket
	exRightBracket
	exRightParen
	exIntOrFloat
	exTypeExpression
	exEnd
)

type optionalValue struct {
	value dgo.Value
}

func (o *optionalValue) String() string {
	return o.value.String() + `?`
}

func (o *optionalValue) Type() dgo.Type {
	return o.value.Type()
}

func (o *optionalValue) Equals(other interface{}) bool {
	if ov, ok := other.(*optionalValue); ok {
		return o.value.Equals(ov.value)
	}
	return false
}

func (o *optionalValue) HashCode() int {
	return o.value.HashCode() * 3
}

func expect(state int) (s string) {
	switch state {
	case exParamsComma:
		s = `one of ',' or ']'`
	case exLeftBracket:
		s = `'['`
	case exRightBracket:
		s = `']'`
	case exListComma:
		s = `one of ',' or '}'`
	case exRightParen:
		s = `')'`
	case exIntOrFloat:
		s = `an literal integer or a float`
	case exTypeExpression:
		s = `a type expression`
	case exEnd:
		s = `end of expression`
	}
	return
}

func badSyntax(t *token, state int) error {
	var ts string
	if t.i == 0 {
		ts = `EOF`
	} else {
		ts = t.String()
	}
	return fmt.Errorf(`expected %s, got %s`, expect(state), ts)
}

type parser struct {
	d  []dgo.Value
	sr *util.StringReader
	pe *token
	lt *token
}

// Parse calls ParseFile with the string "<dgo type expression>" as the fileName
func Parse(content string) dgo.Type {
	return ParseFile(``, content)
}

// ParseFile parses the given content into a px.Value. The content must be a string representation
// of a Puppet literal or a type assignment expression. Valid literals are float, integer, string,
// boolean, undef, array, hash, parameter lists, and type expressions.
//
// Double quoted strings containing interpolation expressions will be parsed into a string verbatim
// without resolving the interpolations.
func ParseFile(fileName, content string) dgo.Type {
	p := &parser{sr: util.NewStringReader(content)}

	defer func() {
		if r := recover(); r != nil {
			es := r
			if err, ok := r.(error); ok {
				es = err.Error()
			}
			tl := 1
			if p.lt != nil && p.lt.s != `` {
				tl = len(p.lt.s)
			}
			fn := ``
			if fileName != `` {
				fn = fmt.Sprintf(`file: %s, `, fileName)
			}
			ln := ``
			if fileName != `` || p.sr.Line() > 1 {
				ln = fmt.Sprintf(`line: %d, `, p.sr.Line())
			}
			panic(fmt.Sprintf("%s: (%s%scolumn: %d)", es, fn, ln, p.sr.Column()-tl))
		}
	}()
	p.parse(p.nextToken())
	return p.popLastType()
}

func (p *parser) peekToken() *token {
	if p.pe == nil {
		p.pe = p.nextToken()
	}
	return p.pe
}

func (p *parser) nextToken() *token {
	var t *token
	if p.pe != nil {
		t = p.pe
		p.pe = nil
	} else {
		t = nextToken(p.sr)
	}
	p.lt = t
	return t
}

func (p *parser) parse(t *token) {
	p.anyOf(t)
	tk := p.nextToken()
	if tk.i != end {
		panic(badSyntax(tk, exEnd))
	}
}

func (p *parser) list() {
	szp := len(p.d)
	for {
		t := p.nextToken()
		if t.i == '}' {
			// Right bracket instead of element indicates an empty array or an extraneous comma. Both are OK
			break
		}
		p.arrayElement(t)
		t = p.nextToken()
		if t.i == '}' {
			break
		}
		if t.i != ',' {
			panic(badSyntax(t, exListComma))
		}
	}

	as := p.d[szp:]
	var tv dgo.Value
	if len(as) > 0 {
		ar := &array{slice: as}
		if _, ok := as[0].(dgo.MapEntry); ok {
			var m dgo.Map
			if m, ok = ar.ToMapFromEntries(); ok {
				m.Each(func(e dgo.MapEntry) {
					hn := e.(*hashNode)
					et := &entryType{}
					k := hn.key
					if kt, ok := k.(dgo.Type); ok {
						et.key = kt
					} else {
						et.key = k.Type()
					}
					v := hn.value
					et.required = true
					if ov, optional := v.(*optionalValue); optional {
						et.required = false
						v = ov.value
					}
					if vt, ok := v.(dgo.Type); ok {
						et.value = vt
					} else {
						et.value = v.Type()
					}
					hn.value = et
				})
				tv = &structType{additional: false, entries: m.(*hashMap)}
			}
		}
		if tv == nil {
			ar = ar.Copy(false).(*array)
			// Convert literal values to types and create a tupleType
			as = ar.slice
			for i := range as {
				v := as[i]
				if _, ok := v.(dgo.Type); !ok {
					as[i] = v.Type()
				}
			}
			tv = (*tupleType)(ar)
		}
	} else {
		tv = &array{}
	}
	p.d = append(p.d[:szp], tv)
}

func (p *parser) params() {
	szp := len(p.d)
	for {
		t := p.nextToken()
		if t.i == ']' {
			// Right bracket instead of element indicates an empty array or an extraneous comma. Both are OK
			break
		}
		p.arrayElement(t)
		t = p.nextToken()
		if t.i == ']' {
			break
		}
		if t.i != ',' {
			panic(badSyntax(t, exParamsComma))
		}
	}
	as := p.d[szp:]
	var tv dgo.Value
	if len(as) > 0 {
		tv = (&array{slice: as}).Copy(false)
	} else {
		tv = &array{}
	}
	p.d = append(p.d[:szp], tv)
}

func (p *parser) arrayElement(t *token) {
	p.anyOf(t)
	optional := p.peekToken().i == '?'
	if optional {
		p.nextToken()
	}
	if p.peekToken().i == ':' {
		// Map entry
		p.nextToken()
		key := p.popLast()
		p.anyOf(p.nextToken())
		val := p.popLast()
		if optional {
			val = &optionalValue{val}
		}
		p.d = append(p.d, &hashNode{key: key, value: val})
	}
}

func (p *parser) anyOf(t *token) {
	p.oneOf(t)
	if p.peekToken().i == '|' {
		szp := len(p.d) - 1
		for {
			p.nextToken()
			p.oneOf(p.nextToken())
			if p.peekToken().i != '|' {
				p.d = append(p.d[:szp], &anyOfType{slice: allTypes(p.d[szp:]), frozen: true})
				break
			}
		}
	}
}

func (p *parser) oneOf(t *token) {
	p.allOf(t)
	if p.peekToken().i == '^' {
		szp := len(p.d) - 1
		for {
			p.nextToken()
			p.allOf(p.nextToken())
			if p.peekToken().i != '^' {
				p.d = append(p.d[:szp], &oneOfType{slice: allTypes(p.d[szp:]), frozen: true})
				break
			}
		}
	}
}

func (p *parser) allOf(t *token) {
	p.unary(t)
	if p.peekToken().i == '&' {
		szp := len(p.d) - 1
		for {
			p.nextToken()
			p.unary(p.nextToken())
			if p.peekToken().i != '&' {
				p.d = append(p.d[:szp], &allOfType{slice: allTypes(p.d[szp:]), frozen: true})
				break
			}
		}
	}
}

func (p *parser) unary(t *token) {
	negate := false
	if t.i == '!' {
		negate = true
		t = p.nextToken()
	}
	p.typeExpression(t)
	if negate {
		p.d = append(p.d, NotType(p.popLastType()))
	}
}

func (p *parser) typeExpression(t *token) {
	var tp dgo.Value
	switch t.i {
	case '{':
		p.list()
		return
	case '(':
		p.anyOf(p.nextToken())
		n := p.nextToken()
		if n.i != ')' {
			panic(badSyntax(n, exRightParen))
		}
		return
	case '[':
		p.params()
		params := p.popLast().(*array)
		p.typeExpression(p.nextToken())
		params.Insert(0, p.popLastType())
		tp = ArrayType(sliceToInterfaces(params)...)
	case integer:
		i := tokenInt(t)
		n := p.peekToken()
		if n.i == dotdot {
			p.nextToken()
			n = p.peekToken()
			switch n.i {
			case integer:
				p.nextToken()
				tp = IntegerRangeType(i, tokenInt(n))
			case float:
				p.nextToken()
				tp = FloatRangeType(float64(i), tokenFloat(n))
			default:
				tp = IntegerRangeType(i, math.MaxInt64) // Unbounded at upper end
			}
		} else {
			tp = Integer(i)
		}
	case float:
		f := tokenFloat(t)
		n := p.peekToken()
		if n.i == dotdot {
			p.nextToken()
			n = p.peekToken()
			if n.i == integer || n.i == float {
				p.nextToken()
				tp = FloatRangeType(f, tokenFloat(n))
			} else {
				p.nextToken()
				tp = FloatRangeType(f, math.MaxFloat64) // Unbounded at upper end
			}
		} else {
			tp = Float(f)
		}
	case dotdot: // Unbounded at lower end
		n := p.peekToken()
		switch n.i {
		case integer:
			p.nextToken()
			tp = IntegerRangeType(math.MinInt64, tokenInt(n))
		case float:
			p.nextToken()
			tp = FloatRangeType(-math.MaxFloat64, tokenFloat(n))
		default:
			panic(badSyntax(n, exIntOrFloat))
		}
	case identifier:
		switch t.s {
		case `map`:
			n := p.nextToken()
			if n.i != '[' {
				panic(badSyntax(n, exLeftBracket))
			}

			// Deal with key type and size constraint
			p.anyOf(p.nextToken())
			keyType := p.popLastType()

			n = p.nextToken()
			var szc *array
			if n.i == ',' {
				// get size arguments
				p.params()
				szc = p.popLast().(*array)
			} else if n.i != ']' {
				panic(badSyntax(n, exRightBracket))
			}

			p.typeExpression(p.nextToken())
			params := &array{slice: []dgo.Value{keyType, p.popLastType()}}
			if szc != nil {
				params.AddAll(szc)
			}
			tp = MapType(sliceToInterfaces(params)...)
		case `type`:
			n := p.nextToken()
			if n.i != '[' {
				panic(badSyntax(n, exLeftBracket))
			}

			// Deal with key type and size constraint
			p.anyOf(p.nextToken())
			typ := p.popLastType()
			if p.nextToken().i != ']' {
				panic(badSyntax(n, exRightBracket))
			}
			tp = &metaType{typ}
		case `any`:
			tp = DefaultAnyType
		case `bool`:
			tp = DefaultBooleanType
		case `int`:
			tp = DefaultIntegerType
		case `float`:
			tp = DefaultFloatType
		case `string`:
			if p.peekToken().i == '[' {
				// get size arguments
				p.nextToken()
				p.params()
				szc := p.popLast().(*array)
				tp = StringType(sliceToInterfaces(szc)...)
			} else {
				tp = DefaultStringType
			}
		case `binary`:
			tp = DefaultBinaryType
		case `true`:
			tp = True
		case `false`:
			tp = False
		case `nil`:
			tp = Nil
		default:
			panic(fmt.Errorf(`unknown identifier '%s'`, t.s))
		}
	case stringLiteral:
		tp = String(t.s)
	case regexpLiteral:
		tp = PatternType(regexp.MustCompile(t.s))
	default:
		panic(badSyntax(t, exTypeExpression))
	}
	p.d = append(p.d, tp)
}

func tokenInt(t *token) int64 {
	i, _ := strconv.ParseInt(t.s, 0, 64)
	return i
}

func tokenFloat(t *token) float64 {
	f, _ := strconv.ParseFloat(t.s, 64)
	return f
}

func sliceToInterfaces(a *array) []interface{} {
	s := a.slice
	is := make([]interface{}, len(s))
	for i := range s {
		is[i] = s[i]
	}
	return is
}

func (p *parser) popLast() dgo.Value {
	last := len(p.d) - 1
	v := p.d[last]
	p.d = p.d[:last]
	return v
}

func (p *parser) popLastType() dgo.Type {
	last := len(p.d) - 1
	v := p.d[last]
	p.d = p.d[:last]
	if tp, ok := v.(dgo.Type); ok {
		return tp
	}
	return v.Type()
}

func allTypes(a []dgo.Value) []dgo.Value {
	c := make([]dgo.Value, len(a))
	for i := range a {
		v := a[i]
		if _, ok := v.(dgo.Type); !ok {
			v = v.Type()
		}
		c[i] = v
	}
	return c
}
