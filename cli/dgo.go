package cli

import (
	"io"

	"github.com/lyraproj/dgo/util"
)

// Dgo creates the global dgo command
func Dgo(out, err io.Writer) Command {
	return &dgoCommand{command{`dgo`, out, err, false}}
}

type dgoCommand struct {
	command
}

func (h *dgoCommand) Do(args []string) int {
	verbose := false
	for len(args) > 0 {
		cmd := args[0]
		args = args[1:]
		switch cmd {
		case `help`, `-?`, `--help`, `-help`:
			return h.Help()
		case `v`, `val`, `valid`, `validate`:
			return Validate(h.out, h.err, verbose).Do(args)
		case `-v`, `--verbose`:
			verbose = true
		default:
			return h.UnknownOption(cmd)
		}
	}
	h.Help()
	return 1
}

func (h *dgoCommand) Help() int {
	util.WriteString(h.out, `dgo: a command line tool to interact with the dgo type system

Usage: 
  dgo [command] [flags]

Available commands:
  help        Shows this help
  validate    Validates input file against a parameter description

Available flags:
  -v, --verbose   Be verbose in output

Use "dgo [command] --help" for more information about a command.
`)
	return 0
}
