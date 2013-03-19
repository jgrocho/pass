package commands

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/jgrocho/pass/options"
	"math/big"
	"os"
	"path/filepath"
)

func init() {
	cmd := new(generate)
	cmd.flags.BoolVar(&cmd.show, "show", false, "Print out passphrase")
	cmd.flags.BoolVar(&cmd.force, "force", false, "Overwrite existing passphrase")
	cmd.flags.UintVar(&cmd.length, "length", 16, "Length of passphrase")
	cmd.flags.StringVar(&cmd.exclude, "exclude", "", "Exclude these characters")
	Register("generate", cmd)
}

type generate struct {
	flags   flag.FlagSet
	show    bool
	force   bool
	length  uint
	exclude string
}

func (cmd *generate) Usage() {
	// Print usage message for generate here.
}

func (cmd *generate) Flags() flag.FlagSet {
	return cmd.flags
}

func ErrNotGenerating(name, message string) *CmdError {
	return &CmdError{4, fmt.Sprintf("not generating %s: %s", name, message)}
}

func (cmd *generate) Run(globals options.Options, args []string) error {
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
			return ErrNotGenerating(name, "could not create directory "+passdir)
		}
	}

	outfile, err := os.Create(passfile)
	if err != nil {
		return ErrNotGenerating(name, "could not open file")
	}
	defer outfile.Close()

	exclude := make(map[byte]bool, len(cmd.exclude))
	for _, b := range []byte(cmd.exclude) {
		exclude[b] = true
	}

	alphabet := make([]byte, '~'-'!'+1)
	for i, r := 0, '!'; r <= '~'; r++ {
		b := byte(r)
		if !exclude[b] {
			alphabet[i] = b
			i++
		}
	}

	max := big.NewInt(int64(len(alphabet)))
	passphrase := make([]byte, cmd.length+1)
	for i := uint(0); i < cmd.length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return ErrNotGenerating(name, "could not generate random number")
		}
		passphrase[i] = alphabet[n.Int64()]
	}
	passphrase[len(passphrase)-1] = '\n'

	if err := encrypt(string(globals.PubRing), bytes.NewReader(passphrase), outfile); err != nil {
		return ErrNotGenerating(name, "could not encrypt passphrase")
	}

	if cmd.show {
		fmt.Printf("%s\n", passphrase)
	} else {
		if err := copyToClipboard(string(passphrase[:len(passphrase)-1])); err != nil {
			return &CmdError{4, "could not copy to clipboard"}
		}
	}

	relative := passfile[len(prefix)+1:]
	if err := addAndCommit(prefix, relative, "Generate password for "+name); err != nil {
		return &CmdError{4, "failed to commit passphrase to repository"}
	}

	return nil
}
