package internal

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/lyraproj/dgo/dgo"
)

func UnmarshalJSON(b []byte) (dgo.Value, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return jsonDecodeValue(dec)
}

type decodeEndToken rune

const endOfArray = decodeEndToken(']')
const endOfMap = decodeEndToken('}')

func (e decodeEndToken) Error() string {
	return fmt.Sprintf(`unexpected delimiter '%c'`, rune(e))
}

func jsonDecodeValue(dec *json.Decoder) (e dgo.Value, err error) {
	var t json.Token
	if t, err = dec.Token(); err != nil {
		return
	}
	switch t := t.(type) {
	case json.Delim:
		switch t {
		case '[':
			return jsonDecodeArray(dec)
		case '{':
			return jsonDecodeMap(dec)
		default:
			return nil, decodeEndToken(t)
		}
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return Integer(i), nil
		} else {
			f, _ := t.Float64()
			return Float(f), nil
		}
	default:
		return Value(t), nil
	}
}

func jsonDecodeArray(dec *json.Decoder) (*array, error) {
	s := make([]dgo.Value, 0)
	for {
		e, err := jsonDecodeValue(dec)
		if err != nil {
			if err == endOfArray {
				return &array{slice: s, frozen: true}, nil
			}
			return nil, err
		}
		s = append(s, e)
	}
}

func jsonDecodeMap(dec *json.Decoder) (*hashMap, error) {
	m := &hashMap{}
	for {
		k, err := jsonDecodeValue(dec)
		if err != nil {
			if err == endOfMap {
				break
			}
			return nil, err
		}
		v, err := jsonDecodeValue(dec)
		if err != nil {
			return nil, err
		}
		nd := &hashNode{key: k, value: v, prev: m.last}
		if m.first == nil {
			m.first = nd
		} else {
			m.last.next = nd
		}
		m.last = nd
		m.len++
	}

	tbl := make([]*hashNode, tableSizeFor(m.len))
	hl := len(tbl) - 1
	for e := m.first; e != nil; e = e.next {
		hk := hl & hash(e.key.HashCode())
		e.hashNext = tbl[hk]
		tbl[hk] = e
	}
	m.table = tbl
	return m, nil
}
