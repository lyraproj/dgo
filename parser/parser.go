package parser

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/lyraproj/dgo/vf"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/util"
)

// States:
const (
	exListComma = iota
	exListEnd
	exParamsComma
	exLeftBracket
	exLeftParen
	exRightBracket
	exRightParen
	exRightAngle
	exIntOrFloat
	exStringLiteral
	exTypeExpression
	exAliasRef
	exEnd
)

type unknownIdentifier struct {
	dgo.Value
}

type optionalValue struct {
	dgo.Value
}

type alias struct {
	dgo.StringType
}

type dType = dgo.Type // To avoid collision with method named Type
type deferredCall struct {
	dType
	args dgo.Arguments
}

type LexFunction func(reader *util.StringReader) *Token

// NewAlias creates a special interim type that represents a type alias used during parsing, and then nowhere else.
func NewAlias(s dgo.String) dgo.Type {
	return &alias{s.Type().(dgo.StringType)}
}

// NewCall creates a special interim type that represents a call during parsing, and then nowhere else.
func NewCall(s dgo.Type, args dgo.Arguments) dgo.Type {
	return &deferredCall{s, args}
}

func (a *alias) GoString() string {
	return a.StringType.(dgo.ExactType).Value().(dgo.String).GoString()
}

type aliasProvider struct {
	sc dgo.AliasMap
}

func (p *aliasProvider) Replace(t dgo.Value) dgo.Value {
	switch t := t.(type) {
	case *deferredCall:
		return vf.New(t.dType, t.args)
	case *alias:
		if ra := p.sc.GetType(internal.String(t.GoString())); ra != nil {
			return ra
		}
		panic(fmt.Errorf(`reference to unresolved type '%s'`, t.GoString()))
	case dgo.AliasContainer:
		t.Resolve(p)
	}
	return t
}

func expect(state int) (s string) {
	switch state {
	case exParamsComma:
		s = `one of ',' or ']'`
	case exLeftBracket:
		s = `'['`
	case exLeftParen:
		s = `'('`
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
	case exStringLiteral:
		s = `a literal string`
	case exTypeExpression:
		s = `a type expression`
	case exAliasRef:
		s = `an identifier`
	case exEnd:
		s = `end of expression`
	}
	return
}

func (p *parser) badSyntax(t *Token, state int) error {
	return fmt.Errorf(`expected %s, got %s`, expect(state), p.TokenString(t))
}

type (
	// Parser is the interface that a parser is required to implement
	Parser interface {
		AliasMap() dgo.AliasMap
		NextToken() *Token
		Parse(t *Token)
		PopLast() dgo.Value
		LastToken() *Token
		StringReader() *util.StringReader
		TokenString(*Token) string
	}

	// ParserBase provides all methods of the Parser interface except the
	// Parse method.
	ParserBase struct {
		lf LexFunction
		sc dgo.AliasMap
		d  []dgo.Value
		sr *util.StringReader
		pe *Token
		lt *Token
	}

	parser struct {
		ParserBase
	}
)

// NewParserBase creates the extendable parser base.
func NewParserBase(am dgo.AliasMap, lf LexFunction, content string) ParserBase {
	if am == nil {
		am = internal.NewAliasMap()
	}
	return ParserBase{lf: lf, sc: am, sr: util.NewStringReader(content)}
}

// Parse calls ParseFile with an empty string as the fileName
func Parse(content string) dgo.Value {
	return ParseFile(nil, ``, content)
}

// ParseFile parses the given content into a dgo.Type. The filename is used in error messages.
//
// The alias map is optional. If given, the parser will recognize the type aliases provided in the map
// and also add any new aliases declared within the parsed content to that map.
func ParseFile(am dgo.AliasMap, fileName, content string) dgo.Value {
	p := &parser{NewParserBase(am, nextToken, content)}
	return DoParse(p, fileName)
}

// DoParse performs the actual parsing and returns the result
func DoParse(p Parser, fileName string) dgo.Value {
	defer func() {
		if r := recover(); r != nil {
			es := r
			if err, ok := r.(error); ok {
				es = err.Error()
			}
			tl := 1
			lt := p.LastToken()
			if lt != nil && lt.Value != `` {
				tl = len(lt.Value)
			}
			fn := ``
			if fileName != `` {
				fn = fmt.Sprintf(`file: %s, `, fileName)
			}
			ln := ``
			sr := p.StringReader()
			if fileName != `` || sr.Line() > 1 {
				ln = fmt.Sprintf(`line: %d, `, sr.Line())
			}
			panic(fmt.Errorf("%s: (%s%scolumn: %d)", es, fn, ln, sr.Column()-tl))
		}
	}()
	p.Parse(p.NextToken())
	tp := p.PopLast()
	return (&aliasProvider{p.AliasMap()}).Replace(tp)
}

