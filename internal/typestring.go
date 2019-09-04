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

func writeTernary(typ dgo.Type, prio int, op string, opPrio int, sb *strings.Builder) {
	if prio >= orPrio {
		sb.WriteByte('(')
	}
	Join(typ.(dgo.TernaryType).Operands(), op, opPrio, sb)
	if prio >= orPrio {
		sb.WriteByte(')')
	}
}

func writeElementsExact(typ dgo.Type, prio int, sb *strings.Builder) {
	if prio >= andPrio {
		sb.WriteByte('(')
	}
	JoinValueTypes(typ.(dgo.Iterable), `&`, andPrio, sb)
	if prio >= andPrio {
		sb.WriteByte(')')
	}
}

var simpleTypes = map[dgo.TypeIdentifier]string{
	dgo.IdAny:     `any`,
	dgo.IdArray:   `[]any`,
	dgo.IdTrue:    `true`,
	dgo.IdFalse:   `false`,
	dgo.IdBoolean: `bool`,
	dgo.IdNil:     `nil`,
	dgo.IdError:   `error`,
	dgo.IdFloat:   `float`,
	dgo.IdMap:     `map[any]any`,
	dgo.IdRegexp:  `regexp`,
	dgo.IdString:  `string`,
	dgo.IdInteger: `int`,
}

type typeToString func(typ dgo.Type, prio int, sb *strings.Builder)

var complexTypes map[dgo.TypeIdentifier]typeToString

func init() {
	complexTypes = map[dgo.TypeIdentifier]typeToString{
		dgo.IdAnyOf: func(typ dgo.Type, prio int, sb *strings.Builder) { writeTernary(typ, prio, `|`, orPrio, sb) },
		dgo.IdOneOf: func(typ dgo.Type, prio int, sb *strings.Builder) { writeTernary(typ, prio, `^`, xorPrio, sb) },
		dgo.IdAllOf: func(typ dgo.Type, prio int, sb *strings.Builder) { writeTernary(typ, prio, `&`, andPrio, sb) },
		dgo.IdArrayExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteByte('{')
			JoinValueTypes(typ.(dgo.ExactType).Value().(dgo.Iterable), `,`, commaPrio, sb)
			sb.WriteByte('}')
		},
		dgo.IdElementsExact:  func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.IdMapKeysExact:   func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.IdMapValuesExact: func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.IdTuple: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteByte('{')
			Join(typ.(dgo.TupleType).ElementTypes(), `,`, commaPrio, sb)
			sb.WriteByte('}')
		},
		dgo.IdArrayElementSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
			at := typ.(dgo.ArrayType)
			if at.Unbounded() {
				sb.WriteString(`[]`)
			} else {
				sb.WriteByte('[')
				writeSizeBoundaries(int64(at.Min()), int64(at.Max()), sb)
				sb.WriteByte(']')
			}
			buildTypeString(at.ElementType(), typePrio, sb)
		},
		dgo.IdMapExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteByte('{')
			JoinValueTypes(typ.(dgo.ExactType).Value().(dgo.Map).Entries(), `,`, commaPrio, sb)
			sb.WriteByte('}')
		},
		dgo.IdStruct: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteByte('{')
			Join(typ.(dgo.StructType).Entries(), `,`, commaPrio, sb)
			sb.WriteByte('}')
		},
		dgo.IdMapEntry: func(typ dgo.Type, prio int, sb *strings.Builder) {
			me := typ.(dgo.MapEntryType)
			buildTypeString(me.KeyType(), commaPrio, sb)
			if !me.Required() {
				sb.WriteByte('?')
			}
			sb.WriteByte(':')
			buildTypeString(me.ValueType(), commaPrio, sb)
		},
		dgo.IdMapEntryExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			me := typ.(dgo.ExactType).Value().(dgo.MapEntry)
			buildTypeString(me.Key().Type(), commaPrio, sb)
			sb.WriteByte(':')
			buildTypeString(me.Value().Type(), commaPrio, sb)
		},
		dgo.IdMapSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
			at := typ.(dgo.MapType)
			sb.WriteString(`map[`)
			buildTypeString(at.KeyType(), commaPrio, sb)
			if !at.Unbounded() {
				sb.WriteByte(',')
				writeSizeBoundaries(int64(at.Min()), int64(at.Max()), sb)
			}
			sb.WriteByte(']')
			buildTypeString(at.ValueType(), typePrio, sb)
		},
		dgo.IdFloatExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(util.Ftoa(typ.(dgo.ExactType).Value().(Float).GoFloat()))
		},
		dgo.IdFloatRange: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.FloatRangeType)
			writeFloatRange(st.Min(), st.Max(), sb)
		},
		dgo.IdIntegerExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(typ.(dgo.ExactType).Value().(fmt.Stringer).String())
		},
		dgo.IdIntegerRange: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.IntegerRangeType)
			writeIntRange(st.Min(), st.Max(), sb)
		},
		dgo.IdRegexpExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(`regexp[`)
			sb.WriteString(strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
			sb.WriteByte(']')
		},
		dgo.IdStringExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
		},
		dgo.IdStringPattern: func(typ dgo.Type, prio int, sb *strings.Builder) {
			RegexpSlashQuote(sb, typ.(dgo.ExactType).Value().(fmt.Stringer).String())
		},
		dgo.IdStringSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.StringType)
			sb.WriteString(`string`)
			if !st.Unbounded() {
				sb.WriteByte('[')
				writeSizeBoundaries(int64(st.Min()), int64(st.Max()), sb)
				sb.WriteByte(']')
			}
		},
		dgo.IdNot: func(typ dgo.Type, prio int, sb *strings.Builder) {
			nt := typ.(dgo.UnaryType)
			sb.WriteByte('!')
			buildTypeString(nt.Operand(), typePrio, sb)
		},
		dgo.IdNative: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(typ.(dgo.NativeType).GoType().String())
		},
		dgo.IdMeta: func(typ dgo.Type, prio int, sb *strings.Builder) {
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
		},
	}
}

func buildTypeString(typ dgo.Type, prio int, sb *strings.Builder) {
	ti := typ.TypeIdentifier()
	if s, ok := simpleTypes[ti]; ok {
		sb.WriteString(s)
	} else if f, ok := complexTypes[ti]; ok {
		f(typ, prio, sb)
	} else {
		panic(fmt.Errorf(`type identifier %d is not handled`, ti))
	}
}
