package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// sizedStringType is a size constrained String type. It only represents all strings whose
	// length is within the inclusive min, max range
	sizedStringType struct {
		min int
		max int
	}

	// defaultStringType represents an string without constraints
	defaultStringType int

	// exactStringType only represents its own string
	exactStringType hstring

	// patternType constrains its instances to those that match the regexp pattern
	patternType struct {
		*regexp.Regexp
	}

	// hstring is a string that caches the hash value when it is computed
	hstring struct {
		s string
		h int
	}
)

// DefaultStringType is the unconstrained String type
const DefaultStringType = defaultStringType(0)

// EnumType returns a StringType that represents all of the given strings
func EnumType(strings []string) dgo.Type {
	switch len(strings) {
	case 0:
		return &notType{DefaultAnyType}
	case 1:
		return (*exactStringType)(makeHString(strings[0]))
	}
	ts := make([]dgo.Value, len(strings))
	for i := range strings {
		et := (*exactStringType)(makeHString(strings[i]))
		ts[i] = et
	}
	return &anyOfType{slice: ts, frozen: true}
}

// String returns the dgo.String for the given string
func String(s string) dgo.String {
	return makeHString(s)
}

func (t defaultStringType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case defaultStringType, *exactStringType, *sizedStringType, *patternType:
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t defaultStringType) Equals(other interface{}) bool {
	return t == other
}

func (t defaultStringType) HashCode() int {
	return int(dgo.TiString)
}

func (t defaultStringType) Instance(value interface{}) bool {
	if _, ok := value.(*hstring); ok {
		return true
	}
	if _, ok := value.(string); ok {
		return true
	}
	return false
}

func (t defaultStringType) Max() int {
	return math.MaxInt64
}

func (t defaultStringType) Min() int {
	return 0
}

func (t defaultStringType) String() string {
	return TypeString(t)
}

func (t defaultStringType) Type() dgo.Type {
	return &metaType{t}
}

func (t defaultStringType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiString
}

func (t defaultStringType) Unbounded() bool {
	return true
}

func (t *exactStringType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*exactStringType); ok {
		return t.s == ot.s
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *exactStringType) Equals(v interface{}) bool {
	if ov, ok := v.(*exactStringType); ok {
		return t.s == ov.s
	}
	return false
}

func (t *exactStringType) HashCode() int {
	return (*hstring)(t).HashCode() * 5
}

func (t *exactStringType) Instance(v interface{}) bool {
	if ov, ok := v.(*hstring); ok {
		return t.s == ov.s
	}
	if ov, ok := v.(string); ok {
		return t.s == ov
	}
	return false
}

func (t *exactStringType) Max() int {
	return len(t.s)
}

func (t *exactStringType) Min() int {
	return len(t.s)
}

func (t *exactStringType) String() string {
	return TypeString(t)
}

func (t *exactStringType) Type() dgo.Type {
	return &metaType{t}
}

func (t *exactStringType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiStringExact
}

func (t *exactStringType) Unbounded() bool {
	return false
}

func (t *exactStringType) Value() dgo.Value {
	return makeHString(t.s)
}

// PatternType returns a StringType that is constrained to strings that match the given
// regular expression pattern
func PatternType(pattern *regexp.Regexp) dgo.Type {
	return &patternType{Regexp: pattern}
}

// StringType returns a new dgo.StringType. It can be called with two optional integer arguments denoting
// the min and max length of the string. If only one integer is given, it represents the min length.
//
// The method can also be called with one string parameter. The returned type will then match that exact
// string and nothing else.
func StringType(args ...interface{}) dgo.StringType {
	switch len(args) {
	case 0:
		return DefaultStringType
	case 1:
		switch a0 := Value(args[0]).(type) {
		case dgo.String:
			return a0.Type().(dgo.StringType) // Exact string type
		case dgo.Integer:
			return SizedStringType(int(a0.GoInt()), math.MaxInt64)
		}
		panic(illegalArgument(`StringType`, `intVal or String`, args, 0))
	case 2:
		if a0, ok := Value(args[0]).(dgo.Integer); ok {
			var a1 dgo.Integer
			if a1, ok = Value(args[1]).(dgo.Integer); ok {
				return SizedStringType(int(a0.GoInt()), int(a1.GoInt()))
			}
			panic(illegalArgument(`StringType`, `intVal`, args, 1))
		}
		panic(illegalArgument(`StringType`, `intVal`, args, 0))
	}
	panic(fmt.Errorf(`illegal number of arguments for StringType. Expected 0 - 2, got %d`, len(args)))
}

