package pcore

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/vf"
)

// Dialect returns the pcore dialect
func Dialect() streamer.Dialect {
	return pcoreDialectSingleton
}

type pcoreDialect int

const pcoreDialectSingleton = pcoreDialect(0)

var typeKey = vf.String(`__ptype`)
var valueKey = vf.String(`__pvalue`)
var refKey = vf.String(`__pref`)
var aliasType = vf.String(`Alias`)
var binaryTyp = vf.String(`Binary`)
var sensitiveTyp = vf.String(`Sensitive`)
var mapType = vf.String(`Hash`)
var timeType = vf.String(`Timestamp`)
var bigFloatType = vf.String("Float")
var bigIntegerType = vf.String("Integer")

func (d pcoreDialect) TypeKey() dgo.String {
	return typeKey
}

func (d pcoreDialect) ValueKey() dgo.String {
	return valueKey
}

func (d pcoreDialect) RefKey() dgo.String {
	return refKey
}

func (d pcoreDialect) AliasTypeName() dgo.String {
	return aliasType
}

func (d pcoreDialect) BigFloatTypeName() dgo.String {
	return bigFloatType
}

func (d pcoreDialect) BigIntTypeName() dgo.String {
	return bigIntegerType
}

func (d pcoreDialect) BinaryTypeName() dgo.String {
	return binaryTyp
}

func (d pcoreDialect) MapTypeName() dgo.String {
	return mapType
}

func (d pcoreDialect) SensitiveTypeName() dgo.String {
	return sensitiveTyp
}

func (d pcoreDialect) TimeTypeName() dgo.String {
	return timeType
}

func (d pcoreDialect) ParseType(aliasMap dgo.AliasAdder, typeString dgo.String) (dt dgo.Type) {
	return Parse(typeString.GoString()).(dgo.Type)
}