func (p *ParserBase) AliasMap() dgo.AliasMap {
	return p.sc
}

func (p *ParserBase) Append(v dgo.Value) {
	p.d = append(p.d, v)
}

func (p *ParserBase) AppendFrom(pos int, v dgo.Value) {
	p.d = append(p.d[:pos], v)
}

func (p *ParserBase) From(pos int) []dgo.Value {
	return p.d[pos:]
}

func (p *ParserBase) Len() int {
	return len(p.d)
}

func (p *ParserBase) PeekToken() *Token {
	if p.pe == nil {
		p.pe = p.NextToken()
	}
	return p.pe
}

func (p *ParserBase) NextToken() *Token {
	var t *Token
	if p.pe != nil {
		t = p.pe
		p.pe = nil
	} else {
		t = p.lf(p.sr)
	}
	p.lt = t
	return t
}

func (p *ParserBase) LastToken() *Token {
	return p.lt
}

func (p *ParserBase) PopLast() dgo.Value {
	last := p.Len() - 1
	v := p.d[last]
	p.d = p.d[:last]
	return v
}

func (p *ParserBase) PopLastType() dgo.Type {
	last := p.Len() - 1
	v := p.d[last]
	p.d = p.d[:last]
	if tp, ok := v.(dgo.Type); ok {
		return tp
	}
	return v.Type()
}

func (p *ParserBase) StringReader() *util.StringReader {
	return p.sr
}

func (p *parser) Parse(t *Token) {
	p.anyOf(t)
	tk := p.NextToken()
	if tk.Type != end {
		panic(p.badSyntax(tk, exEnd))
	}
}

func (p *parser) TokenString(t *Token) string {
	return tokenString(t)
}

func (p *parser) list(endChar int) {
	szp := p.Len()
	ellipsis := false
	expectEntry := 0
	if endChar == '}' {
		expectEntry = 1
	}
	for {
		t := p.NextToken()
		if t.Type == dotdotdot {
			ellipsis = true
			t = p.NextToken()
			if t.Type == endChar && expectEntry != 0 {
				break
			}
			if expectEntry == 2 {
				panic(p.badSyntax(t, exListEnd))
			}
			expectEntry = 0
			p.anyOf(t)
			t = p.NextToken()
			et := internal.ArrayType([]interface{}{p.PopLastType(), 0, math.MaxInt64})
			p.Append(et)
		}

		if t.Type == endChar {
			// Right bracket instead of element indicates an empty array or an extraneous comma. Both are OK
			break
		}
		expectEntry = p.arrayElement(t, expectEntry)
		t = p.NextToken()
		if t.Type == endChar {
			break
		}
		if t.Type != ',' {
			panic(p.badSyntax(t, exListComma))
		}
	}

	as := p.From(szp)
	var tv dgo.Value
	if len(as) > 0 {
		if expectEntry == 2 {
			tv = makeStructType(as, ellipsis)
		}
		if tv == nil {
			tv = makeTupleType(as, ellipsis)
		}
	} else {
		if expectEntry == 0 {
			tv = internal.DefaultTupleType
		} else {
			tv = makeStructType(nil, ellipsis)
		}
	}
	p.AppendFrom(szp, tv)
}

func makeTupleType(as []dgo.Value, variadic bool) dgo.TupleType {
	// Convert literal values to types and create a tupleType
	ln := len(as)
	ts := make([]interface{}, ln)
	for i := range as {
		v := as[i]
		t, ok := v.(dgo.Type)
		if !ok {
			t = v.Type()
		}
		ts[i] = t
	}
	if variadic {
		ln--
		ts[ln] = ts[ln].(dgo.ArrayType).ElementType()
		return internal.VariadicTupleType(ts)
	}
	return internal.TupleType(ts)
}