func (t *patternType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case *exactStringType:
		return t.IsInstance(ot.s)
	case *patternType:
		return t.String() == ot.String()
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *patternType) Equals(v interface{}) bool {
	if ov, ok := v.(*patternType); ok {
		return t.String() == ov.String()
	}
	return false
}

func (t *patternType) HashCode() int {
	return stringHash(t.String())
}

func (t *patternType) Instance(v interface{}) bool {
	if sv, ok := v.(*hstring); ok {
		return t.MatchString(sv.s)
	}
	if sv, ok := v.(string); ok {
		return t.MatchString(sv)
	}
	return false
}

func (t *patternType) IsInstance(v string) bool {
	return t.MatchString(v)
}

func (t *patternType) Type() dgo.Type {
	return &metaType{t}
}

func (t *patternType) String() string {
	return TypeString(t)
}

func (t *patternType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiStringPattern
}

func (t *patternType) Value() dgo.Value {
	return (*regexpVal)(t.Regexp)
}

// SizedStringType returns a StringType that is constrained to strings whose length is within the
// inclusive range given by min and max.
func SizedStringType(min, max int) dgo.StringType {
	if min < 0 {
		min = 0
	}
	if max < min {
		tmp := max
		max = min
		min = tmp
	}
	if min == 0 && max == math.MaxInt64 {
		return DefaultStringType
	}
	return &sizedStringType{min: min, max: max}
}

func (t *sizedStringType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case *exactStringType:
		return t.IsInstance(ot.s)
	case *sizedStringType:
		return t.min <= ot.min && t.max >= ot.max
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *sizedStringType) Equals(v interface{}) bool {
	if ob, ok := v.(*sizedStringType); ok {
		return *t == *ob
	}
	return false
}

func (t *sizedStringType) HashCode() int {
	h := int(dgo.TiStringSized)
	if t.min > 0 {
		h = h*31 + t.min
	}
	if t.max < math.MaxInt64 {
		h = h*31 + t.max
	}
	return h
}

func (t *sizedStringType) Instance(v interface{}) bool {
	if sv, ok := v.(*hstring); ok {
		return t.IsInstance(sv.s)
	}
	if sv, ok := v.(string); ok {
		return t.IsInstance(sv)
	}
	return false
}

func (t *sizedStringType) IsInstance(v string) bool {
	l := len(v)
	return t.min <= l && l <= t.max
}

func (t *sizedStringType) Max() int {
	return t.max
}

func (t *sizedStringType) Min() int {
	return t.min
}

func (t *sizedStringType) String() string {
	return TypeString(t)
}

func (t *sizedStringType) Type() dgo.Type {
	return &metaType{t}
}

func (t *sizedStringType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiStringSized
}

func (t *sizedStringType) Unbounded() bool {
	return t.min == 0 && t.max == math.MaxInt64
}

func stringHash(s string) int {
	h := 1
	for i := range s {
		h = 31*h + int(s[i])
	}
	return h
}

func makeHString(s string) *hstring {
	return &hstring{s: s}
}

func (v *hstring) AppendTo(w *util.Indenter) {
	w.Append(strconv.Quote(v.s))
}

func (v *hstring) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	switch ov := other.(type) {
	case *hstring:
		switch {
		case v.s > ov.s:
			r = 1
		case v.s < ov.s:
			r = -1
		default:
			r = 0
		}
	case string:
		switch {
		case v.s > ov:
			r = 1
		case v.s < ov:
			r = -1
		default:
			r = 0
		}
	case nilValue:
		r = 1
	default:
		ok = false
	}
	return
}

func (v *hstring) Equals(other interface{}) bool {
	// comparison for *hstring must be first here or the HashMap will get a penalty. It
	// must always use *hstring to get the hash code
	if ov, ok := other.(*hstring); ok {
		return v.s == ov.s
	}
	if s, ok := other.(string); ok {
		return v.s == s
	}
	return false
}

func (v *hstring) GoString() string {
	return v.s
}

func (v *hstring) HashCode() int {
	if v.h == 0 {
		v.h = stringHash(v.s)
	}
	return v.h
}

func (v *hstring) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.s)
}

func (v *hstring) MarshalYAML() (interface{}, error) {
	n := &yaml.Node{}
	n.SetString(v.s)
	return n, nil
}

func (v *hstring) String() string {
	return v.s
}

func (v *hstring) Type() dgo.Type {
	return (*exactStringType)(v)
}

func (v *hstring) UnmarshalJSON(b []byte) error {
	v.h = 0
	return json.Unmarshal(b, &v.s)
}

func (v *hstring) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode && n.Tag == `!!str` {
		v.h = 0
		v.s = n.Value
		return nil
	}
	return errors.New(`"expecting data to be a string"`)
}
