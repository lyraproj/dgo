package stringer

import (
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/util"
)

const (
	commaPrio = iota
	orPrio
	xorPrio
	andPrio
	typePrio
)

type typeToString func(sb *typeBuilder, typ dgo.Type, prio int)

type typeBuilder struct {
	io.Writer
	aliasMap dgo.AliasMap
	seen     []dgo.Value
}

// TypeString produces a string with the go-like syntax for the given type.
func TypeString(typ dgo.Type) string {
	return TypeStringWithAliasMap(typ, internal.DefaultAliases())
}

// TypeStringOn produces a string with the go-like syntax for the given type onto the given io.Writer.
func TypeStringOn(typ dgo.Type, w io.Writer) {
	newTypeBuilder(w, internal.DefaultAliases()).buildTypeString(typ, 0)
}

// TypeStringWithAliasMap produces a string with the go-like syntax for the given type.
func TypeStringWithAliasMap(typ dgo.Type, am dgo.AliasMap) string {
	s := strings.Builder{}
	newTypeBuilder(&s, am).buildTypeString(typ, 0)
	return s.String()
}

func anyOf(sb *typeBuilder, typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `|`, orPrio)
}

func oneOf(sb *typeBuilder, typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `^`, xorPrio)
}

func allOf(sb *typeBuilder, typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `&`, andPrio)
}

func allOfValue(sb *typeBuilder, typ dgo.Type, prio int) {
	sb.writeTernary(typ, valueAsType, prio, `&`, andPrio)
}

func array(sb *typeBuilder, typ dgo.Type, _ int) {
	at := typ.(dgo.ArrayType)
	if at.Unbounded() {
		util.WriteString(sb, `[]`)
	} else {
		util.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(at.Min()), int64(at.Max()))
		util.WriteByte(sb, ']')
	}
	sb.buildTypeString(at.ElementType(), typePrio)
}

func arrayExact(sb *typeBuilder, typ dgo.Type, _ int) {
	util.WriteByte(sb, '{')
	sb.joinValueTypes(typ.(dgo.Iterable), `,`, commaPrio)
	util.WriteByte(sb, '}')
}

func binary(sb *typeBuilder, typ dgo.Type, _ int) {
	st := typ.(dgo.BinaryType)
	util.WriteString(sb, `binary`)
	if !st.Unbounded() {
		util.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(st.Min()), int64(st.Max()))
		util.WriteByte(sb, ']')
	}
}

func stringValue(sb *typeBuilder, typ dgo.Type, _ int) {
	_, _ = sb.Write([]byte(typ.String()))
}

func exactValue(sb *typeBuilder, typ dgo.Type, _ int) {
	util.WriteString(sb, typ.(dgo.ExactType).ExactValue().String())
}

func tuple(sb *typeBuilder, typ dgo.Type, _ int) {
	sb.writeTupleArgs(typ.(dgo.TupleType), '{', '}')
}

func _map(sb *typeBuilder, typ dgo.Type, _ int) {
	at := typ.(dgo.MapType)
	util.WriteString(sb, `map[`)
	sb.buildTypeString(at.KeyType(), commaPrio)
	if !at.Unbounded() {
		util.WriteByte(sb, ',')
		sb.writeSizeBoundaries(int64(at.Min()), int64(at.Max()))
	}
	util.WriteByte(sb, ']')
	sb.buildTypeString(at.ValueType(), typePrio)
}

func mapExact(sb *typeBuilder, typ dgo.Type, _ int) {
	util.WriteByte(sb, '{')
	sb.joinValueTypes(typ.(dgo.Map), `,`, commaPrio)
	util.WriteByte(sb, '}')
}

func _struct(sb *typeBuilder, typ dgo.Type, _ int) {
	util.WriteByte(sb, '{')
	st := typ.(dgo.StructMapType)
	sb.joinStructMapEntries(st)
	if st.Additional() {
		if st.Len() > 0 {
			util.WriteByte(sb, ',')
		}
		util.WriteString(sb, `...`)
	}
	util.WriteByte(sb, '}')
}

func mapEntryExact(sb *typeBuilder, typ dgo.Type, _ int) {
	me := typ.(dgo.MapEntry)
	sb.buildTypeString(typeAsType(me.Key()), commaPrio)
	util.WriteByte(sb, ':')
	sb.buildTypeString(typeAsType(me.Value()), commaPrio)
}

func floatRange(sb *typeBuilder, typ dgo.Type, _ int) {
	st := typ.(dgo.FloatType)
	sb.writeFloatRange(st.Min(), st.Max(), st.Inclusive())
}

func integerRange(sb *typeBuilder, typ dgo.Type, _ int) {
	st := typ.(dgo.IntegerType)
	sb.writeIntRange(st.Min(), st.Max(), st.Inclusive())
}

