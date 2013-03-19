package commands

import (
	"flag"
	"github.com/jgrocho/pass/options"
)

func init() {
	cmd := new(initialize)
	// Define flags here
	Register("init", cmd)
}

type initialize struct {
	flags flag.FlagSet
}

func (cmd *initialize) Usage() {
	// Print usage message for initialize here.
}

func (cmd *initialize) Flags() flag.FlagSet {
	return cmd.flags
}

func (cmd *initialize) Run(globals options.Options, args []string) error {
	if err := initRepo(string(globals.Prefix)); err != nil {
		return &CmdError{4, "failed to initialize repository"}
	}
	return nil
}
