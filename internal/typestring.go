package internal

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

const (
	commaPrio = iota
	orPrio
	xorPrio
	andPrio
	typePrio
)

func TypeString(typ dgo.Type) string {
	bld := &strings.Builder{}
	buildTypeString(typ, 0, bld)
	return bld.String()
}

func Join(v dgo.Array, s string, prio int, sb *strings.Builder) {
	v.EachWithIndex(func(v dgo.Value, i int) {
		if i > 0 {
			sb.WriteString(s)
		}
		buildTypeString(v.(dgo.Type), prio, sb)
	})
}

func JoinValueTypes(v dgo.Iterable, s string, prio int, sb *strings.Builder) {
	first := true
	v.Each(func(v dgo.Value) {
		if first {
			first = false
		} else {
			sb.WriteString(s)
		}
		buildTypeString(v.Type(), prio, sb)
	})
}

func writeSizeBoundaries(min, max int64, sb *strings.Builder) {
	sb.WriteString(strconv.FormatInt(min, 10))
	if max != math.MaxInt64 {
		sb.WriteByte(',')
		sb.WriteString(strconv.FormatInt(max, 10))
	}
}

func writeIntRange(min, max int64, sb *strings.Builder) {
	if min != math.MinInt64 {
		sb.WriteString(strconv.FormatInt(min, 10))
	}
	sb.WriteString(`..`)
	if max != math.MaxInt64 {
		sb.WriteString(strconv.FormatInt(max, 10))
	}
}

func writeFloatRange(min, max float64, sb *strings.Builder) {
	if min != -math.MaxFloat64 {
		sb.WriteString(util.Ftoa(min))
	}
	sb.WriteString(`..`)
	if max != math.MaxFloat64 {
		sb.WriteString(util.Ftoa(max))
	}
}

