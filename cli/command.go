package cli

import (
	"io"
	"strings"

	"github.com/lyraproj/dgo/util"
)

// Command is the interface implemented by all dgo commands. Both global and sub commands.
type Command interface {
	// Do parses the arguments and runs the command
	Do([]string) int

	// Help prints the help for the command and returns an 0
	Help() int

	// MissingOption prints the missing option error for the given option and returns 1
	MissingOption(string) int

	// MissingOptionArgument prints the missing option argument error for the given option and returns 1
	MissingOptionArgument(string) int

	// RunWithCatch runs the given function, recovers error panics and reports them. It returns a non zero exit status
	// in case an error was recovered
	RunWithCatch(func() int) int

	// UnknownOption prints the unknown option error for the given option and returns 1
	UnknownOption(string) int
}

type command struct {
	name  string
	out   io.Writer
	err   io.Writer
	debug bool
}

func (h *command) RunWithCatch(runner func() int) int {
	exitCode := 0
	err := util.Catch(func() {
		exitCode = runner()
	})
	if err != nil {
		util.Fprintf(h.err, "Error: %s\n", err.Error())
		return 1
	}
	return exitCode
}

func (h *command) MissingOption(opt string) int {
	util.Fprintf(h.err, "Error: Missing required option %s\nUse '%s --help' for more help\n", opt, h.name)
	return 1
}

func (h *command) MissingOptionArgument(opt string) int {
	util.Fprintf(h.err, "Error: Option %s requires an argument\nUse '%s --help' for more help\n", opt, h.name)
	return 1
}

func (h *command) UnknownOption(opt string) int {
	if strings.HasPrefix(opt, `-`) {
		util.Fprintf(h.err, "Error: UnknownOption flag: %s\nUse '%s --help' for more help\n", opt, h.name)
	} else {
		util.Fprintf(h.err, "Error: UnknownOption argument: %s\nUse '%s --help' for more help\n", opt, h.name)
	}
	return 1
}