func makeStructType(as []dgo.Value, ellipsis bool) dgo.StructMapType {
	l := len(as)
	entries := make([]dgo.StructMapEntry, l)
	for i := range as {
		hn := as[i].(dgo.MapEntry)
		k := hn.Key()
		v := hn.Value()
		kt, isType := k.(dgo.Type)
		if !isType {
			kt = k.Type()
		}
		ov, optional := v.(*optionalValue)
		if optional {
			v = ov.Value
		}
		vt, isType := v.(dgo.Type)
		if !isType {
			vt = v.Type()
		}
		entries[i] = internal.StructMapEntry(kt, vt, !optional)
	}
	return internal.StructMapType(ellipsis, entries)
}

func (p *parser) params() {
	szp := p.Len()
	for {
		t := p.NextToken()
		if t.Type == ']' {
			// Right bracket instead of element indicates an empty array or an extraneous comma. Both are OK
			break
		}
		p.arrayElement(t, 0)
		t = p.NextToken()
		if t.Type == ']' {
			break
		}
		if t.Type != ',' {
			panic(p.badSyntax(t, exParamsComma))
		}
	}
	as := p.From(szp)
	tv := internal.WrapSlice(as).Copy(false)
	p.AppendFrom(szp, tv)
}

func (p *parser) arrayElement(t *Token, expectEntry int) int {
	var key dgo.Value
	nt := p.PeekToken()
	if t.Type == identifier && nt.Type == ':' || nt.Type == '?' {
		if expectEntry == 0 {
			panic(errors.New(`mix of elements and map entries`))
		}
		key = internal.String(t.Value)
	} else {
		p.anyOf(t)
		nt = p.PeekToken()
	}

	optional := nt.Type == '?'
	if optional {
		p.NextToken()
	}

	if p.PeekToken().Type == ':' {
		if expectEntry == 0 {
			panic(errors.New(`mix of elements and map entries`))
		}
		// Map mapEntry
		p.NextToken()
		if key == nil {
			key = p.PopLast()
		}
		p.anyOf(p.NextToken())
		val := p.PopLast()
		if optional {
			val = &optionalValue{val}
		}
		p.Append(internal.NewMapEntry(key, val))
		expectEntry = 2
	} else {
		if expectEntry == 2 {
			panic(errors.New(`mix of elements and map entries`))
		}
		expectEntry = 0
	}
	return expectEntry
}

func (p *parser) anyOf(t *Token) {
	p.oneOf(t)
	if p.PeekToken().Type == '|' {
		szp := p.Len() - 1
		for {
			p.NextToken()
			p.oneOf(p.NextToken())
			if p.PeekToken().Type != '|' {
				p.AppendFrom(szp, internal.AnyOfType(allTypes(p.From(szp))))
				break
			}
		}
	}
}

func (p *parser) oneOf(t *Token) {
	p.allOf(t)
	if p.PeekToken().Type == '^' {
		szp := p.Len() - 1
		for {
			p.NextToken()
			p.allOf(p.NextToken())
			if p.PeekToken().Type != '^' {
				p.AppendFrom(szp, internal.OneOfType(allTypes(p.From(szp))))
				break
			}
		}
	}
}

func (p *parser) allOf(t *Token) {
	p.unary(t)
	if p.PeekToken().Type == '&' {
		szp := p.Len() - 1
		for {
			p.NextToken()
			p.unary(p.NextToken())
			if p.PeekToken().Type != '&' {
				p.AppendFrom(szp, internal.AllOfType(allTypes(p.From(szp))))
				break
			}
		}
	}
}

func (p *parser) unary(t *Token) {
	negate := false
	ciString := false
	if t.Type == '!' {
		negate = true
		t = p.NextToken()
	}
	if t.Type == '~' {
		ciString = true
		t = p.NextToken()
	}

	p.typeExpression(t)

	if ciString {
		if s, ok := p.PopLast().(dgo.String); ok {
			p.Append(internal.CiStringType(s))
		} else {
			panic(p.badSyntax(t, exStringLiteral))
		}
	}
	if negate {
		nt := internal.NotType(p.PopLastType())
		p.Append(nt)
	}
}

func (p *parser) mapExpression() dgo.Value {
	n := p.NextToken()
	if n.Type != '[' {
		panic(p.badSyntax(n, exLeftBracket))
	}

	// Deal with key type and size constraint
	p.anyOf(p.NextToken())
	keyType := p.PopLastType()

	n = p.NextToken()
	var szc dgo.Array
	if n.Type == ',' {
		// get size arguments
		p.params()
		szc = p.PopLast().(dgo.Array)
	} else if n.Type != ']' {
		panic(p.badSyntax(n, exRightBracket))
	}

	p.typeExpression(p.NextToken())
	params := internal.WrapSlice([]dgo.Value{keyType, p.PopLastType()})
	if szc != nil {
		params.AddAll(szc)
	}
	return internal.MapType(params.InterfaceSlice())
}

