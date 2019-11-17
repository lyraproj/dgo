package pcore

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/parser"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

// States:
const (
	exElement     = 0 // Expect value literal
	exParam       = 1 // Expect value literal
	exKey         = 2 // Expect exKey or end of hash
	exValue       = 3 // Expect value
	exEntryValue  = 4 // Expect value
	exRocket      = 5 // Expect rocket
	exListComma   = 6 // Expect comma or end of array
	exParamsComma = 7 // Expect comma or end of parameter list
	exHashComma   = 8 // Expect comma or end of hash
	exName        = 9
	exEqual       = 10
	exListFirst   = 11
	exHashFirst   = 12
	exParamsFirst = 13
	exEnd         = 14
)

func expect(state int) (s string) {
	switch state {
	case exElement, exParam:
		s = `a literal`
	case exKey:
		s = `a hash key`
	case exValue, exEntryValue:
		s = `a hash value`
	case exRocket:
		s = `'=>'`
	case exListComma:
		s = `one of ',' or ']'`
	case exParamsComma:
		s = `one of ',' or ')'`
	case exHashComma:
		s = `one of ',' or '}'`
	case exName:
		s = `a type name`
	case exEqual:
		s = `'='`
	case exListFirst:
		s = `']' or a literal`
	case exParamsFirst:
		s = `')' or a literal`
	case exHashFirst:
		s = `'}' or a hash key`
	case exEnd:
		s = `end of expression`
	}
	return
}

func badSyntax(t *parser.Token, state int) error {
	var ts string
	if t.Type == 0 {
		ts = `EOT`
	} else {
		ts = t.Value
		if ts == `` {
			ts = fmt.Sprintf(`'%c'`, rune(t.Type))
		}
	}
	return fmt.Errorf(`expected %s, got %s`, expect(state), ts)
}

// Parse calls ParseFile with an empty string as the fileName
func Parse(content string) dgo.Value {
	return ParseFile(nil, ``, content)
}

// ParseType parses the given content into a dgo.Type.
func ParseType(content string) dgo.Type {
	return typ.AsType(Parse(content))
}

// ParseFile parses the given content into a dgo.Type. The filename is used in error messages.
//
// The alias map is optional. If given, the pcoreParser will recognize the type aliases provided in the map
// and also add any new aliases declared within the parsed content to that map.
func ParseFile(am dgo.AliasMap, fileName, content string) dgo.Value {
	p := &pcoreParser{parser.NewParserBase(am, nextToken, content)}
	return parser.DoParse(p, fileName)
}

type pcoreParser struct {
	parser.ParserBase
}

func (p *pcoreParser) Parse(t *parser.Token) {
	p.element(t)
	tk := p.NextToken()
	if tk.Type == rocket {
		// Accept top level x => y expression as a singleton hash
		key := p.PopLast()
		p.element(p.NextToken())
		tk = p.NextToken()
		if tk.Type == end {
			p.Append(vf.Map(key, p.PopLast()))
		}
	}
	if tk.Type != end {
		panic(badSyntax(tk, exEnd))
	}
}

func (p *pcoreParser) TokenString(t *parser.Token) string {
	return tokenString(t)
}

func (p *pcoreParser) array() {
	p.list(']', exListComma)
}

func (p *pcoreParser) params() {
	p.list(')', exParamsComma)
}

func (p *pcoreParser) list(et int, bs int) {
	szp := p.Len()
	arrayHash := false

	var rockLhs dgo.Value
	for {
		tk := p.NextToken()
		if rockLhs == nil && tk.Type == et {
			// Right bracket instead of element indicates an empty array or an extraneous comma. Both are OK
			break
		}
		p.element(tk)
		if rockLhs != nil {
			// Last two elements is a hash entry
			p.Append(vf.MapEntry(rockLhs, p.PopLast()))
			rockLhs = nil
			arrayHash = true
		}

		// Comma, rocket, or right bracket must follow element
		tk = p.NextToken()
		switch tk.Type {
		case et:
		case ',':
			continue
		case rocket:
			rockLhs = p.PopLast()
			continue
		default:
			panic(badSyntax(tk, bs))
		}
		break
	}

	a := vf.WrapSlice(p.From(szp)).Copy(false)
	if arrayHash {
		// there's at least one hash entry in the array
		a = convertMapEntries(a)
	}
	p.AppendFrom(szp, a)
}

