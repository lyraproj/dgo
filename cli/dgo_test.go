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

func TestDgo_validate_bad_port(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	require.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_bad_port.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	s := out.String()
	require.NoMatch(t, `'host' FAILED\!`, s)
	require.Match(t, `'port' FAILED\!(?:.|\s)*2222`, s)
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
