package examples

import (
	"errors"
	"fmt"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/newtype"
)

// convertTypeStringsToTypes finds all occurrences of key '$type' recursively and converts the
// corresponding value into a dgo.Type
func convertTypeStringsToTypes(v dgo.Value) dgo.Value {
	switch cv := v.(type) {
	case dgo.Map:
		v = cv.Map(func(e dgo.MapEntry) interface{} {
			if e.Key().Equals(`$type`) {
				return newtype.Parse(e.Value().String())
			}
			return convertTypeStringsToTypes(e.Value())
		})
	case dgo.Array:
		v = cv.Map(func(e dgo.Value) interface{} {
			return convertTypeStringsToTypes(e)
		})
	}
	return v
}

// validateParameterValues validates the given parameters using the given description. Errors are generated
// if a required parameter is missing, not recognized, or if it is of incorrect type
func validateParameterValues(description dgo.Map, parameters dgo.Value) (errs []error) {
	pm, ok := parameters.(dgo.Map)
	if !ok {
		return []error{errors.New(`parameters is not a Map`)}
	}
	description.Each(func(e dgo.MapEntry) {
		pd := e.Value().(dgo.Map)
		if v := pm.Get(e.Key()); v != nil {
			t := pd.Get(`$type`)
			if !t.(dgo.Type).Instance(v) {
				errs = append(errs, fmt.Errorf(`parameter '%s' is not an instance of type %s`, e.Key(), t))
			}
		} else if rq := pd.Get(`required`); rq != nil && rq.(dgo.Boolean).GoBool() {
			errs = append(errs, fmt.Errorf(`missing required parameter '%s'`, e.Key()))
		}
	})
	pm.EachKey(func(k dgo.Value) {
		if description.Get(k) == nil {
			errs = append(errs, fmt.Errorf(`unknown parameter '%s'`, k))
		}
	})
	return
}