func buildTypeString(typ dgo.Type, prio int, sb *strings.Builder) {
	switch typ.TypeIdentifier() {
	case dgo.IdAny:
		sb.WriteString(`any`)
	case dgo.IdAnyOf:
		if prio >= orPrio {
			sb.WriteByte('(')
		}
		Join(typ.(dgo.TernaryType).Operands(), `|`, orPrio, sb)
		if prio >= orPrio {
			sb.WriteByte(')')
		}
	case dgo.IdOneOf:
		if prio >= xorPrio {
			sb.WriteByte('(')
		}
		Join(typ.(dgo.TernaryType).Operands(), `^`, xorPrio, sb)
		if prio >= xorPrio {
			sb.WriteByte(')')
		}
	case dgo.IdAllOf:
		if prio >= andPrio {
			sb.WriteByte('(')
		}
		Join(typ.(dgo.TernaryType).Operands(), `&`, andPrio, sb)
		if prio >= andPrio {
			sb.WriteByte(')')
		}
	case dgo.IdArray:
		sb.WriteString(`[]any`)
	case dgo.IdArrayExact:
		sb.WriteByte('{')
		JoinValueTypes(typ.(dgo.ExactType).Value().(dgo.Iterable), `,`, commaPrio, sb)
		sb.WriteByte('}')
	case dgo.IdElementsExact, dgo.IdMapKeysExact, dgo.IdMapValuesExact:
		if prio >= andPrio {
			sb.WriteByte('(')
		}
		JoinValueTypes(typ.(dgo.Iterable), `&`, andPrio, sb)
		if prio >= andPrio {
			sb.WriteByte(')')
		}
	case dgo.IdTuple:
		sb.WriteByte('{')
		Join(typ.(dgo.TupleType).ElementTypes(), `,`, commaPrio, sb)
		sb.WriteByte('}')
	case dgo.IdArrayElementSized:
		at := typ.(dgo.ArrayType)
		if at.Unbounded() {
			sb.WriteString(`[]`)
		} else {
			sb.WriteByte('[')
			writeSizeBoundaries(int64(at.Min()), int64(at.Max()), sb)
			sb.WriteByte(']')
		}
		buildTypeString(at.ElementType(), typePrio, sb)
	case dgo.IdMap:
		sb.WriteString(`map[any]any`)
	case dgo.IdMapExact:
		sb.WriteByte('{')
		JoinValueTypes(typ.(dgo.ExactType).Value().(dgo.Map).Entries(), `,`, commaPrio, sb)
		sb.WriteByte('}')
	case dgo.IdStruct:
		sb.WriteByte('{')
		Join(typ.(dgo.StructType).Entries(), `,`, commaPrio, sb)
		sb.WriteByte('}')
	case dgo.IdMapEntry:
		me := typ.(dgo.MapEntryType)
		buildTypeString(me.KeyType(), commaPrio, sb)
		if !me.Required() {
			sb.WriteByte('?')
		}
		sb.WriteByte(':')
		buildTypeString(me.ValueType(), commaPrio, sb)
	case dgo.IdMapEntryExact:
		me := typ.(dgo.ExactType).Value().(dgo.MapEntry)
		buildTypeString(me.Key().Type(), commaPrio, sb)
		sb.WriteByte(':')
		buildTypeString(me.Value().Type(), commaPrio, sb)
	case dgo.IdMapSized:
		at := typ.(dgo.MapType)
		sb.WriteString(`map[`)
		buildTypeString(at.KeyType(), commaPrio, sb)
		if !at.Unbounded() {
			sb.WriteByte(',')
			writeSizeBoundaries(int64(at.Min()), int64(at.Max()), sb)
		}
		sb.WriteByte(']')
		buildTypeString(at.ValueType(), typePrio, sb)
	case dgo.IdTrue:
		sb.WriteString(`true`)
	case dgo.IdFalse:
		sb.WriteString(`false`)
	case dgo.IdBoolean:
		sb.WriteString(`bool`)
	case dgo.IdNil:
		sb.WriteString(`nil`)
	case dgo.IdError:
		sb.WriteString(`error`)
	case dgo.IdFloat:
		sb.WriteString(`float`)
	case dgo.IdFloatExact:
		sb.WriteString(util.Ftoa(typ.(dgo.ExactType).Value().(Float).GoFloat()))
	case dgo.IdFloatRange:
		st := typ.(dgo.FloatRangeType)
		writeFloatRange(st.Min(), st.Max(), sb)
	case dgo.IdInteger:
		sb.WriteString(`int`)
	case dgo.IdIntegerExact:
		sb.WriteString(typ.(dgo.ExactType).Value().(fmt.Stringer).String())
	case dgo.IdIntegerRange:
		st := typ.(dgo.IntegerRangeType)
		writeIntRange(st.Min(), st.Max(), sb)
	case dgo.IdRegexp:
		sb.WriteString(`regexp`)
	case dgo.IdRegexpExact:
		sb.WriteString(`regexp[`)
		sb.WriteString(strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
		sb.WriteByte(']')
	case dgo.IdString:
		sb.WriteString(`string`)
	case dgo.IdStringExact:
		sb.WriteString(strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
	case dgo.IdStringPattern:
		RegexpSlashQuote(sb, typ.(dgo.ExactType).Value().(fmt.Stringer).String())
	case dgo.IdStringSized:
		st := typ.(dgo.StringType)
		sb.WriteString(`string`)
		if !st.Unbounded() {
			sb.WriteByte('[')
			writeSizeBoundaries(int64(st.Min()), int64(st.Max()), sb)
			sb.WriteByte(']')
		}
	case dgo.IdNot:
		nt := typ.(dgo.UnaryType)
		sb.WriteByte('!')
		buildTypeString(nt.Operand(), typePrio, sb)
	case dgo.IdNative:
		sb.WriteString(typ.(dgo.NativeType).GoType().String())
	case dgo.IdMeta:
		nt := typ.(dgo.UnaryType)
		sb.WriteString(`type`)
		if op := nt.Operand(); op != DefaultAnyType {
			if op == nil {
				sb.WriteString(`[type]`)
			} else {
				sb.WriteByte('[')
				buildTypeString(op, prio, sb)
				sb.WriteByte(']')
			}
		}
	}
}
