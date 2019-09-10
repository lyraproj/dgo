package internal

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func badNode() *yaml.Node {
	n := &yaml.Node{Kind: yaml.AliasNode}
	n.Alias = n
	return n
}

func okNode() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!!string`, Value: `a`}
}

func unknownTagNode() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!something:here`, Value: `a`}
}

func badSequenceNode() *yaml.Node {
	return &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{badNode()}}
}

func badMappingNodeKey() *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{badNode(), okNode()}}
}

func badMappingNodeValue() *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{okNode(), badNode()}}
}

func Test_decodeBadNode(t *testing.T) {
	_, err := yamlDecodeValue(badNode())
	if err == nil {
		t.Error(`expected error`)
	}
}

func Test_decodeBadSequenceNode(t *testing.T) {
	_, err := yamlDecodeValue(badSequenceNode())
	if err == nil {
		t.Error(`expected error`)
	}
}

func Test_decodeBadMappingNode(t *testing.T) {
	_, err := yamlDecodeValue(badMappingNodeKey())
	if err == nil {
		t.Error(`expected error`)
	}
	_, err = yamlDecodeValue(badMappingNodeValue())
	if err == nil {
		t.Error(`expected error`)
	}
}

func Test_decodeUnknownTagNode(t *testing.T) {
	v, err := yamlDecodeValue(unknownTagNode())
	if err != nil {
		t.Error(err.Error())
	}
	if v.String() != `a` {
		t.Error(`expected string`)
	}
}
