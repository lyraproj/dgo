package parser

import (
	"fmt"
	"testing"

	"github.com/lyraproj/dgo/util"
)

func Test_nextToken_unicodeError(t *testing.T) {
	sr := util.NewStringReader(string([]byte{0x82, 0xff}))

	defer func() {
		err, ok := recover().(error)
		if !(ok && err.Error() == `unicode error`) {
			t.Error(`expected panic did no occur`)
		}
	}()
	nextToken(sr)
}

func Example_nextToken() {
	const src = `constants: {
    first: 0,
    second: 0x32,
    third: 2e4,
    fourth: 2.3e-2,
    fifth: "hello world",
    value: "String\nWith \\Escape",
    raw: ` + "`" + `String\nWith \\Escape` + "`" + `,
    array: [a, b, c],
    signed: -23,
    positive: +1.2,
    map: {x:, 3},
    group: (x, 3),
    limit: boo<x, 3>
  }`
	sr := util.NewStringReader(src)
	for {
		tf := nextToken(sr)
		if tf.Type == end {
			break
		}
		fmt.Println(tokenString(tf))
	}
	//nolint:gocritic
	// Output:
	//constants
	//':'
	//'{'
	//first
	//':'
	//0
	//','
	//second
	//':'
	//0x32
	//','
	//third
	//':'
	//2e4
	//','
	//fourth
	//':'
	//2.3e-2
	//','
	//fifth
	//':'
	//"hello world"
	//','
	//value
	//':'
	//"String\nWith \\Escape"
	//','
	//raw
	//':'
	//"String\\nWith \\\\Escape"
	//','
	//array
	//':'
	//'['
	//a
	//','
	//b
	//','
	//c
	//']'
	//','
	//signed
	//':'
	//-23
	//','
	//positive
	//':'
	//1.2
	//','
	//map
	//':'
	//'{'
	//x
	//':'
	//','
	//3
	//'}'
	//','
	//group
	//':'
	//'('
	//x
	//','
	//3
	//')'
	//','
	//limit
	//':'
	//boo
	//'<'
	//x
	//','
	//3
	//'>'
	//'}'
}

func TestIsLetterOrDigit(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{name: "Underscore", r: '_', want: false},
		{name: "Zero", r: '0', want: true},
		{name: "Nine", r: '9', want: true},
		{name: "A", r: 'A', want: true},
		{name: "Z", r: 'Z', want: true},
		{name: "a", r: 'a', want: true},
		{name: "z", r: 'z', want: true},
		{name: "Non ASCII Å", r: 'Å', want: false},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLetterOrDigit(tt.r); got != tt.want {
				t.Errorf("IsLetterOrDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLowerCase(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{name: "Underscore", r: '_', want: false},
		{name: "A", r: 'A', want: false},
		{name: "Z", r: 'Z', want: false},
		{name: "a", r: 'a', want: true},
		{name: "z", r: 'z', want: true},
		{name: "Non ASCII Å", r: 'Å', want: false},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLowerCase(tt.r); got != tt.want {
				t.Errorf("IsLowerCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUpperCase(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{name: "Underscore", r: '_', want: false},
		{name: "A", r: 'A', want: true},
		{name: "Z", r: 'Z', want: true},
		{name: "a", r: 'a', want: false},
		{name: "z", r: 'z', want: false},
		{name: "Non ASCII Å", r: 'Å', want: false},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUpperCase(tt.r); got != tt.want {
				t.Errorf("IsUpperCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
