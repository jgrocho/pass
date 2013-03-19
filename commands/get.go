package commands

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jgrocho/pass/options"
	"os"
)

func init() {
	cmd := new(get)
	cmd.flags.BoolVar(&cmd.show, "show", false, "Print out passphrase")
	Register("get", cmd)
}

type get struct {
	flags flag.FlagSet
	show  bool
}

func (cmd *get) Usage() {
	// Print usage message for get here.
}

func (cmd *get) Flags() flag.FlagSet {
	return cmd.flags
}

func (cmd *get) Run(globals options.Options, args []string) error {
	prefix := string(globals.Prefix)
	name, passfile, err := getNameAndFile(prefix, args)
	if err != nil {
		return err
	}

	if _, err := os.Stat(passfile); err != nil && os.IsNotExist(err) {
		return &CmdError{4, "passphrase does not exist for " + name}
	}

	input, err := os.Open(passfile)
	if err != nil {
		return &CmdError{4, "could not read passphrase for " + name}
	}

	data, err := decrypt(string(globals.SecRing), input)
	if err != nil {
		return err
	}
	buffered := bufio.NewReader(data)
	passphrase, err := buffered.ReadString('\n')
	if err != nil {
		return &CmdError{4, "could not read passphrase for " + name}
	}

	if cmd.show {
		fmt.Print(passphrase)
	} else {
		if err := copyToClipboard(passphrase[:len(passphrase)-1]); err != nil {
			return &CmdError{4, "could not copy to clipboard"}
		}
	}

	return nil
}
