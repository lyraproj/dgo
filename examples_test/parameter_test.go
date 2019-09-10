package examples

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
	"gopkg.in/yaml.v3"
)

// This is the type that describes how a parameters map is constituted
const parametersType = `map[string]{name: string[1], $type: dgo, required?: bool}`

// Sample parameter map
const sampleParameters = `
host:
  $type: string[1]
  name: sample/service_host
  required: true
port:
  $type: 1..999
  name: sample/service_port
`

func TestValidateParameterValues(t *testing.T) {
	const sampleValues = `
host: example.com
port: 22
`
	params, err := vf.UnmarshalYAML([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectNoErrors(t, validateParameterValues(loadParamsDesc(t), params))
}

func TestValidateParameterValues_failRequired(t *testing.T) {
	const sampleValues = `
port: 22
`
	params, err := vf.UnmarshalYAML([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectError(t, `missing required parameter 'host'`, validateParameterValues(loadParamsDesc(t), params))
}

func TestValidateParameterValues_failNotRecognized(t *testing.T) {
	const sampleValues = `
host: example.com
port: 22
login: foo:bar
`
	params, err := vf.UnmarshalYAML([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectError(t, `unknown parameter 'login'`, validateParameterValues(loadParamsDesc(t), params))
}

func TestValidateParameterValues_failInvalidHostType(t *testing.T) {
	const sampleValues = `
host: 85493
port: 22
`
	params, err := vf.UnmarshalYAML([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectError(t, `parameter 'host' is not an instance of type string[1]`, validateParameterValues(loadParamsDesc(t), params))
}

func TestValidateParameterValues_failInvalidPortType(t *testing.T) {
	const sampleValues = `
host: example.com
port: 1022
`
	params, err := vf.UnmarshalYAML([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectError(t, `parameter 'port' is not an instance of type 1..999`, validateParameterValues(loadParamsDesc(t), params))
}

func loadParamsDesc(t *testing.T) dgo.Map {
	t.Helper()
	paramsDesc := vf.MutableMap(parametersType)
	if err := yaml.Unmarshal([]byte(sampleParameters), paramsDesc); err != nil {
		t.Fatal(err)
	}
	return convertTypeStringsToTypes(paramsDesc).(dgo.Map)
}

func expectError(t *testing.T, error string, errors []error) {
	t.Helper()
	if len(errors) != 1 || error != errors[0].Error() {
		t.Errorf(`expected "%s" error`, error)
		expectNoErrors(t, errors)
	}
}

func expectNoErrors(t *testing.T, errors []error) {
	t.Helper()
	for _, err := range errors {
		t.Error(err)
	}
}
