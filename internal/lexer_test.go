package internal

import (
	"fmt"
	"testing"

	"github.com/lyraproj/got/util"
)

func Test_unicodeError(t *testing.T) {
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
    map: {x:, 3}
    group: (x, 3)
    limit: boo<x, 3>
  }`
	sr := util.NewStringReader(src)
	for {
		tf := nextToken(sr)
		if tf.i == end {
			break
		}
		fmt.Println(tf)
	}
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
	//map
	//':'
	//'{'
	//x
	//':'
	//','
	//3
	//'}'
	//group
	//':'
	//'('
	//x
	//','
	//3
	//')'
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