func sensitive(sb *typeBuilder, typ dgo.Type, prio int) {
	util.WriteString(sb, `sensitive`)
	if op := typ.(dgo.UnaryType).Operand(); internal.DefaultAnyType != op {
		util.WriteByte(sb, '[')
		sb.buildTypeString(op, prio)
		util.WriteByte(sb, ']')
	}
}

func stringSized(sb *typeBuilder, typ dgo.Type, _ int) {
	st := typ.(dgo.StringType)
	util.WriteString(sb, `string`)
	if !st.Unbounded() {
		util.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(st.Min()), int64(st.Max()))
		util.WriteByte(sb, ']')
	}
}

func ciString(sb *typeBuilder, typ dgo.Type, _ int) {
	util.WriteByte(sb, '~')
	util.WriteString(sb, strconv.Quote(typ.(dgo.String).GoString()))
}

func native(sb *typeBuilder, typ dgo.Type, _ int) {
	rt := typ.(dgo.NativeType).GoType()
	util.WriteString(sb, `native`)
	if rt != nil {
		util.WriteByte(sb, '[')
		util.WriteString(sb, strconv.Quote(rt.String()))
		util.WriteByte(sb, ']')
	}
}

func not(sb *typeBuilder, typ dgo.Type, _ int) {
	nt := typ.(dgo.UnaryType)
	util.WriteByte(sb, '!')
	sb.buildTypeString(nt.Operand(), typePrio)
}

func meta(sb *typeBuilder, typ dgo.Type, prio int) {
	nt := typ.(dgo.UnaryType)
	util.WriteString(sb, `type`)
	if op := nt.Operand(); internal.DefaultAnyType != op {
		if op == nil {
			util.WriteString(sb, `[type]`)
		} else {
			util.WriteByte(sb, '[')
			sb.buildTypeString(op, prio)
			util.WriteByte(sb, ']')
		}
	}
}

func function(sb *typeBuilder, typ dgo.Type, prio int) {
	ft := typ.(dgo.FunctionType)
	util.WriteString(sb, `func`)
	sb.writeTupleArgs(ft.In(), '(', ')')
	out := ft.Out()
	if out.Len() > 0 {
		util.WriteByte(sb, ' ')
		if out.Len() == 1 && !out.Variadic() {
			sb.buildTypeString(out.ElementTypeAt(0), prio)
		} else {
			sb.writeTupleArgs(ft.Out(), '(', ')')
		}
	}
}

func named(sb *typeBuilder, typ dgo.Type, _ int) {
	nt := typ.(dgo.NamedType)
	util.WriteString(sb, nt.Name())
	if params := nt.Parameters(); params != nil {
		util.WriteByte(sb, '[')
		sb.joinValueTypes(params, `,`, commaPrio)
		util.WriteByte(sb, ']')
	}
}

func newTypeBuilder(w io.Writer, am dgo.AliasMap) *typeBuilder {
	return &typeBuilder{Writer: w, aliasMap: am}
}

var complexTypes map[dgo.TypeIdentifier]typeToString

func init() {
	complexTypes = map[dgo.TypeIdentifier]typeToString{
		dgo.TiAnyOf:         anyOf,
		dgo.TiOneOf:         oneOf,
		dgo.TiAllOf:         allOf,
		dgo.TiAllOfValue:    allOfValue,
		dgo.TiArray:         array,
		dgo.TiArrayExact:    arrayExact,
		dgo.TiBinary:        binary,
		dgo.TiBinaryExact:   stringValue,
		dgo.TiBooleanExact:  stringValue,
		dgo.TiTuple:         tuple,
		dgo.TiMap:           _map,
		dgo.TiMapExact:      mapExact,
		dgo.TiMapEntryExact: mapEntryExact,
		dgo.TiStruct:        _struct,
		dgo.TiFloatExact:    stringValue,
		dgo.TiFloatRange:    floatRange,
		dgo.TiIntegerExact:  stringValue,
		dgo.TiIntegerRange:  integerRange,
		dgo.TiRegexpExact:   stringValue,
		dgo.TiTimeExact:     stringValue,
		dgo.TiSensitive:     sensitive,
		dgo.TiStringExact:   stringValue,
		dgo.TiStringPattern: stringValue,
		dgo.TiStringSized:   stringSized,
		dgo.TiCiString:      ciString,
		dgo.TiNative:        native,
		dgo.TiNot:           not,
		dgo.TiMeta:          meta,
		dgo.TiFunction:      function,
		dgo.TiErrorExact:    stringValue,
		dgo.TiNamed:         named,
		dgo.TiNamedExact:    exactValue}
}

