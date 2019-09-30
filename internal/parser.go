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
	exIntOrFloat
	exTypeExpression
	exAliasRef
	exEnd
)

type unknownIdentifier struct {
	hstring
}

type optionalValue struct {
	dgo.Value
}

type alias struct {
	exactStringType
}

type aliasProvider struct {
	sc dgo.AliasMap
}

func (p *aliasProvider) Replace(t dgo.Type) dgo.Type {
	if a, ok := t.(*alias); ok {
		if ra := p.sc.GetType(String(a.s)); ra != nil {
			return ra
		}
		panic(fmt.Errorf(`reference to unresolved type '%s'`, a.s))
	}
	if rt, ok := t.(dgo.AliasContainer); ok {
		rt.Resolve(p)
	}
	return t
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
	case exIntOrFloat:
		s = `an integer or a float`
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
	return fmt.Errorf(`expected %s, got %s`, expect(state), t)
}

type parser struct {
	sc dgo.AliasMap
	d  []dgo.Value
	sr *util.StringReader
	pe *token
	lt *token
}

// Parse calls ParseFile with the string "<dgo type expression>" as the fileName
func Parse(content string) dgo.Type {
	return ParseFile(nil, ``, content)
}

// ParseFile parses the given content into a dgo.Type. The filename is used in error messages.
//
// The alias map is optional. If given, the parser will recognize the type aliases provided in the map
// and also add any new aliases declared within the parsed content to that map.
func ParseFile(am dgo.AliasMap, fileName, content string) dgo.Type {
	if am == nil {
		am = &aliasMap{}
	}
	p := &parser{sc: am, sr: util.NewStringReader(content)}

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
			panic(fmt.Errorf("%s: (%s%scolumn: %d)", es, fn, ln, p.sr.Column()-tl))
		}
	}()
	p.parse(p.nextToken())
	tp := p.popLastType()
	return (&aliasProvider{p.sc}).Replace(tp)
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
	expectEntry := 1
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
		expectEntry = p.arrayElement(t, expectEntry)
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
		if expectEntry == 2 {
			tv = makeStructType(as, ellipsis)
		}
		if tv == nil {
			tv = makeTupleType(&array{slice: as})
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
		hn := as[i].(dgo.MapEntry)
		k := hn.Key()
		v := hn.Value()
		var kt, vt dgo.Value
		var isType bool
		if kt, isType = k.(dgo.Type); !isType {
			kt = k.Type()
		}
		keys[i] = kt
		if ov, optional := v.(*optionalValue); optional {
			v = ov.Value
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
		p.arrayElement(t, 0)
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

func (p *parser) arrayElement(t *token, expectEntry int) int {
	var key dgo.Value
	nt := p.peekToken()
	if t.i == identifier && nt.i == ':' || nt.i == '?' {
		if expectEntry == 0 {
			panic(errors.New(`mix of elements and map entries`))
		}
		key = String(t.s)
	} else {
		p.anyOf(t)
		nt = p.peekToken()
	}

	optional := nt.i == '?'
	if optional {
		p.nextToken()
	}

	if p.peekToken().i == ':' {
		if expectEntry == 0 {
			panic(errors.New(`mix of elements and map entries`))
		}
		// Map mapEntry
		p.nextToken()
		if key == nil {
			key = p.popLast()
		}
		p.anyOf(p.nextToken())
		val := p.popLast()
		if optional {
			val = &optionalValue{val}
		}
		p.d = append(p.d, &mapEntry{key: key, value: val})
		expectEntry = 2
	} else {
		if expectEntry == 2 {
			panic(errors.New(`mix of elements and map entries`))
		}
		expectEntry = 0
	}
	return expectEntry
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
	if p.peekToken().i != '[' {
		return &metaType{DefaultAnyType}
	}
	p.nextToken()

	// Deal with key type and size constraint
	p.anyOf(p.nextToken())
	typ := p.popLastType()
	t := p.nextToken()
	if t.i != ']' {
		panic(badSyntax(t, exRightBracket))
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

func (p *parser) identifier(t *token, returnUnknown bool) dgo.Value {
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
		if returnUnknown {
			tp = &unknownIdentifier{hstring: hstring{s: t.s}}
		} else {
			tp = p.aliasReference(t)
		}
	}
	return tp
}

func (p *parser) float(t *token) dgo.Value {
	var tp dgo.Value
	f := tokenFloat(t)
	n := p.peekToken()
	if n.i == dotdot || n.i == dotdotdot {
		inclusive := n.i == dotdot
		p.nextToken()
		n = p.peekToken()
		if n.i == integer || n.i == float {
			p.nextToken()
			tp = FloatRangeType(f, tokenFloat(n), inclusive)
		} else {
			p.nextToken()
			tp = FloatRangeType(f, math.MaxFloat64, inclusive) // Unbounded at upper end
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
		inclusive := n.i == dotdot
		p.nextToken()
		x := p.peekToken()
		switch x.i {
		case integer:
			p.nextToken()
			tp = IntegerRangeType(i, tokenInt(x), inclusive)
		case float:
			p.nextToken()
			tp = FloatRangeType(float64(i), tokenFloat(x), inclusive)
		default:
			tp = IntegerRangeType(i, math.MaxInt64, inclusive) // Unbounded at upper end
		}
	} else {
		tp = intVal(i)
	}
	return tp
}

func (p *parser) dotRange(t *token) dgo.Value {
	var tp dgo.Value
	n := p.peekToken()
	inclusive := t.i == dotdot
	switch n.i {
	case integer:
		p.nextToken()
		tp = IntegerRangeType(math.MinInt64, tokenInt(n), inclusive)
	case float:
		p.nextToken()
		tp = FloatRangeType(-math.MaxFloat64, tokenFloat(n), inclusive)
	default:
		panic(badSyntax(n, exIntOrFloat))
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
	if tp := p.sc.GetType(String(t.s)); tp != nil {
		return tp
	}
	return &alias{exactStringType{t.s, 0}}
}

func (p *parser) aliasDeclaration(t *token) dgo.Value {
	// Should result in an unknown identifier or name is reserved
	tp := p.identifier(t, true)
	if un, ok := tp.(*unknownIdentifier); ok {
		n := un.s
		if p.sc.GetType(String(n)) == nil {
			s := String(n)
			p.nextToken() // skip '='
			p.sc.Add(&alias{exactStringType: exactStringType{s: n}}, s)
			p.anyOf(p.nextToken())
			tp = p.popLastType()
			p.sc.Add(tp.(dgo.Type), s)
			return tp
		}
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
		n := p.nextToken()
		if n.i != '>' {
			panic(badSyntax(n, exRightAngle))
		}
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