func (p *parser) meta() dgo.Value {
	if p.PeekToken().Type != '[' {
		return internal.DefaultMetaType
	}
	p.NextToken()

	// Deal with key type and size constraint
	p.anyOf(p.NextToken())
	tp := p.PopLastType()
	t := p.NextToken()
	if t.Type != ']' {
		panic(p.badSyntax(t, exRightBracket))
	}
	return internal.MetaType(tp)
}

func (p *parser) string() dgo.Value {
	if p.PeekToken().Type == '[' {
		// get size arguments
		p.NextToken()
		p.params()
		szc := p.PopLast().(dgo.Array)
		return internal.StringType(szc.InterfaceSlice())
	}
	return internal.DefaultStringType
}

func (p *parser) sensitive() dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.params()
		szc := p.PopLast().(dgo.Array)
		return internal.SensitiveType(szc.InterfaceSlice())
	}
	if !isExpressionEnd(rune(tt)) {
		p.anyOf(p.NextToken())
		return internal.Sensitive(internal.ExactValue(p.PopLast()))
	}
	return internal.DefaultSensitiveType
}

func (p *parser) funcExpression() dgo.Value {
	t := p.NextToken()
	if t.Type != '(' {
		panic(p.badSyntax(t, exLeftParen))
	}
	p.list(')')
	args := p.PopLastType().(dgo.TupleType)
	var returns dgo.TupleType = internal.EmptyTupleType
	t = p.PeekToken()
	switch {
	case t.Type == '(':
		p.NextToken()
		p.list(')')
		returns = p.PopLastType().(dgo.TupleType)
	case isExpressionEnd(rune(t.Type)):
		break
	default:
		p.anyOf(p.NextToken())
		returns = internal.TupleType([]interface{}{p.PopLastType()})
	}
	return internal.FunctionType(args, returns)
}

var identifierToTypeMap = map[string]dgo.Value{
	`any`:    internal.DefaultAnyType,
	`bool`:   internal.DefaultBooleanType,
	`int`:    internal.DefaultIntegerType,
	`float`:  internal.DefaultFloatType,
	`dgo`:    internal.DefaultDgoStringType,
	`binary`: internal.DefaultBinaryType,
	`true`:   internal.True,
	`false`:  internal.False,
	`nil`:    internal.Nil,
}

func (p *parser) identifier(t *Token, returnUnknown bool) dgo.Value {
	tp, ok := identifierToTypeMap[t.Value]
	if ok {
		return tp
	}
	switch t.Value {
	case `map`:
		tp = p.mapExpression()
	case `type`:
		tp = p.meta()
	case `string`:
		tp = p.string()
	case `sensitive`:
		tp = p.sensitive()
	case `func`:
		tp = p.funcExpression()
	default:
		if returnUnknown {
			tp = &unknownIdentifier{internal.String(t.Value)}
		} else {
			tp = p.namedType(t)
		}
	}
	return tp
}

func (p *parser) namedType(t *Token) dgo.Value {
	tp := p.aliasReference(t)
	if nt, ok := tp.(dgo.NamedType); ok {
		if !isExpressionEnd(rune(p.PeekToken().Type)) {
			p.anyOf(p.NextToken())
			tp = nt.New(internal.ExactValue(p.PopLast()))
		}
	}
	return tp
}

func (p *parser) float(t *Token) dgo.Value {
	var tp dgo.Value
	f := tokenFloat(t)
	n := p.PeekToken()
	if n.Type == dotdot || n.Type == dotdotdot {
		inclusive := n.Type == dotdot
		p.NextToken()
		n = p.PeekToken()
		if n.Type == integer || n.Type == float {
			p.NextToken()
			tp = internal.FloatRangeType(f, tokenFloat(n), inclusive)
		} else {
			p.NextToken()
			tp = internal.FloatRangeType(f, math.MaxFloat64, inclusive) // Unbounded at upper end
		}
	} else {
		tp = internal.Float(f)
	}
	return tp
}

