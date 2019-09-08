package internal

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

// States:
const (
	exListComma = iota
	exListEnd
	exParamsComma
	exLeftBracket
	exRightBracket
	exRightParen
	exRightAngle
	exInt
	exIntOrFloat
	exFloatDotDotDot
	exTypeExpression
	exAliasRef
	exEnd
)

type unknownIdentifier struct {
	hstring
}

type optionalValue struct {
	value dgo.Value
}

type alias struct {
	exactStringType
}

type aliasProvider struct {
	sc map[string]dgo.Type
}

func (p *aliasProvider) Replace(t dgo.Type) dgo.Type {
	if a, ok := t.(*alias); ok {
		if ra, ok := p.sc[a.s]; ok {
			return ra
		}
		panic(fmt.Errorf(`reference to unresolved alias '%s'`, a.s))
	}
	if rt, ok := t.(dgo.AliasContainer); ok {
		rt.Resolve(p)
	}
	return t
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
	case exListEnd:
		s = `'}'`
	case exRightParen:
		s = `')'`
	case exRightAngle:
		s = `'>'`
	case exInt:
		s = `an integer`
	case exIntOrFloat:
		s = `an integer or a float`
	case exFloatDotDotDot:
		s = `float boundary on ... range`
	case exTypeExpression:
		s = `a type expression`
	case exAliasRef:
		s = `an identifier`
	case exEnd:
		s = `end of expression`
	}
	return
}

func badSyntax(t *token, state int) error {
	if state == exFloatDotDotDot {
		return errors.New(expect(state))
	}
	return fmt.Errorf(`expected %s, got %s`, expect(state), t)
}

type parser struct {
	d  []dgo.Value
	sr *util.StringReader
	sc map[string]dgo.Type
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
	tp := p.popLastType()
	if p.sc != nil {
		tp = (&aliasProvider{p.sc}).Replace(tp)
	}
	return tp
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
	ellipsis := false
	for {
		t := p.nextToken()
		if t.i == dotdotdot {
			ellipsis = true
			t = p.nextToken()
			if t.i == '}' {
				break
			}
			panic(badSyntax(t, exListEnd))
		}

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
			tv = makeStructType(as, ellipsis)
		}
		if tv == nil {
			tv = makeTupleType(ar)
		}
	} else {
		tv = makeStructType(nil, ellipsis)
	}
	p.d = append(p.d[:szp], tv)
}

func makeTupleType(ar dgo.Array) dgo.TupleType {
	// Convert literal values to types and create a tupleType
	as := sliceCopy(ar.GoSlice())
	for i := range as {
		v := as[i]
		if _, ok := v.(dgo.Type); !ok {
			as[i] = v.Type()
		}
	}
	return &tupleType{slice: as, frozen: true}
}

