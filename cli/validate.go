package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/lyraproj/dgo/newtype"

	"github.com/lyraproj/dgo/vf"
	"gopkg.in/yaml.v3"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/util"
)

// Help prints the validate command help on stdout
func (h *validateCommand) Help() int {
	util.WriteString(h.out, `validate input against the dgo type system

Usage:
  dgo validate [flags]

Aliases:
  v, val, valid

Flags:
  -i, --input  Relative or absolute path to a yaml file containing input to validate
  -s, --spec   Relative or absolute path to a yaml or dgo file with the parameter definitions

Global Flags:
  -v, --verbose   Be verbose in output
`)
	return 0
}

// Validate is the Dgo sub command that reads and validates YAML input using a spec written in
// YAML or Dgo
func Validate(out, err io.Writer, debug bool) Command {
	return &validateCommand{command{`dgo validate`, out, err, debug}}
}

type validateCommand struct {
	command
}

func (h *validateCommand) run(input, spec string) int {
	// read function that succeed or panics
	read := func(name string) []byte {
		/* #nosec */
		bs, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}
		return bs
	}

	var iMap dgo.Map
	switch {
	case strings.HasSuffix(input, `.yaml`), strings.HasSuffix(input, `.json`):
		m := vf.MutableMap(nil)
		data := read(input)
		if err := yaml.Unmarshal(data, m); err != nil {
			panic(err)
		}
		iMap = m
		if h.verbose {
			bld := util.NewIndenter(`  `)
			bld.Append(`Got input yaml with:`)
			b2 := bld.Indent()
			b2.NewLine()
			b2.AppendIndented(string(data))
			util.WriteString(h.out, bld.String())
		}
	default:
		panic(fmt.Errorf(`invalid file name '%s', expected file name to end with .yaml or .json`, input))
	}

	var sType dgo.StructType
	switch {
	case strings.HasSuffix(spec, `.yaml`), strings.HasSuffix(spec, `.json`):
		m := vf.MutableMap(nil)
		if err := yaml.Unmarshal(read(spec), m); err != nil {
			panic(err)
		}
		sType = newtype.StructFromMap(false, m)
	case strings.HasSuffix(spec, `.dgo`):
		tp := newtype.ParseFile(spec, string(read(spec)))
		if st, ok := tp.(dgo.StructType); ok {
			sType = st
		} else {
			panic(fmt.Errorf(`file '%s' does not contain a struct definition`, spec))
		}
	default:
		panic(fmt.Errorf(`invalid file name '%s', expected file name to end with .yaml, .json, or .dgo`, input))
	}

	ok := true
	if h.verbose {
		bld := util.NewIndenter(`  `)
		ok = sType.ValidateVerbose(iMap, bld)
		util.WriteString(h.out, bld.String())
	} else {
		vs := sType.Validate(nil, iMap)
		if len(vs) > 0 {
			ok = false
			for _, err := range vs {
				util.WriteString(h.out, err.Error())
				util.WriteRune(h.out, '\n')
			}
		}
	}
	if ok {
		return 0
	}
	return 1
}

// Do parses the validate command line options and runs the validation
func (h *validateCommand) Do(args []string) int {
	if len(args) == 0 {
		return h.MissingOption(`--input`)
	}

	input := ``
	spec := ``
	for len(args) > 0 {
		opt := args[0]
		args = args[1:]
		switch opt {
		case `help`, `-?`, `--help`, `-help`:
			return h.Help()
		case `-i`, `--input`:
			if len(args) == 0 {
				return h.MissingOptionArgument(opt)
			}
			input = args[0]
			args = args[1:]
		case `-s`, `--spec`:
			if len(args) == 0 {
				return h.MissingOptionArgument(opt)
			}
			spec = args[0]
			args = args[1:]
		case `-v`, `--verbose`:
			h.verbose = true
		default:
			return h.UnknownOption(opt)
		}
	}
	if input == `` {
		return h.MissingOption(`--input`)
	}
	if spec == `` {
		return h.MissingOption(`--spec`)
	}
	return h.RunWithCatch(func() int {
		return h.run(input, spec)
	})
}
