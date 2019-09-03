// Package newtype (Type Factory) contains all factory methods for creating types
package newtype

import (
	"regexp"

	"github.com/lyraproj/got/dgo"
	"github.com/lyraproj/got/internal"
)

// String returns a new got.StringType. It can be called with two optional integer arguments denoting
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

// IntegerRange returns a got.Type that is limited to the inclusive range given by min and max
func IntegerRange(min, max int64) dgo.IntegerRangeType {
	return internal.IntegerRangeType(min, max)
}

// FloatRange returns a got.Type that is limited to the inclusive range given by min and max
func FloatRange(min, max float64) dgo.FloatRangeType {
	return internal.FloatRangeType(min, max)
}
