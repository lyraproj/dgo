package streamer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/lyraproj/dgo/dgo"
)

const (
	firstInArray = iota
	firstInObject
	afterElement
	afterValue
	afterKey
)

// JSON creates a new streamer that will produce JSON when
// receiving values
func JSON(out io.Writer) Consumer {
	return &jsonStreamer{out: out, state: firstInArray, dialect: DgoDialect()}
}

type jsonStreamer struct {
	out     io.Writer
	dialect Dialect
	state   int
}

func (j *jsonStreamer) AddArray(len int, doer dgo.Doer) {
	j.delimit(func() {
		j.state = firstInArray
		assertOk(j.out.Write([]byte{'['}))
		doer()
		assertOk(j.out.Write([]byte{']'}))
	})
}

func (j *jsonStreamer) AddMap(len int, doer dgo.Doer) {
	j.delimit(func() {
		assertOk(j.out.Write([]byte{'{'}))
		j.state = firstInObject
		doer()
		assertOk(j.out.Write([]byte{'}'}))
	})
}

func (j *jsonStreamer) Add(element dgo.Value) {
	j.delimit(func() {
		j.write(element)
	})
}

func (j *jsonStreamer) AddRef(ref int) {
	j.delimit(func() {
		assertOk(fmt.Fprintf(j.out, `{"%s":%d}`, j.dialect.RefKey(), ref))
	})
}

func (j *jsonStreamer) CanDoBinary() bool {
	return false
}

func (j *jsonStreamer) CanDoComplexKeys() bool {
	return false
}

func (j *jsonStreamer) CanDoTime() bool {
	return false
}

func (j *jsonStreamer) StringDedupThreshold() int {
	return 20
}

func (j *jsonStreamer) delimit(doer dgo.Doer) {
	switch j.state {
	case firstInArray:
		doer()
		j.state = afterElement
	case firstInObject:
		doer()
		j.state = afterKey
	case afterKey:
		assertOk(j.out.Write([]byte{':'}))
		doer()
		j.state = afterValue
	case afterValue:
		assertOk(j.out.Write([]byte{','}))
		doer()
		j.state = afterKey
	default: // Element
		assertOk(j.out.Write([]byte{','}))
		doer()
	}
}

func (j *jsonStreamer) write(e dgo.Value) {
	var v []byte
	var err error
	switch e := e.(type) {
	case dgo.String:
		v, err = json.Marshal(e.GoString())
	case dgo.Float:
		v, err = json.Marshal(e.GoFloat())
	case dgo.Integer:
		v, err = json.Marshal(e.GoInt())
	case dgo.Boolean:
		v, err = json.Marshal(e.GoBool())
	default:
		v = []byte(`null`)
	}
	assertOk(0, err)
	assertOk(j.out.Write(v))
}

func assertOk(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