func makeStructType(as []dgo.Value, ellipsis bool) dgo.StructType {
	l := len(as)
	keys := make([]dgo.Value, l)
	values := make([]dgo.Value, l)
	required := make([]byte, l)
	for i := range as {
		hn, ok := as[i].(dgo.MapEntry)
		if !ok {
			return nil
		}
		k := hn.Key()
		v := hn.Value()
		var kt, vt dgo.Value
		var isType bool
		if kt, isType = k.(dgo.Type); !isType {
			kt = k.Type()
		}
		keys[i] = kt
		if ov, optional := v.(*optionalValue); optional {
			v = ov.value
		} else {
			required[i] = 1
		}
		if vt, isType = v.(dgo.Type); !isType {
			vt = v.Type()
		}
		values[i] = vt
	}
	return &structType{additional: ellipsis, keys: array{slice: keys, frozen: true}, values: array{slice: values, frozen: true}, required: required}
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
	var id dgo.Value
	if t.i == identifier && p.peekToken().i != '=' {
		// Special handling of identifiers at this point
		id = p.identifier(t, true)
		p.d = append(p.d, id)
	} else {
		p.anyOf(t)
	}
	optional := p.peekToken().i == '?'
	if optional {
		p.nextToken()
	}
	if p.peekToken().i == ':' {
		// Map mapEntry
		p.nextToken()
		key := p.popLast()
		if unknown, ok := key.(*unknownIdentifier); ok {
			key = String(unknown.s)
		}
		p.anyOf(p.nextToken())
		val := p.popLast()
		if optional {
			val = &optionalValue{val}
		}
		p.d = append(p.d, &mapEntry{key: key, value: val})
	} else if unknown, ok := id.(*unknownIdentifier); ok {
		panic(fmt.Errorf(`unknown identifier '%s'`, unknown.s))
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

func (p *parser) mapExpression() dgo.Value {
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
	return MapType(sliceToInterfaces(params)...)
}

func (p *parser) meta() dgo.Value {
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
	return &metaType{typ}
}

func (p *parser) string() dgo.Value {
	if p.peekToken().i == '[' {
		// get size arguments
		p.nextToken()
		p.params()
		szc := p.popLast().(*array)
		return StringType(sliceToInterfaces(szc)...)
	}
	return DefaultStringType
}

func (p *parser) identifier(t *token, acceptUnknown bool) dgo.Value {
	var tp dgo.Value
	switch t.s {
	case `map`:
		tp = p.mapExpression()
	case `type`:
		tp = p.meta()
	case `any`:
		tp = DefaultAnyType
	case `bool`:
		tp = DefaultBooleanType
	case `int`:
		tp = DefaultIntegerType
	case `float`:
		tp = DefaultFloatType
	case `dgo`:
		tp = DefaultDgoStringType
	case `string`:
		tp = p.string()
	case `binary`:
		tp = DefaultBinaryType
	case `true`:
		tp = True
	case `false`:
		tp = False
	case `nil`:
		tp = Nil
	default:
		if !acceptUnknown {
			panic(fmt.Errorf(`unknown identifier '%s'`, t.s))
		}
		tp = &unknownIdentifier{hstring: hstring{s: t.s}}
	}
	return tp
}

func (p *parser) float(t *token) dgo.Value {
	var tp dgo.Value
	f := tokenFloat(t)
	n := p.peekToken()
	if n.i == dotdotdot {
		panic(badSyntax(t, exFloatDotDotDot))
	}
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
		tp = floatVal(f)
	}
	return tp
}

func (p *parser) integer(t *token) dgo.Value {
	var tp dgo.Value
	i := tokenInt(t)
	n := p.peekToken()
	if n.i == dotdot || n.i == dotdotdot {
		p.nextToken()
		x := p.peekToken()
		switch x.i {
		case integer:
			p.nextToken()
			u := tokenInt(x)
			if n.i == dotdotdot {
				u--
			}
			tp = IntegerRangeType(i, u)
		case float:
			if n.i == dotdotdot {
				panic(badSyntax(x, exFloatDotDotDot))
			}
			p.nextToken()
			tp = FloatRangeType(float64(i), tokenFloat(x))
		default:
			// .. or ... doesn't matter since there's no upper bound to apply the difference on.
			tp = IntegerRangeType(i, math.MaxInt64) // Unbounded at upper end
		}
	} else {
		tp = intVal(i)
	}
	return tp
}

func (p *parser) dotRange(t *token) dgo.Value {
	var tp dgo.Value
	n := p.peekToken()
	switch n.i {
	case integer:
		p.nextToken()
		u := tokenInt(n)
		if t.i == dotdotdot {
			u--
		}
		tp = IntegerRangeType(math.MinInt64, u)
	case float:
		if t.i == dotdotdot {
			panic(badSyntax(n, exFloatDotDotDot))
		}
		p.nextToken()
		tp = FloatRangeType(-math.MaxFloat64, tokenFloat(n))
	default:
		ex := exIntOrFloat
		if t.i == dotdotdot {
			ex = exInt
		}
		panic(badSyntax(n, ex))
	}
	return tp
}

func (p *parser) array(t *token) dgo.Value {
	p.params()
	params := p.popLast().(*array)
	p.typeExpression(p.nextToken())
	params.Insert(0, p.popLastType())
	return ArrayType(sliceToInterfaces(params)...)
}

func (p *parser) aliasReference(t *token) dgo.Value {
	if t.i != identifier {
		panic(badSyntax(t, exAliasRef))
	}
	if p.sc != nil {
		if tp, ok := p.sc[t.s]; ok {
			t = p.nextToken()
			if t.i != '>' {
				panic(badSyntax(t, exRightAngle))
			}
			return tp
		}
	}
	panic(fmt.Errorf(`unknown alias '%s'`, t.s))
}

func (p *parser) aliasDeclaration(t *token) dgo.Value {
	// Should result in an unknown identifier or name is reserved
	tp := p.identifier(t, true)
	if un, ok := tp.(*unknownIdentifier); ok {
		p.nextToken() // skip '='
		if p.sc == nil {
			p.sc = make(map[string]dgo.Type)
		}
		p.sc[un.s] = &alias{exactStringType: exactStringType{s: un.s}}
		p.parse(p.nextToken())
		tp = p.popLastType()
		p.sc[un.s] = tp.(dgo.Type)
		return tp
	}
	panic(fmt.Errorf(`attempt redeclare identifier '%s'`, t.s))
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
		tp = p.array(t)
	case '<':
		tp = p.aliasReference(p.nextToken())
	case integer:
		tp = p.integer(t)
	case float:
		tp = p.float(t)
	case dotdot, dotdotdot: // Unbounded at lower end
		tp = p.dotRange(t)
	case identifier:
		if p.peekToken().i == '=' {
			tp = p.aliasDeclaration(t)
		} else {
			tp = p.identifier(t, false)
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