func (p *pcoreParser) hash() {
	szp := p.Len()
	for {
		tk := p.NextToken()
		// Right curly brace instead of element indicates an empty hash or an extraneous comma. Both are OK
		if tk.Type == '}' {
			break
		}
		p.element(tk)
		tk = p.NextToken()

		// rocket must follow key
		if tk.Type != rocket {
			panic(badSyntax(tk, exRocket))
		}

		p.element(p.NextToken())

		// Comma or right curly brace must follow value
		tk = p.NextToken()
		switch tk.Type {
		case '}':
		case ',':
			continue
		default:
			panic(badSyntax(tk, exHashComma))
		}
		break
	}
	p.AppendFrom(szp, vf.WrapSlice(p.From(szp)).ToMap())
}

func (p *pcoreParser) aliasDeclaration(t *parser.Token) dgo.Value {
	if p.knownType(t) == nil {
		n := toDgoName(t.Value)
		if tf.Named(n) == nil {
			s := vf.String(n)
			am := p.AliasMap()
			if am.GetType(s) == nil {
				am.Add(parser.NewAlias(s), s)
				p.element(p.NextToken())
				tp := p.PopLast()
				am.Add(tp.(dgo.Type), s)
				return tp
			}
		}
	}
	panic(fmt.Errorf(`attempt to redeclare identifier '%s'`, t.Value))
}

func (p *pcoreParser) element(t *parser.Token) {
	switch t.Type {
	case '{':
		p.hash()
	case '[':
		p.array()
	case '(':
		p.params()
	case integer:
		i, err := strconv.ParseInt(t.Value, 0, 64)
		if err != nil {
			panic(err)
		}
		p.Append(vf.Integer(i))
	case float:
		f, err := strconv.ParseFloat(t.Value, 64)
		if err != nil {
			panic(err)
		}
		p.Append(vf.Float(f))
	case identifier:
		switch t.Value {
		case `true`:
			p.Append(vf.True)
		case `type`:
			if p.PeekToken().Type == name {
				t = p.NextToken()
				et := p.NextToken()
				if et.Type != '=' {
					panic(badSyntax(t, exEqual))
				}
				p.Append(p.aliasDeclaration(t))
			} else {
				p.Append(vf.String(t.Value))
			}
		case `false`:
			p.Append(vf.False)
		case `undef`:
			p.Append(vf.Nil)
		default:
			p.Append(vf.String(t.Value))
		}
	case stringLiteral:
		p.Append(vf.String(t.Value))
	case regexpLiteral:
		p.Append(vf.Regexp(regexp.MustCompile(t.Value)))
	case name:
		p.Append(p.identifier(t))
	default:
		panic(badSyntax(t, exElement))
	}
}

func arrayType(p *pcoreParser) dgo.Value {
	if p.PeekToken().Type == '[' {
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Array`, 1, 3)

		var et dgo.Type
		var min int64 = 0
		var max int64 = math.MaxInt64
		argc := args.Len()
		switch argc {
		case 1, 2:
			if i, ok := getIfInt(args, 0, 0); ok {
				min = i
				if argc > 1 {
					max = getInt(`Array`, args, 1, math.MaxInt64)
				}
			} else {
				et = args.Arg(`Array`, 0, typ.Type).(dgo.Type)
				if argc > 1 {
					min = getInt(`Array`, args, 1, 0)
				}
			}
		case 3:
			et = args.Arg(`Array`, 0, typ.Type).(dgo.Type)
			min = getInt(`Array`, args, 1, 0)
			max = getInt(`Array`, args, 2, math.MaxInt64)
		}
		return tf.Array(et, min, max)
	}
	return typ.Array
}

func binaryType(p *pcoreParser) dgo.Value {
	if p.PeekToken().Type == '[' {
		// get size arguments
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Binary`, 1, 2)
		min := getInt(`Binary`, args, 0, 0)
		max := int64(math.MaxInt64)
		if args.Len() > 1 {
			max = getInt(`Binary`, args, 1, max)
		}
		return tf.Binary(min, max)
	}
	return typ.Binary
}

