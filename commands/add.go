package commands

import (
	"flag"
	"fmt"
	"github.com/jgrocho/pass/options"
	"github.com/jgrocho/passphrase"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	cmd := new(add)
	cmd.flags.BoolVar(&cmd.force, "force", false, "Overwrite existing passphrase")
	cmd.flags.StringVar(&cmd.edit, "edit", "", "Edit file with editor")
	Register("add", cmd)
}

type add struct {
	flags flag.FlagSet
	force bool
	edit  string
}

func (cmd *add) Usage() {
	// Print usage message for add here.
}

func (cmd *add) Flags() flag.FlagSet {
	return cmd.flags
}

func ErrNotAdding(name, message string) *CmdError {
	return &CmdError{4, fmt.Sprintf("not adding %s: %s", name, message)}
}

func (cmd *add) Run(globals options.Options, args []string) error {
	prefix := string(globals.Prefix)
	name, passfile, err := getNameAndFile(prefix, args)
	if err != nil {
		return err
	}

	if _, err := os.Stat(passfile); err == nil && !cmd.force {
		return &CmdError{4, "an entry already exists for " + name + ". Use -force to overwrite it"}
	}

	passdir := filepath.Dir(passfile)
	if _, err := os.Stat(passdir); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(passdir, 0777); err != nil {
			return ErrNotAdding(name, "could not create directory " + passdir)
		}
	}

	outfile, err := os.Create(passfile)
	if err != nil {
		return ErrNotAdding(name, "could not open file")
	}
	defer outfile.Close()

	var input io.Reader = os.Stdin
	if cmd.edit != "" {
		tempfile, err := ioutil.TempFile("", "pass_" + name)
		defer os.Remove(tempfile.Name())
		if err != nil {
			return ErrNotAdding(name, "could not create temporary file for editing")
		}
		editor := exec.Command(cmd.edit, tempfile.Name())
		editor.Stdin = os.Stdin
		editor.Stdout = os.Stdout
		editor.Stderr = os.Stderr
		if err := editor.Run(); err != nil {
			return ErrNotAdding(name, "failed running editor: " + cmd.edit)
		}
		input = tempfile
	} else if isTerminal(os.Stdin) {
		passwd, err := passphrase.GetPassphrase("", "", "Please enter a password for " + name, "", true, true)
		if err != nil {
			return ErrNotAdding(name, "could not get a passphrase")
		}
		input = strings.NewReader(passwd + "\n")
	}

	if err := encrypt(string(globals.PubRing), input, outfile); err != nil {
		return ErrNotAdding(name, "could not encrypt passphrase")
	}

	relative := passfile[len(prefix)+1:]
	if err := addAndCommit(prefix, relative, "Add password for " + name); err != nil {
		return &CmdError{4, "failed to commit passphrase to repository"}
	}

	return nil
}
