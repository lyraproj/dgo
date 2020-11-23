package streamer_test

import (
	"testing"

	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestDataCollector(t *testing.T) {
	a := vf.Values(
		typ.String,
		vf.Binary([]byte{1, 2, 3}, false),
		vf.Sensitive(`secret`),
		vf.DurationFromString(`2h48m3.3s`),
		vf.TimeFromString(`2019-10-06T16:15:00.123-01:00`))

	c := streamer.DataCollector()
	streamer.New(nil, nil).Stream(a, c)

	assert.Equal(t, vf.Values(
		vf.Map(`__type`, `string`),
		vf.Map(`__type`, `binary`, `__value`, `AQID`),
		vf.Map(`__type`, `sensitive`, `__value`, `secret`),
		vf.Map(`__type`, `duration`, `__value`, `2h48m3.3s`),
		vf.Map(`__type`, `time`, `__value`, `2019-10-06T16:15:00.123-01:00`),
	), c.Value())
}
