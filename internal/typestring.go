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

// TypeString produces a string with the go-like syntax for the given type.
func TypeString(typ dgo.Type) string {
	bld := &strings.Builder{}
	buildTypeString(typ, 0, bld)
	return bld.String()
}

func join(v dgo.Array, s string, prio int, sb *strings.Builder) {
	v.EachWithIndex(func(v dgo.Value, i int) {
		if i > 0 {
			sb.WriteString(s)
		}
		buildTypeString(v.(dgo.Type), prio, sb)
	})
}

func joinValueTypes(v dgo.Iterable, s string, prio int, sb *strings.Builder) {
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
	join(typ.(dgo.TernaryType).Operands(), op, opPrio, sb)
	if prio >= orPrio {
		sb.WriteByte(')')
	}
}

func writeElementsExact(typ dgo.Type, prio int, sb *strings.Builder) {
	if prio >= andPrio {
		sb.WriteByte('(')
	}
	joinValueTypes(typ.(dgo.Iterable), `&`, andPrio, sb)
	if prio >= andPrio {
		sb.WriteByte(')')
	}
}

var simpleTypes = map[dgo.TypeIdentifier]string{
	dgo.TiAny:     `any`,
	dgo.TiArray:   `[]any`,
	dgo.TiTrue:    `true`,
	dgo.TiFalse:   `false`,
	dgo.TiBoolean: `bool`,
	dgo.TiNil:     `nil`,
	dgo.TiError:   `error`,
	dgo.TiFloat:   `float`,
	dgo.TiMap:     `map[any]any`,
	dgo.TiRegexp:  `regexp`,
	dgo.TiString:  `string`,
	dgo.TiInteger: `int`,
}

type typeToString func(typ dgo.Type, prio int, sb *strings.Builder)

var complexTypes map[dgo.TypeIdentifier]typeToString

func init() {
	complexTypes = map[dgo.TypeIdentifier]typeToString{
		dgo.TiAnyOf: func(typ dgo.Type, prio int, sb *strings.Builder) { writeTernary(typ, prio, `|`, orPrio, sb) },
		dgo.TiOneOf: func(typ dgo.Type, prio int, sb *strings.Builder) { writeTernary(typ, prio, `^`, xorPrio, sb) },
		dgo.TiAllOf: func(typ dgo.Type, prio int, sb *strings.Builder) { writeTernary(typ, prio, `&`, andPrio, sb) },
		dgo.TiArrayExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteByte('{')
			joinValueTypes(typ.(dgo.ExactType).Value().(dgo.Iterable), `,`, commaPrio, sb)
			sb.WriteByte('}')
		},
		dgo.TiElementsExact:  func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.TiMapKeysExact:   func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.TiMapValuesExact: func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.TiTuple: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteByte('{')
			join(typ.(dgo.TupleType).ElementTypes(), `,`, commaPrio, sb)
			sb.WriteByte('}')
		},
		dgo.TiArrayElementSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
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
		dgo.TiMapExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteByte('{')
			joinValueTypes(typ.(dgo.ExactType).Value().(dgo.Map).Entries(), `,`, commaPrio, sb)
			sb.WriteByte('}')
		},
		dgo.TiStruct: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteByte('{')
			join(typ.(dgo.StructType).Entries(), `,`, commaPrio, sb)
			sb.WriteByte('}')
		},
		dgo.TiMapEntry: func(typ dgo.Type, prio int, sb *strings.Builder) {
			me := typ.(dgo.MapEntryType)
			buildTypeString(me.KeyType(), commaPrio, sb)
			if !me.Required() {
				sb.WriteByte('?')
			}
			sb.WriteByte(':')
			buildTypeString(me.ValueType(), commaPrio, sb)
		},
		dgo.TiMapEntryExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			me := typ.(dgo.ExactType).Value().(dgo.MapEntry)
			buildTypeString(me.Key().Type(), commaPrio, sb)
			sb.WriteByte(':')
			buildTypeString(me.Value().Type(), commaPrio, sb)
		},
		dgo.TiMapSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
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
		dgo.TiFloatExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(util.Ftoa(typ.(dgo.ExactType).Value().(floatVal).GoFloat()))
		},
		dgo.TiFloatRange: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.FloatRangeType)
			writeFloatRange(st.Min(), st.Max(), sb)
		},
		dgo.TiIntegerExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(typ.(dgo.ExactType).Value().(fmt.Stringer).String())
		},
		dgo.TiIntegerRange: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.IntegerRangeType)
			writeIntRange(st.Min(), st.Max(), sb)
		},
		dgo.TiRegexpExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(`regexp[`)
			sb.WriteString(strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
			sb.WriteByte(']')
		},
		dgo.TiStringExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
		},
		dgo.TiStringPattern: func(typ dgo.Type, prio int, sb *strings.Builder) {
			RegexpSlashQuote(sb, typ.(dgo.ExactType).Value().(fmt.Stringer).String())
		},
		dgo.TiStringSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.StringType)
			sb.WriteString(`string`)
			if !st.Unbounded() {
				sb.WriteByte('[')
				writeSizeBoundaries(int64(st.Min()), int64(st.Max()), sb)
				sb.WriteByte(']')
			}
		},
		dgo.TiNot: func(typ dgo.Type, prio int, sb *strings.Builder) {
			nt := typ.(dgo.UnaryType)
			sb.WriteByte('!')
			buildTypeString(nt.Operand(), typePrio, sb)
		},
		dgo.TiNative: func(typ dgo.Type, prio int, sb *strings.Builder) {
			sb.WriteString(typ.(dgo.NativeType).GoType().String())
		},
		dgo.TiMeta: func(typ dgo.Type, prio int, sb *strings.Builder) {
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
