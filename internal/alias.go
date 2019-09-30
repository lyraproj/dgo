package internal

import "github.com/lyraproj/dgo/dgo"

type aliasMap struct {
	typeNames  hashMap
	namedTypes hashMap
}

// NewAliasMap creates a new dgo.Alias map to be used as a scope when parsing types
func NewAliasMap() dgo.AliasMap {
	return &aliasMap{}
}

func (a *aliasMap) GetName(t dgo.Type) dgo.String {
	if v := a.typeNames.Get(t); v != nil {
		return v.(dgo.String)
	}
	return nil
}

func (a *aliasMap) GetType(n dgo.String) dgo.Type {
	if v := a.namedTypes.Get(n); v != nil {
		return v.(dgo.Type)
	}
	return nil
}

func (a *aliasMap) Add(t dgo.Type, name dgo.String) {
	a.typeNames.Put(t, name)
	a.namedTypes.Put(name, t)
}