func callableType(p *pcoreParser) dgo.Value {
	var rt dgo.Type
	if p.PeekToken().Type == '[' {
		// get size arguments
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Callable`, 1, 3)
		argc := args.Len()
		first := args.Get(0)
		if params, ok := first.(dgo.TupleType); ok {
			var blockType dgo.Type
			if argc > 1 {
				blockType, ok = args.Get(1).(dgo.Type)
				if argc > 2 {
					rt, ok = args.Get(2).(dgo.Type)
				}
			}

			if ok {
				if blockType != nil {
					params = paramsWithBlock(params, blockType)
				}
				return tf.Function(params, returnType(rt))
			}
		}

		if argc == 1 || argc == 2 {
			// check for [[params, block], return]
			if iv, ok := first.(dgo.Array); ok {
				if argc == 2 {
					rt = args.Arg(`Callable`, 1, typ.Type).(dgo.Type)
				}
				argc = iv.Len()
				args = vf.ArgumentsFromArray(iv)
			}
		}
		params := tupleFromArgs(args)
		if rt == nil {
			rt = typ.Any
		}
		return tf.Function(params, returnType(rt))
	}
	return typ.Function
}

func enumType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Enum`, 1, math.MaxInt64)
		n := args.Len()
		last := args.Get(n - 1)
		ci := false
		if b, ok := last.(dgo.Boolean); ok {
			ci = b.GoBool()
			n--
		}
		ss := make([]string, n)
		for i := 0; i < n; i++ {
			ss[i] = args.Arg(`Enum`, i, typ.String).(dgo.String).GoString()
		}
		if ci {
			return tf.CiEnum(ss...)
		}
		return tf.Enum(ss...)
	}
	return typ.String
}

func floatType(p *pcoreParser) dgo.Value {
	if p.PeekToken().Type == '[' {
		// get size arguments
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Float`, 1, 2)
		from := getFloat(`Float`, args, 0, -math.MaxFloat64)
		to := math.MaxFloat64
		if args.Len() == 2 {
			to = getFloat(`Float`, args, 1, math.MaxFloat64)
		}
		return tf.FloatRange(from, to, true)
	}
	return typ.Float
}

func hashType(p *pcoreParser) dgo.Value {
	if p.PeekToken().Type == '[' {
		// get size arguments
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Hash`, 1, 4)
		switch args.Len() {
		case 1:
			return tf.Map(getInt(`Hash`, args, 0, 0))
		case 2:
			if i, ok := getIfInt(args, 0, 0); ok {
				return tf.Map(i, getInt(`Hash`, args, 1, math.MaxInt64))
			}
			return tf.Map(args.Arg(`Hash`, 0, typ.Type), args.Arg(`Hash`, 1, typ.Type))
		case 3:
			return tf.Map(args.Arg(`Hash`, 0, typ.Type), args.Arg(`Hash`, 1, typ.Type), getInt(`Hash`, args, 2, 0))
		default:
			return tf.Map(
				args.Arg(`Hash`, 0, typ.Type),
				args.Arg(`Hash`, 1, typ.Type),
				getInt(`Hash`, args, 2, 0),
				getInt(`Hash`, args, 3, math.MaxInt64))
		}
	}
	return typ.Map
}

func integerType(p *pcoreParser) dgo.Value {
	if p.PeekToken().Type == '[' {
		// get size arguments
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Integer`, 1, 2)
		from := getInt(`Integer`, args, 0, math.MinInt64)
		to := int64(math.MaxInt64)
		if args.Len() == 2 {
			to = getInt(`Integer`, args, 1, math.MaxInt64)
		}
		return tf.IntegerRange(from, to, true)
	}
	return typ.Integer
}

func numberType(p *pcoreParser) dgo.Value {
	tp := floatType(p).(dgo.FloatRangeType)
	max := int64(math.MaxInt64)
	if tp.Max() != math.MaxFloat64 {
		max = int64(tp.Max())
	}
	return tf.AnyOf(tf.IntegerRange(int64(tp.Min()), max, true), tp)
}

func notUndefType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`NotUndef`, 1, 1)
		tp := args.Arg(`Type`, 0, typ.Type).(dgo.Type)
		if tp.Assignable(typ.Nil) {
			tp = tf.AllOf(tf.Not(typ.Nil), tp)
		}
		return tp
	}
	return tf.Not(typ.Nil)
}

func optionalType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Optional`, 1, 1)
		v := args.Arg(`Optional`, 0, tf.AnyOf(typ.String, typ.Type))
		var tp dgo.Type
		if s, ok := v.(dgo.String); ok {
			tp = s.Type()
		} else {
			tp = v.(dgo.Type)
		}
		if !tp.Assignable(typ.Nil) {
			tp = tf.AnyOf(typ.Nil, tp)
		}
		return tp
	}
	return typ.Any
}

func sensitiveType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		return tf.Sensitive(p.PopLast().(dgo.Array).InterfaceSlice()...)
	}
	return typ.Sensitive
}

