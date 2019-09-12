// Package newtype (Type Factory) contains the factory methods for creating dgo Types
package newtype

import (
	"regexp"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// String returns a new dgo.StringType. It can be called with two optional integer arguments denoting
// the min and max length of the string. If only one integer is given, it represents the min length.
//
// The method can also be called with one string parameter. The returned type will then match that exact
// string and nothing else.
func String(args ...interface{}) dgo.StringType {
	return internal.StringType(args...)
}

// Pattern returns a StringType that is constrained to strings that match the given
// regular expression pattern
func Pattern(pattern *regexp.Regexp) dgo.Type {
	return internal.PatternType(pattern)
}

// Enum returns an AnyOf that represents all of the given strings
func Enum(strings ...string) dgo.Type {
	return internal.EnumType(strings)
}

// IntegerRange returns a dgo.Type that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func IntegerRange(min, max int64, inclusive bool) dgo.IntegerRangeType {
	return internal.IntegerRangeType(min, max, inclusive)
}

// FloatRange returns a dgo.FloatRangeType that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func FloatRange(min, max float64, inclusive bool) dgo.FloatRangeType {
	return internal.FloatRangeType(min, max, inclusive)
}
