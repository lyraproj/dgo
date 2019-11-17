package pcore

import (
	"fmt"

	"github.com/lyraproj/dgo/util"
)

func Example_nextToken() {
	const src = `# This is scanned code.
  constants => {
    first => 0,
    second => 0x32,
    third => 2e4,
    fourth => 2.3e-2,
    fifth => 'hello',
    sixth => "world",
    type => Foo::Bar,
    value => "String\nWith \\Escape",
    array => [a, b, c],
    call => Boo::Bar(x, 3)
  }`
	sr := util.NewStringReader(src)
	for {
		tf := nextToken(sr)
		if tf.Type == end {
			break
		}
		fmt.Println(tokenString(tf))
	}
	// Output:
	//identifier: 'constants'
	//rocket: '=>'
	//{: ''
	//identifier: 'first'
	//rocket: '=>'
	//integer: '0'
	//,: ''
	//identifier: 'second'
	//rocket: '=>'
	//integer: '0x32'
	//,: ''
	//identifier: 'third'
	//rocket: '=>'
	//float: '2e4'
	//,: ''
	//identifier: 'fourth'
	//rocket: '=>'
	//float: '2.3e-2'
	//,: ''
	//identifier: 'fifth'
	//rocket: '=>'
	//string: 'hello'
	//,: ''
	//identifier: 'sixth'
	//rocket: '=>'
	//string: 'world'
	//,: ''
	//identifier: 'type'
	//rocket: '=>'
	//name: 'Foo::Bar'
	//,: ''
	//identifier: 'value'
	//rocket: '=>'
	//string: 'String
	//With \Escape'
	//,: ''
	//identifier: 'array'
	//rocket: '=>'
	//[: ''
	//identifier: 'a'
	//,: ''
	//identifier: 'b'
	//,: ''
	//identifier: 'c'
	//]: ''
	//,: ''
	//identifier: 'call'
	//rocket: '=>'
	//name: 'Boo::Bar'
	//(: ''
	//identifier: 'x'
	//,: ''
	//integer: '3'
	//): ''
	//}: ''
}
