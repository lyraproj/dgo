package util

import (
	"strconv"
)

// Ftoa returns the given float as a string with almost all trailing zeroes removed. The resulting string will however
// always contain either the letter 'E' or a dot.
func Ftoa(f float64) string {
	s := strconv.FormatFloat(f, 'G', -1, 64)
	for i := range s {
		switch s[i] {
		case 'e', 'E', '.':
			return s
		}
	}
	return s + `.0`
}
