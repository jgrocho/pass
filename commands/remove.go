package commands

import (
	"flag"
	"github.com/jgrocho/pass/options"
	"log"
	"os"
	"path/filepath"
)

func init() {
	cmd := new(remove)
	// Define flags here
	Register("rm", cmd)
	Register("remove", cmd)
}

type remove struct {
	flags flag.FlagSet
}

func (cmd *remove) Usage() {
	// Print usage message for remove here.
}

func (cmd *remove) Flags() flag.FlagSet {
	return cmd.flags
}

func (cmd *remove) Run(globals options.Options, args []string) error {
	prefix := string(globals.Prefix)
	name, passfile, err := getNameAndFile(prefix, args)
	if err != nil {
		return err
	}

	if _, err := os.Stat(passfile); err != nil && os.IsNotExist(err) {
		return &CmdError{4, "passphrase does not exist for " + name}
	}

	relative := passfile[len(prefix)+1:]
	if err := removeAndCommit(prefix, relative, "Remove password for " + name); err != nil {
		log.Printf("%s\n", err)
		return &CmdError{4, "failed to remove passphrase from repository"}
	}

	if err := os.Remove(passfile); err != nil {
		return &CmdError{4, "could not delete passphrase for " + name}
	}

	prefixInfo, err := os.Stat(prefix)
	if err != nil {
		return nil
	}
	passdir := filepath.Dir(passfile)
	passdirInfo, err := os.Stat(passdir)
	if err != nil {
		return nil
	}
	for !os.SameFile(prefixInfo, passdirInfo) {
		os.Remove(passdir)
		passdir = filepath.Dir(passdir)
		passdirInfo, err = os.Stat(passdir)
		if err != nil {
			return nil
		}
	}

	return nil
}
