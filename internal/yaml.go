package internal

import (
	"github.com/lyraproj/got/dgo"
	"gopkg.in/yaml.v3"
)

func UnmarshalYAML(b []byte) (dgo.Value, error) {
	var n yaml.Node
	err := yaml.Unmarshal(b, &n)
	if err != nil {
		return nil, err
	}
	return yamlDecodeValue(&n)
}

func yamlEncodeValue(v dgo.Value) (*yaml.Node, error) {
	switch v := v.(type) {
	case yaml.Marshaler:
		yv, err := v.MarshalYAML()
		if err != nil {
			return nil, err
		}
		if n, ok := yv.(*yaml.Node); ok {
			return n, nil
		}
	case dgo.Type:
		return yamlMarshalType(v)
	case dgo.Native:
		if ym, ok := v.GoValue().(yaml.Marshaler); ok {
			yv, err := ym.MarshalYAML()
			if err != nil {
				return nil, err
			}
			if n, ok := yv.(*yaml.Node); ok {
				return n, nil
			}
			return yamlEncodeValue(value(yv))
		}
	}
	n := &yaml.Node{}
	n.SetString(v.String())
	return n, nil
}

func yamlDecodeScalar(n *yaml.Node) (dgo.Value, error) {
	var v dgo.Value
	switch n.Tag {
	case `!!null`:
		v = Nil
	case `!!bool`:
		var x bool
		if err := n.Decode(&x); err != nil {
			return nil, err
		}
		v = Boolean(x)
	case `!!int`:
		var x int64
		if err := n.Decode(&x); err != nil {
			return nil, err
		}
		v = Integer(x)
	case `!!float`:
		var x float64
		if err := n.Decode(&x); err != nil {
			return nil, err
		}
		v = Float(x)
	case `!!str`:
		v = makeHString(n.Value)
	/* TODO: timestamp and binary
	case `!!timestamp`:
		var x time.Time
		if err := n.Decode(&x); err != nil {
			panic(err)
		}
		v = Timestamp(x)
	*/
	case `!!binary`:
		v = BinaryFromString(n.Value)
	default:
		var x interface{}
		if err := n.Decode(&x); err != nil {
			return nil, err
		}
		v = Value(x)
	}
	return v, nil
}

func yamlDecodeValue(n *yaml.Node) (v dgo.Value, err error) {
	switch n.Kind {
	case yaml.DocumentNode:
		v, err = yamlDecodeValue(n.Content[0])
	case yaml.SequenceNode:
		v, err = yamlDecodeArray(n)
	case yaml.MappingNode:
		v, err = yamlDecodeMap(n)
	default:
		v, err = yamlDecodeScalar(n)
	}
	return
}

func yamlDecodeArray(n *yaml.Node) (a *array, err error) {
	ms := n.Content
	es := make([]dgo.Value, len(ms))
	for i, me := range ms {
		es[i], err = yamlDecodeValue(me)
		if err != nil {
			return
		}
	}
	a = &array{slice: es, frozen: true}
	return
}

func yamlDecodeMap(n *yaml.Node) (*hashMap, error) {
	ms := n.Content
	top := len(ms)
	tbl := make([]*hashNode, tableSizeFor(top/2))
	hl := len(tbl) - 1
	m := &hashMap{table: tbl, len: top / 2, frozen: true}

	for i := 0; i < top; {
		k, err := yamlDecodeValue(ms[i])
		if err != nil {
			return nil, err
		}
		i++
		v, err := yamlDecodeValue(ms[i])
		if err != nil {
			return nil, err
		}
		i++
		hk := hl & hash(k.HashCode())
		nd := &hashNode{key: k, value: v, hashNext: tbl[hk], prev: m.last}
		if m.first == nil {
			m.first = nd
		} else {
			m.last.next = nd
		}
		m.last = nd
		tbl[hk] = nd
	}
	return m, nil
}

func yamlMarshalType(t dgo.Type) (*yaml.Node, error) {
	return &yaml.Node{Tag: `!puppet.com,2019:got/type`, Kind: yaml.ScalarNode, Value: TypeString(t)}, nil
}
