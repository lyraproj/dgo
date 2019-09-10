package internal_test

import (
	"fmt"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"

	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/vf"
)

func ExampleToIndentedStringERP() {
	v := vf.Map(
		`name`, `Bob`,
		`address`, vf.Map(
			`street`, `Foil street`,
			`zip`, `54321`,
			`city`, `Smallfylke`,
			`gender`, `male`),
		`age`, 32)
	fmt.Println(internal.ToIndentedStringERP(v))
	// Output:
	// {
	//   "name": "Bob",
	//   "address": {
	//     "street": "Foil street",
	//     "zip": "54321",
	//     "city": "Smallfylke",
	//     "gender": "male"
	//   },
	//   "age": 32
	// }
}

func TestToIndentedStringERP_nonStringer(t *testing.T) {
	ei := internal.NewERPIndenter(` `)
	ei.AppendValue(struct{ A string }{`hello`})
	require.Equal(t, `struct { A string }{A:"hello"}`, ei.String())
}
