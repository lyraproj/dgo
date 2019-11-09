package internal

import (
	"github.com/lyraproj/dgo/dgo"
)

type arguments struct {
	array
}

// Arguments returns an immutable Arguments instance that represents the given slice
func Arguments(values []interface{}) dgo.Arguments {
	a := Values(values).(*array)
	return &arguments{*a}
}

func (a *arguments) AssertSize(funcName string, min, max int) {
	l := a.Len()
	if min > l || l > max {
		panic(illegalArgumentCount(funcName, min, max, l))
	}
}

func (a *arguments) Arg(funcName string, n int, typ dgo.Type) dgo.Value {
	v := a.Get(n)
	if typ.Instance(v) {
		return v
	}
	panic(illegalArgument(funcName, typ.String(), sliceToInterfaces(&a.array), n))
}