func tupleFromArgs(args dgo.Array) dgo.TupleType {
	l := args.Len()
	var min int64
	var max int64
	if n, ok := getIfInt(args, l-1, math.MaxInt64); ok {
		args.Pop()
		l--
		if min, ok = getIfInt(args, l-1, int64(l-2)); ok {
			args.Pop()
			l--
			max = n
		} else {
			if n == math.MaxInt64 {
				min = 0
			} else {
				min = n
			}
		}
	} else {
		min = int64(l)
		max = min
	}
	ta := args.Map(func(v dgo.Value) interface{} {
		if tp, ok := v.(dgo.Type); ok {
			return tp
		}
		panic(tf.IllegalAssignment(typ.Type, v))
	})
	if min == int64(l) && min == max {
		return tf.Tuple(ta.InterfaceSlice()...)
	}
	tps := ta.InterfaceSlice()
	return tf.VariadicTuple(tps...)
}

func tupleType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		return tupleFromArgs(p.PopLast().(dgo.Array))
	}
	return typ.Array
}

func patternType(p *pcoreParser) dgo.Value {
	if p.PeekToken().Type == '[' {
		// get size arguments
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Pattern`, 1, math.MaxInt64)
		return tf.AnyOf(
			args.Map(func(v dgo.Value) interface{} { return tf.Pattern(toRegexp(v)) }).InterfaceSlice()...)
	}
	return typ.String
}

func regexpType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Regexp`, 1, 1)
		return vf.Regexp(toRegexp(args.Get(0))).Type()
	}
	return typ.Regexp
}

func stringType(p *pcoreParser) dgo.Value {
	if p.PeekToken().Type == '[' {
		// get size arguments
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`String`, 1, 2)
		return tf.String(args.InterfaceSlice()...)
	}
	return typ.String
}

