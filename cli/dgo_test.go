package cli_test

import (
	"strings"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"

	"github.com/lyraproj/dgo/cli"
)

func TestDgo_noArgs(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{}))
	require.Match(t, `dgo: a command `, out.String())
}

func TestDgo_unknownOption(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`what`}))
	require.Match(t, `Error: Unknown option: what`, err.String())
}

func TestDgo_unknownFlag(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`-what`}))
	require.Match(t, `Error: Unknown flag: -what`, err.String())
}

func TestDgo_validate_missingOption(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`validate`}))
	require.Match(t, `Error: Missing required option --input`, err.String())
	require.Equal(t, 1, dgo.Do([]string{`validate`, `--input`, `foo`}))
	require.Match(t, `Error: Missing required option --spec`, err.String())
	require.Equal(t, 1, dgo.Do([]string{`validate`, `--spec`, `foo`}))
	require.Match(t, `Error: Missing required option --input`, err.String())
}

func TestDgo_validate_missingOptionArgument(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`validate`, `--input`}))
	require.Match(t, `Error: Option --input requires an argument`, err.String())
	require.Equal(t, 1, dgo.Do([]string{`validate`, `--input`, `foo`, `--spec`}))
	require.Match(t, `Error: Option --spec requires an argument`, err.String())
}

func TestDgo_validate_unknownOption(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`validate`, `--output`}))
	require.Match(t, `Error: Unknown flag: --output`, err.String())
}

func TestDgo_help(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 0, dgo.Do([]string{`help`}))
	require.Match(t, `dgo: a command `, out.String())
	require.Equal(t, 0, dgo.Do([]string{`--help`}))
	require.Match(t, `dgo: a command `, out.String())
	require.Equal(t, 0, dgo.Do([]string{`-?`}))
	require.Match(t, `dgo: a command `, out.String())
}

func TestDgo_validate_noArgs(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`validate`}))
	require.Match(t, `^Error: Missing required option --input`, err.String())
}

func TestDgo_validate_help(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 0, dgo.Do([]string{`validate`, `help`}))
	require.Match(t, `dgo validate \[flags\]`, out.String())
}

func TestDgo_validate_ok(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 0, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	s := out.String()
	require.Match(t, `'host' OK\!`, s)
	require.Match(t, `'port' OK\!`, s)
}

func TestDgo_validate_ok_brief(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 0, dgo.Do([]string{`validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	require.Equal(t, ``, out.String())
}

func TestDgo_validate_bad_port(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`validate`, `--input`, `testdata/service_bad_port.yaml`, `--spec`, `testdata/servicespec.yaml`, `--verbose`}))
	s := out.String()
	require.NoMatch(t, `'host' FAILED\!`, s)
	require.Match(t, `'port' FAILED\!(?:.|\s)*2222`, s)
}

func TestDgo_validate_bad_port_brief(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`validate`, `--input`, `testdata/service_bad_port.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	require.Equal(t, "parameter 'port' is not an instance of type 1..999\n", out.String())
}

func TestDgo_validate_bad_port_dgo(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_bad_port.yaml`, `--spec`, `testdata/servicespec.dgo`}))
	s := out.String()
	require.NoMatch(t, `'host' FAILED\!`, s)
	require.Match(t, `'port' FAILED\!(?:.|\s)*2222`, s)
}

func TestDgo_validate_extraneous_param(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_extraneous_param.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	s := out.String()
	require.Match(t, `'host' OK\!`, s)
	require.Match(t, `'port' OK\!`, s)
	require.Match(t, `'login' FAILED\!(?:.|\s)*key is not found in definition`, s)
}

func TestDgo_validate_missing_host(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_missing_host.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	s := out.String()
	require.Match(t, `'port' OK\!`, s)
	require.Match(t, `'host' FAILED\!(?:.|\s)*required key not found in input`, s)
}

func TestDgo_validate_no_input_file(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/server.yaml`, `--spec`, `testdata/servicespec.dgo`}))
	s := err.String()
	require.Match(t, `server\.yaml.*no such file or directory`, s)
}

func TestDgo_validate_no_spec_file(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/serverspec.dgo`}))
	s := err.String()
	require.Match(t, `serverspec\.dgo.*no such file or directory`, s)
}

func TestDgo_validate_input_not_map(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_array.yaml`, `--spec`, `testdata/servicespec.dgo`}))
	s := err.String()
	require.Match(t, `Error: expecting data to be a map`, s)
}

func TestDgo_validate_spec_not_map(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec_array.yaml`}))
	s := err.String()
	require.Match(t, `Error: expecting data to be a map`, s)
}

func TestDgo_validate_spec_bad_type(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec_bad_type.dgo`}))
	s := err.String()
	require.Match(t, `does not contain a struct definition`, s)
}

func TestDgo_validate_input_extension(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.pson`, `--spec`, `testdata/servicespec.dgo`}))
	s := err.String()
	require.Match(t, `expected file name to end with \.yaml or \.json`, s)
}

func TestDgo_validate_spec_extension(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec.go`}))
	s := err.String()
	require.Match(t, `expected file name to end with \.yaml, \.json, or \.dgo`, s)
}

func TestDgo_validate_bad_dgo(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec_bad.dgo`}))
	s := err.String()
	require.Match(t, `Error: mix of elements and map entries`, s)
}
