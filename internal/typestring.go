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
			util.WriteString(sb, s)
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
			util.WriteString(sb, s)
		}
		buildTypeString(v.Type(), prio, sb)
	})
}

func writeSizeBoundaries(min, max int64, sb *strings.Builder) {
	util.WriteString(sb, strconv.FormatInt(min, 10))
	if max != math.MaxInt64 {
		util.WriteByte(sb, ',')
		util.WriteString(sb, strconv.FormatInt(max, 10))
	}
}

func writeIntRange(min, max int64, sb *strings.Builder) {
	if min != math.MinInt64 {
		util.WriteString(sb, strconv.FormatInt(min, 10))
	}
	util.WriteString(sb, `..`)
	if max != math.MaxInt64 {
		util.WriteString(sb, strconv.FormatInt(max, 10))
	}
}

func writeFloatRange(min, max float64, sb *strings.Builder) {
	if min != -math.MaxFloat64 {
		util.WriteString(sb, util.Ftoa(min))
	}
	util.WriteString(sb, `..`)
	if max != math.MaxFloat64 {
		util.WriteString(sb, util.Ftoa(max))
	}
}

func writeTernary(typ dgo.Type, prio int, op string, opPrio int, sb *strings.Builder) {
	if prio >= orPrio {
		util.WriteByte(sb, '(')
	}
	join(typ.(dgo.TernaryType).Operands(), op, opPrio, sb)
	if prio >= orPrio {
		util.WriteByte(sb, ')')
	}
}

func writeElementsExact(typ dgo.Type, prio int, sb *strings.Builder) {
	if prio >= andPrio {
		util.WriteByte(sb, '(')
	}
	joinValueTypes(typ.(dgo.Iterable), `&`, andPrio, sb)
	if prio >= andPrio {
		util.WriteByte(sb, ')')
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
			util.WriteByte(sb, '{')
			joinValueTypes(typ.(dgo.ExactType).Value().(dgo.Iterable), `,`, commaPrio, sb)
			util.WriteByte(sb, '}')
		},
		dgo.TiElementsExact:  func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.TiMapKeysExact:   func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.TiMapValuesExact: func(typ dgo.Type, prio int, sb *strings.Builder) { writeElementsExact(typ, prio, sb) },
		dgo.TiTuple: func(typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteByte(sb, '{')
			join(typ.(dgo.TupleType).ElementTypes(), `,`, commaPrio, sb)
			util.WriteByte(sb, '}')
		},
		dgo.TiArrayElementSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
			at := typ.(dgo.ArrayType)
			if at.Unbounded() {
				util.WriteString(sb, `[]`)
			} else {
				util.WriteByte(sb, '[')
				writeSizeBoundaries(int64(at.Min()), int64(at.Max()), sb)
				util.WriteByte(sb, ']')
			}
			buildTypeString(at.ElementType(), typePrio, sb)
		},
		dgo.TiMapExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteByte(sb, '{')
			joinValueTypes(typ.(dgo.ExactType).Value().(dgo.Map).Entries(), `,`, commaPrio, sb)
			util.WriteByte(sb, '}')
		},
		dgo.TiStruct: func(typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteByte(sb, '{')
			join(typ.(dgo.StructType).Entries(), `,`, commaPrio, sb)
			util.WriteByte(sb, '}')
		},
		dgo.TiMapEntry: func(typ dgo.Type, prio int, sb *strings.Builder) {
			me := typ.(dgo.MapEntryType)
			buildTypeString(me.KeyType(), commaPrio, sb)
			if !me.Required() {
				util.WriteByte(sb, '?')
			}
			util.WriteByte(sb, ':')
			buildTypeString(me.ValueType(), commaPrio, sb)
		},
		dgo.TiMapEntryExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			me := typ.(dgo.ExactType).Value().(dgo.MapEntry)
			buildTypeString(me.Key().Type(), commaPrio, sb)
			util.WriteByte(sb, ':')
			buildTypeString(me.Value().Type(), commaPrio, sb)
		},
		dgo.TiMapSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
			at := typ.(dgo.MapType)
			util.WriteString(sb, `map[`)
			buildTypeString(at.KeyType(), commaPrio, sb)
			if !at.Unbounded() {
				util.WriteByte(sb, ',')
				writeSizeBoundaries(int64(at.Min()), int64(at.Max()), sb)
			}
			util.WriteByte(sb, ']')
			buildTypeString(at.ValueType(), typePrio, sb)
		},
		dgo.TiFloatExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, util.Ftoa(typ.(dgo.ExactType).Value().(floatVal).GoFloat()))
		},
		dgo.TiFloatRange: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.FloatRangeType)
			writeFloatRange(st.Min(), st.Max(), sb)
		},
		dgo.TiIntegerExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, typ.(dgo.ExactType).Value().(fmt.Stringer).String())
		},
		dgo.TiIntegerRange: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.IntegerRangeType)
			writeIntRange(st.Min(), st.Max(), sb)
		},
		dgo.TiRegexpExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, `regexp[`)
			util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
			util.WriteByte(sb, ']')
		},
		dgo.TiStringExact: func(typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
		},
		dgo.TiStringPattern: func(typ dgo.Type, prio int, sb *strings.Builder) {
			RegexpSlashQuote(sb, typ.(dgo.ExactType).Value().(fmt.Stringer).String())
		},
		dgo.TiStringSized: func(typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.StringType)
			util.WriteString(sb, `string`)
			if !st.Unbounded() {
				util.WriteByte(sb, '[')
				writeSizeBoundaries(int64(st.Min()), int64(st.Max()), sb)
				util.WriteByte(sb, ']')
			}
		},
		dgo.TiNot: func(typ dgo.Type, prio int, sb *strings.Builder) {
			nt := typ.(dgo.UnaryType)
			util.WriteByte(sb, '!')
			buildTypeString(nt.Operand(), typePrio, sb)
		},
		dgo.TiNative: func(typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, typ.(dgo.NativeType).GoType().String())
		},
		dgo.TiMeta: func(typ dgo.Type, prio int, sb *strings.Builder) {
			nt := typ.(dgo.UnaryType)
			util.WriteString(sb, `type`)
			if op := nt.Operand(); op != DefaultAnyType {
				if op == nil {
					util.WriteString(sb, `[type]`)
				} else {
					util.WriteByte(sb, '[')
					buildTypeString(op, prio, sb)
					util.WriteByte(sb, ']')
				}
			}
		},
	}
}

func buildTypeString(typ dgo.Type, prio int, sb *strings.Builder) {
	ti := typ.TypeIdentifier()
	if s, ok := simpleTypes[ti]; ok {
		util.WriteString(sb, s)
	} else if f, ok := complexTypes[ti]; ok {
		f(typ, prio, sb)
	} else {
		panic(fmt.Errorf(`type identifier %d is not handled`, ti))
	}
}