func structType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Struct`, 1, 1)
		m := args.Arg(`Struct`, 0, typ.Map).(dgo.Map)
		entries := make([]dgo.StructMapEntry, m.Len())
		i := 0
		m.EachEntry(func(e dgo.MapEntry) {
			k, optional := optionalValue(e.Key())
			entries[i] = tf.StructMapEntry(k, e.Value(), !optional)
			i++
		})
		return tf.StructMap(false, entries...)
	}
	return typ.Map
}

func typeType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		args := vf.ArgumentsFromArray(p.PopLast().(dgo.Array))
		args.AssertSize(`Type`, 1, 1)
		return tf.Meta(args.Arg(`Type`, 0, typ.Type).(dgo.Type))
	}
	return typ.Type
}

func variantType(p *pcoreParser) dgo.Value {
	tt := p.PeekToken().Type
	if tt == '[' {
		p.NextToken()
		p.array()
		args := p.PopLast().(dgo.Array)
		return tf.AnyOf(args.Map(func(v dgo.Value) interface{} {
			if tp, ok := v.(dgo.Type); ok {
				return tp
			}
			panic(tf.IllegalAssignment(typ.Type, v))
		}).InterfaceSlice()...)
	}
	return typ.AnyOf
}

var identifierToTypeMap = map[string]dgo.Value{
	`Any`:       typ.Any,
	`Boolean`:   typ.Boolean,
	`False`:     typ.False,
	`Timestamp`: typ.Time,
	`True`:      typ.True,
	`Undef`:     typ.Nil,
}

var identifierToFuncMap map[string]func(p *pcoreParser) dgo.Value

func init() {
	identifierToFuncMap = map[string]func(p *pcoreParser) dgo.Value{
		`Array`:     arrayType,
		`Binary`:    binaryType,
		`Callable`:  callableType,
		`Float`:     floatType,
		`Enum`:      enumType,
		`Hash`:      hashType,
		`Integer`:   integerType,
		`Number`:    numberType,
		`NotUndef`:  notUndefType,
		`Optional`:  optionalType,
		`Pattern`:   patternType,
		`Regexp`:    regexpType,
		`Sensitive`: sensitiveType,
		`String`:    stringType,
		`Struct`:    structType,
		`Tuple`:     tupleType,
		`Type`:      typeType,
		`Variant`:   variantType,
	}
}

func (p *pcoreParser) knownType(t *parser.Token) dgo.Value {
	tp, ok := identifierToTypeMap[t.Value]
	if ok {
		return tp
	}
	fn, ok := identifierToFuncMap[t.Value]
	if ok {
		return fn(p)
	}
	return nil
}

func (p *pcoreParser) identifier(t *parser.Token) dgo.Value {
	if tp := p.knownType(t); tp != nil {
		return tp
	}
	return p.namedType(t.Value)
}

func (p *pcoreParser) namedType(n string) dgo.Value {
	ns := strings.Split(n, `::`)
	for i := range ns {
		pn := ns[i]
		ns[i] = strings.ToLower(pn[:1]) + pn[1:]
	}

	n = toDgoName(n)
	tp := p.aliasReference(n)
	if p.PeekToken().Type == '(' {
		p.NextToken()
		p.params()
		tp = parser.NewCall(tp, vf.ArgumentsFromArray(p.PopLast().(dgo.Array)))
	}
	return tp
}

func (p *pcoreParser) aliasReference(n string) dgo.Type {
	if tp := tf.Named(n); tp != nil {
		if p.PeekToken().Type == '[' {
			p.NextToken()
			p.array()
			tp = tf.Parameterized(tp, p.PopLast().(dgo.Array))
		}
		return tp
	}
	vn := vf.String(n)
	if tp := p.AliasMap().GetType(vn); tp != nil {
		return tp
	}
	return parser.NewAlias(vn)
}

// convertMapEntries converts consecutive MapEntry elements found in an array to a Map. This is
// to permit the x => y notation inside an Array.
func convertMapEntries(av dgo.Array) dgo.Array {
	es := make([]dgo.Value, 0, av.Len())

	var en dgo.Map
	av.Each(func(v dgo.Value) {
		if he, ok := v.(dgo.MapEntry); ok {
			if en == nil {
				en = vf.MutableMap(he.Key(), he.Value())
			} else {
				en.Put(he.Key(), he.Value())
			}
		} else {
			if en != nil {
				es = append(es, en)
				en = nil
			}
			es = append(es, v)
		}
	})
	if en != nil {
		es = append(es, en)
	}
	return vf.WrapSlice(es)
}

func getIfInt(args dgo.Array, arg int, dflt int64) (int64, bool) {
	switch v := args.Get(arg).(type) {
	case dgo.Integer:
		return v.GoInt(), true
	case dgo.String:
		if v.GoString() == `default` {
			return dflt, true
		}
	}
	return 0, false
}

func getInt(fn string, args dgo.Arguments, arg int, dflt int64) int64 {
	switch v := args.Get(arg).(type) {
	case dgo.Integer:
		return v.GoInt()
	case dgo.String:
		if v.GoString() == `default` {
			return dflt
		}
	}
	// Trigger panic since argument isn't an integer
	args.Arg(fn, arg, typ.Integer)
	return 0
}

func getFloat(fn string, args dgo.Arguments, arg int, dflt float64) float64 {
	switch v := args.Get(arg).(type) {
	case dgo.Integer:
		return float64(v.GoInt())
	case dgo.Float:
		return v.GoFloat()
	case dgo.String:
		if v.GoString() == `default` {
			return dflt
		}
	}
	// Trigger panic since argument isn't an integer
	args.Arg(fn, arg, typ.Number)
	return 0
}

func optionalValue(v dgo.Value) (dgo.Value, bool) {
	if tt, ok := v.(dgo.TernaryType); ok {
		if tt.Operator() == dgo.OpOr {
			ops := tt.Operands()
			if ops.Len() == 2 {
				if typ.Nil.Equals(ops.Get(0)) {
					return ops.Get(1), true
				}
				if typ.Nil.Equals(ops.Get(1)) {
					return ops.Get(0), true
				}
			}
		}
	}
	return v, false
}

func paramsWithBlock(params dgo.TupleType, block dgo.Type) dgo.TupleType {
	tva := params.ElementTypes().InterfaceSlice()
	ln := len(tva) - 1
	tvc := make([]interface{}, ln+1)
	copy(tvc, tva)
	if params.Variadic() {
		// Move variadic type to new last pos
		tvc[ln] = tvc[ln-1]
		tvc[ln-1] = block
		params = tf.VariadicTuple(tvc...)
	} else {
		tvc[ln] = block
		params = tf.Tuple(tvc...)
	}
	return params
}

func returnType(rt dgo.Type) dgo.TupleType {
	var ret dgo.TupleType
	if rt == nil {
		ret = typ.EmptyTuple
	} else {
		ret = tf.Tuple(rt)
	}
	return ret
}

func toDgoName(pcoreName string) string {
	switch pcoreName {
	case `Integer`:
		return `int`
	case `Boolean`:
		return `bool`
	}
	ns := strings.Split(pcoreName, `::`)
	for i := range ns {
		n := ns[i]
		ns[i] = strings.ToLower(n[:1]) + n[1:]
	}
	return strings.Join(ns, `.`)
}

func toRegexp(v dgo.Value) (rx *regexp.Regexp) {
	switch v := v.(type) {
	case dgo.String:
		rx = regexp.MustCompile(v.GoString())
	case dgo.Regexp:
		rx = v.GoRegexp()
	default:
		panic(tf.IllegalAssignment(tf.AnyOf(typ.String, typ.Regexp), v))
	}
	return
}