func (sb *typeBuilder) joinTypes(v dgo.Iterable, s string, prio int) {
	sb.joinX(v, typeAsType, s, prio)
}

func (sb *typeBuilder) joinValueTypes(v dgo.Iterable, s string, prio int) {
	sb.joinX(v, valueAsType, s, prio)
}

func (sb *typeBuilder) joinX(v dgo.Iterable, tc func(dgo.Value) dgo.Type, s string, prio int) {
	first := true
	v.Each(func(v dgo.Value) {
		if first {
			first = false
		} else {
			util.WriteString(sb, s)
		}
		sb.buildTypeString(tc(v), prio)
	})
}

func (sb *typeBuilder) joinStructMapEntries(v dgo.StructMapType) {
	first := true
	v.EachEntryType(func(e dgo.StructMapEntry) {
		if first {
			first = false
		} else {
			util.WriteByte(sb, ',')
		}
		sb.buildTypeString(e.Key().(dgo.Type), commaPrio)
		if !e.Required() {
			util.WriteByte(sb, '?')
		}
		util.WriteByte(sb, ':')
		sb.buildTypeString(e.Value().(dgo.Type), commaPrio)
	})
}

func (sb *typeBuilder) writeSizeBoundaries(min, max int64) {
	util.WriteString(sb, strconv.FormatInt(min, 10))
	if max != math.MaxInt64 {
		util.WriteByte(sb, ',')
		util.WriteString(sb, strconv.FormatInt(max, 10))
	}
}

func (sb *typeBuilder) writeIntRange(min, max int64, inclusive bool) {
	if min != math.MinInt64 {
		util.WriteString(sb, strconv.FormatInt(min, 10))
	}
	op := `...`
	if inclusive {
		op = `..`
	}
	util.WriteString(sb, op)
	if max != math.MaxInt64 {
		util.WriteString(sb, strconv.FormatInt(max, 10))
	}
}

func (sb *typeBuilder) writeFloatRange(min, max float64, inclusive bool) {
	if min != -math.MaxFloat64 {
		util.WriteString(sb, util.Ftoa(min))
	}
	op := `...`
	if inclusive {
		op = `..`
	}
	util.WriteString(sb, op)
	if max != math.MaxFloat64 {
		util.WriteString(sb, util.Ftoa(max))
	}
}

func (sb *typeBuilder) writeTupleArgs(tt dgo.TupleType, leftSep, rightSep byte) {
	es := tt.ElementTypes()
	if tt.Variadic() {
		n := es.Len() - 1
		sep := leftSep
		for i := 0; i < n; i++ {
			util.WriteByte(sb, sep)
			sep = ','
			sb.buildTypeString(es.Get(i).(dgo.Type), commaPrio)
		}
		util.WriteByte(sb, sep)
		util.WriteString(sb, `...`)
		sb.buildTypeString(es.Get(n).(dgo.Type), commaPrio)
		util.WriteByte(sb, rightSep)
	} else {
		util.WriteByte(sb, leftSep)
		sb.joinTypes(es, `,`, commaPrio)
		util.WriteByte(sb, rightSep)
	}
}

func (sb *typeBuilder) writeTernary(typ dgo.Type, tc func(dgo.Value) dgo.Type, prio int, op string, opPrio int) {
	if prio >= orPrio {
		util.WriteByte(sb, '(')
	}
	sb.joinX(typ.(dgo.TernaryType).Operands(), tc, op, opPrio)
	if prio >= orPrio {
		util.WriteByte(sb, ')')
	}
}

func (sb *typeBuilder) buildTypeString(typ dgo.Type, prio int) {
	if tn := sb.aliasMap.GetName(typ); tn != nil {
		util.WriteString(sb, tn.GoString())
		return
	}

	ti := typ.TypeIdentifier()
	if f, ok := complexTypes[ti]; ok {
		if util.RecursionHit(sb.seen, typ) {
			util.WriteString(sb, `<recursive self reference to `)
			util.WriteString(sb, ti.String())
			util.WriteString(sb, ` type>`)
			return
		}
		os := sb.seen
		sb.seen = append(sb.seen, typ)
		f(sb, typ, prio)
		sb.seen = os
	} else {
		util.WriteString(sb, ti.String())
	}
}

func typeAsType(v dgo.Value) dgo.Type {
	return v.(dgo.Type)
}

func valueAsType(v dgo.Value) dgo.Type {
	return v.Type()
}

// The use of an init() function is motivated by the internal package's need to call the TypeString function
// and the stringer package's need to use the internal package. The only way to break that dependency would be
// to merge  the two packages. There's however no way to use the internal package without the stringer package
// being initialized first so the circularity between them is harmless.
func init() {
	internal.TypeString = TypeString
	internal.TypeStringOn = TypeStringOn
}
