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
	return internal.StringType(args)
}

// Pattern returns a StringType that is constrained to strings that match the given
// regular expression pattern
func Pattern(pattern *regexp.Regexp) dgo.Type {
	return internal.PatternType(pattern)
}

// CiString returns a StringType that is constrained to strings that are equal to the given string under
// Unicode case-folding.
func CiString(s interface{}) dgo.StringType {
	return internal.CiStringType(s)
}

// Enum returns a Type that represents all of the given strings.
func Enum(strings ...string) dgo.Type {
	return internal.EnumType(strings)
}

// CiEnum returns a Type that represents all strings that are equal to one of the given strings
// under Unicode case-folding.
func CiEnum(strings ...string) dgo.Type {
	return internal.CiEnumType(strings)
}

// IntegerRange returns a dgo.Type that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func IntegerRange(min, max int64, inclusive bool) dgo.IntegerRangeType {
	return internal.IntegerRangeType(min, max, inclusive)
}

// IntEnum returns a Type that represents any of the given integers
func IntEnum(ints ...int) dgo.Type {
	return internal.IntEnumType(ints)
}

// FloatRange returns a dgo.FloatRangeType that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func FloatRange(min, max float64, inclusive bool) dgo.FloatRangeType {
	return internal.FloatRangeType(min, max, inclusive)
}
