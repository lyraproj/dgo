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
	buildTypeString(nil, typ, 0, bld)
	return bld.String()
}

func joinTypes(seen []dgo.Value, v dgo.Iterable, s string, prio int, sb *strings.Builder) {
	joinX(seen, v, typeAsType, s, prio, sb)
}

func joinValueTypes(seen []dgo.Value, v dgo.Iterable, s string, prio int, sb *strings.Builder) {
	joinX(seen, v, valueAsType, s, prio, sb)
}

func joinX(seen []dgo.Value, v dgo.Iterable, tc func(dgo.Value) dgo.Type, s string, prio int, sb *strings.Builder) {
	first := true
	v.Each(func(v dgo.Value) {
		if first {
			first = false
		} else {
			util.WriteString(sb, s)
		}
		buildTypeString(seen, tc(v), prio, sb)
	})
}

func joinStructMapEntries(seen []dgo.Value, v dgo.StructMapType, sb *strings.Builder) {
	first := true
	v.Each(func(e dgo.StructMapEntry) {
		if first {
			first = false
		} else {
			util.WriteByte(sb, ',')
		}
		buildTypeString(seen, e.Key().(dgo.Type), commaPrio, sb)
		if !e.Required() {
			util.WriteByte(sb, '?')
		}
		util.WriteByte(sb, ':')
		buildTypeString(seen, e.Value().(dgo.Type), commaPrio, sb)
	})
}

func writeSizeBoundaries(min, max int64, sb *strings.Builder) {
	util.WriteString(sb, strconv.FormatInt(min, 10))
	if max != math.MaxInt64 {
		util.WriteByte(sb, ',')
		util.WriteString(sb, strconv.FormatInt(max, 10))
	}
}