func (p *parser) integer(t *Token) dgo.Value {
	var tp dgo.Value
	i := tokenInt(t)
	n := p.PeekToken()
	if n.Type == dotdot || n.Type == dotdotdot {
		inclusive := n.Type == dotdot
		p.NextToken()
		x := p.PeekToken()
		switch x.Type {
		case integer:
			p.NextToken()
			tp = internal.IntegerRangeType(i, tokenInt(x), inclusive)
		case float:
			p.NextToken()
			tp = internal.FloatRangeType(float64(i), tokenFloat(x), inclusive)
		default:
			tp = internal.IntegerRangeType(i, math.MaxInt64, inclusive) // Unbounded at upper end
		}
	} else {
		tp = internal.Integer(i)
	}
	return tp
}

func (p *parser) dotRange(t *Token) dgo.Value {
	var tp dgo.Value
	n := p.PeekToken()
	inclusive := t.Type == dotdot
	switch n.Type {
	case integer:
		p.NextToken()
		tp = internal.IntegerRangeType(math.MinInt64, tokenInt(n), inclusive)
	case float:
		p.NextToken()
		tp = internal.FloatRangeType(-math.MaxFloat64, tokenFloat(n), inclusive)
	default:
		panic(p.badSyntax(n, exIntOrFloat))
	}
	return tp
}

func (p *parser) array(t *Token) dgo.Value {
	p.params()
	params := p.PopLast().(dgo.Array)
	p.typeExpression(p.NextToken())
	params.Insert(0, p.PopLastType())
	return internal.ArrayType(params.InterfaceSlice())
}

func (p *parser) aliasReference(t *Token) dgo.Value {
	if t.Type != identifier {
		panic(p.badSyntax(t, exAliasRef))
	}
	if tp := internal.NamedType(t.Value); tp != nil {
		return tp
	}
	vn := internal.String(t.Value)
	if tp := p.sc.GetType(vn); tp != nil {
		return tp
	}
	return NewAlias(vn)
}

func (p *parser) aliasDeclaration(t *Token) dgo.Value {
	// Should result in an unknown identifier or name is reserved
	tp := p.identifier(t, true)
	if un, ok := tp.(*unknownIdentifier); ok {
		s := un.Value.(dgo.String)
		if internal.NamedType(s.GoString()) == nil && p.sc.GetType(s) == nil {
			p.NextToken() // skip '='
			p.sc.Add(NewAlias(s), s)
			p.anyOf(p.NextToken())
			tp = p.PopLastType()
			p.sc.Add(tp.(dgo.Type), s)
			return tp
		}
	}
	panic(fmt.Errorf(`attempt to redeclare identifier '%s'`, t.Value))
}

func (p *parser) typeExpression(t *Token) {
	var tp dgo.Value
	switch t.Type {
	case '{':
		p.list('}')
		return
	case '(':
		p.anyOf(p.NextToken())
		n := p.NextToken()
		if n.Type != ')' {
			panic(p.badSyntax(n, exRightParen))
		}
		return
	case '[':
		tp = p.array(t)
	case '<':
		tp = p.aliasReference(p.NextToken())
		n := p.NextToken()
		if n.Type != '>' {
			panic(p.badSyntax(n, exRightAngle))
		}
	case integer:
		tp = p.integer(t)
	case float:
		tp = p.float(t)
	case dotdot, dotdotdot: // Unbounded at lower end
		tp = p.dotRange(t)
	case identifier:
		if p.PeekToken().Type == '=' {
			tp = p.aliasDeclaration(t)
		} else {
			tp = p.identifier(t, false)
		}
	case stringLiteral:
		tp = internal.String(t.Value)
	case regexpLiteral:
		tp = internal.PatternType(regexp.MustCompile(t.Value))
	default:
		panic(p.badSyntax(t, exTypeExpression))
	}
	p.Append(tp)
}

func tokenInt(t *Token) int64 {
	i, _ := strconv.ParseInt(t.Value, 0, 64)
	return i
}

func tokenFloat(t *Token) float64 {
	f, _ := strconv.ParseFloat(t.Value, 64)
	return f
}

func allTypes(a []dgo.Value) []interface{} {
	c := make([]interface{}, len(a))
	for i := range a {
		v := a[i]
		if tv, ok := v.(dgo.Type); ok {
			c[i] = tv
		} else {
			c[i] = v.Type()
		}
	}
	return c
}

func init() {
	internal.Parse = Parse
}