func writeIntRange(min, max int64, inclusive bool, sb *strings.Builder) {
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

func writeFloatRange(min, max float64, inclusive bool, sb *strings.Builder) {
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

func writeTupleArgs(seen []dgo.Value, tt dgo.TupleType, leftSep, rightSep byte, sb *strings.Builder) {
	es := tt.ElementTypes()
	if tt.Variadic() {
		n := es.Len() - 1
		sep := leftSep
		for i := 0; i < n; i++ {
			util.WriteByte(sb, sep)
			sep = ','
			buildTypeString(seen, es.Get(i).(dgo.Type), commaPrio, sb)
		}
		util.WriteByte(sb, sep)
		util.WriteString(sb, `...`)
		vt := es.Get(n).(dgo.ArrayType)
		buildTypeString(seen, vt.ElementType(), commaPrio, sb)
		util.WriteByte(sb, rightSep)
		if !vt.Unbounded() {
			util.WriteByte(sb, '[')
			writeSizeBoundaries(int64(tt.Min()), int64(tt.Max()), sb)
			util.WriteByte(sb, ']')
		}
	} else {
		util.WriteByte(sb, leftSep)
		joinTypes(seen, es, `,`, commaPrio, sb)
		util.WriteByte(sb, rightSep)
	}
}

func writeTernary(seen []dgo.Value, typ dgo.Type, tc func(dgo.Value) dgo.Type, prio int, op string, opPrio int, sb *strings.Builder) {
	if prio >= orPrio {
		util.WriteByte(sb, '(')
	}
	joinX(seen, typ.(dgo.TernaryType).Operands(), tc, op, opPrio, sb)
	if prio >= orPrio {
		util.WriteByte(sb, ')')
	}
}

type typeToString func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder)

var complexTypes map[dgo.TypeIdentifier]typeToString

func init() {
	complexTypes = map[dgo.TypeIdentifier]typeToString{
		dgo.TiAnyOf: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			writeTernary(seen, typ, typeAsType, prio, `|`, orPrio, sb)
		},
		dgo.TiOneOf: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			writeTernary(seen, typ, typeAsType, prio, `^`, xorPrio, sb)
		},
		dgo.TiAllOf: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			writeTernary(seen, typ, typeAsType, prio, `&`, andPrio, sb)
		},
		dgo.TiAllOfValue: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			writeTernary(seen, typ, valueAsType, prio, `&`, andPrio, sb)
		},
		dgo.TiArrayExact: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteByte(sb, '{')
			joinValueTypes(seen, typ.(dgo.ExactType).Value().(dgo.Iterable), `,`, commaPrio, sb)
			util.WriteByte(sb, '}')
		},
		dgo.TiArray: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			at := typ.(dgo.ArrayType)
			if at.Unbounded() {
				util.WriteString(sb, `[]`)
			} else {
				util.WriteByte(sb, '[')
				writeSizeBoundaries(int64(at.Min()), int64(at.Max()), sb)
				util.WriteByte(sb, ']')
			}
			buildTypeString(seen, at.ElementType(), typePrio, sb)
		},
		dgo.TiBinary: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.BinaryType)
			util.WriteString(sb, `binary`)
			if !st.Unbounded() {
				util.WriteByte(sb, '[')
				writeSizeBoundaries(int64(st.Min()), int64(st.Max()), sb)
				util.WriteByte(sb, ']')
			}
		},
		dgo.TiTuple: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			writeTupleArgs(seen, typ.(dgo.TupleType), '{', '}', sb)
		},
		dgo.TiMap: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			at := typ.(dgo.MapType)
			util.WriteString(sb, `map[`)
			buildTypeString(seen, at.KeyType(), commaPrio, sb)
			if !at.Unbounded() {
				util.WriteByte(sb, ',')
				writeSizeBoundaries(int64(at.Min()), int64(at.Max()), sb)
			}
			util.WriteByte(sb, ']')
			buildTypeString(seen, at.ValueType(), typePrio, sb)
		},
		dgo.TiMapExact: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteByte(sb, '{')
			joinValueTypes(seen, typ.(dgo.ExactType).Value().(dgo.Map), `,`, commaPrio, sb)
			util.WriteByte(sb, '}')
		},
		dgo.TiStruct: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteByte(sb, '{')
			st := typ.(dgo.StructMapType)
			joinStructMapEntries(seen, st, sb)
			if st.Additional() {
				if st.Len() > 0 {
					util.WriteByte(sb, ',')
				}
				util.WriteString(sb, `...`)
			}
			util.WriteByte(sb, '}')
		},
		dgo.TiMapEntryExact: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			me := typ.(dgo.ExactType).Value().(dgo.MapEntry)
			buildTypeString(seen, me.Key().Type(), commaPrio, sb)
			util.WriteByte(sb, ':')
			buildTypeString(seen, me.Value().Type(), commaPrio, sb)
		},
		dgo.TiFloatExact: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, util.Ftoa(typ.(dgo.ExactType).Value().(floatVal).GoFloat()))
		},
		dgo.TiFloatRange: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.FloatRangeType)
			writeFloatRange(st.Min(), st.Max(), st.Inclusive(), sb)
		},
		dgo.TiIntegerExact: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, typ.(dgo.ExactType).Value().(fmt.Stringer).String())
		},
		dgo.TiIntegerRange: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.IntegerRangeType)
			writeIntRange(st.Min(), st.Max(), st.Inclusive(), sb)
		},
		dgo.TiRegexpExact: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, typ.TypeIdentifier().String())
			util.WriteByte(sb, '[')
			util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
			util.WriteByte(sb, ']')
		},
		dgo.TiTimeExact: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, typ.TypeIdentifier().String())
			util.WriteByte(sb, '[')
			util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
			util.WriteByte(sb, ']')
		},
		dgo.TiSensitive: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, `sensitive`)
			if op := typ.(dgo.UnaryType).Operand(); DefaultAnyType != op {
				util.WriteByte(sb, '[')
				buildTypeString(seen, op, prio, sb)
				util.WriteByte(sb, ']')
			}
		},
		dgo.TiStringExact: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
		},
		dgo.TiStringPattern: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			RegexpSlashQuote(sb, typ.(dgo.ExactType).Value().(fmt.Stringer).String())
		},
		dgo.TiStringSized: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			st := typ.(dgo.StringType)
			util.WriteString(sb, `string`)
			if !st.Unbounded() {
				util.WriteByte(sb, '[')
				writeSizeBoundaries(int64(st.Min()), int64(st.Max()), sb)
				util.WriteByte(sb, ']')
			}
		},
		dgo.TiCiString: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			util.WriteByte(sb, '~')
			util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).Value().(fmt.Stringer).String()))
		},
		dgo.TiNot: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			nt := typ.(dgo.UnaryType)
			util.WriteByte(sb, '!')
			buildTypeString(seen, nt.Operand(), typePrio, sb)
		},
		dgo.TiNative: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			rt := typ.(dgo.NativeType).GoType()
			util.WriteString(sb, `native`)
			if rt != nil {
				util.WriteByte(sb, '[')
				util.WriteString(sb, strconv.Quote(rt.String()))
				util.WriteByte(sb, ']')
			}
		},
		dgo.TiMeta: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			nt := typ.(dgo.UnaryType)
			util.WriteString(sb, `type`)
			if op := nt.Operand(); DefaultAnyType != op {
				if op == nil {
					util.WriteString(sb, `[type]`)
				} else {
					util.WriteByte(sb, '[')
					buildTypeString(seen, op, prio, sb)
					util.WriteByte(sb, ']')
				}
			}
		},
		dgo.TiFunction: func(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
			ft := typ.(dgo.FunctionType)
			util.WriteString(sb, `func`)
			writeTupleArgs(seen, ft.In(), '(', ')', sb)
			switch ft.Out().Max() {
			case 0:
			case 1:
				buildTypeString(seen, ft.Out().Element(0), prio, sb)
			default:
				util.WriteByte(sb, '(')
				joinTypes(seen, ft.Out().ElementTypes(), `,`, commaPrio, sb)
				util.WriteByte(sb, ')')
			}
		},
	}
}

func buildTypeString(seen []dgo.Value, typ dgo.Type, prio int, sb *strings.Builder) {
	ti := typ.TypeIdentifier()
	if f, ok := complexTypes[ti]; ok {
		if recursionHit(seen, typ) {
			util.WriteString(sb, `<recursive self reference to `)
			util.WriteString(sb, ti.String())
			util.WriteString(sb, ` type>`)
			return
		}
		f(append(seen, typ), typ, prio, sb)
	} else {
		util.WriteString(sb, ti.String())
	}
}
