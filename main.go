package main

import (
	"flag"
	"fmt"
	"github.com/jgrocho/pass/commands"
	"github.com/jgrocho/pass/options"
	"os"
	"os/user"
	"path/filepath"
)

const (
	VER_MAJOR uint = 0
	VER_MINOR      = 1
	VER_PATCH      = 0
)

var (
	flags   flag.FlagSet
	help    bool
	version bool
	globals options.Options
)

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic("pass: could not determine current user")
	}
	homedir := currentUser.HomeDir

	flags.BoolVar(&help, "help", false, "Show this help message")
	flags.BoolVar(&version, "version", false, "Show version number")

	globals.Prefix = options.FilePath(filepath.Join(homedir, ".pass"))
	flags.Var(&globals.Prefix, "prefix", "Set directory to store data")

	globals.PubRing = options.FilePath(filepath.Join(homedir, ".gnupg", "pubring.gpg"))
	flags.Var(&globals.PubRing, "pubring", "Set public ring file")

	globals.SecRing = options.FilePath(filepath.Join(homedir, ".gnupg", "secring.gpg"))
	flags.Var(&globals.SecRing, "secring", "Set private ring file")
}

func printUsage(out *os.File) {
	fmt.Fprintf(out, "%s: usage\n", os.Args[0])
}

func printVersion() {
	fmt.Printf("%s: version %d.%d.%d\n", os.Args[0], VER_MAJOR, VER_MINOR, VER_PATCH)
}

func processCommand(args []string) error {
	if len(args) < 1 {
		return commands.ErrNoCommand
	}

	command := args[0]
	cmd := commands.Get(command)
	if cmd == nil {
		return commands.ErrUnknownCommand(command)
	}

	cmdFlags := cmd.Flags()
	cmdHelp := cmdFlags.Bool("help", false, "Show this help message")
	cmdFlags.Parse(args[1:])

	if *cmdHelp {
		cmd.Usage()
		return nil
	}

	return cmd.Run(globals, cmdFlags.Args())
}

func main() {
	flags.Parse(os.Args[1:])

	if help {
		printUsage(os.Stdout)
		return
	}

	if version {
		printVersion()
		return
	}

	if err := processCommand(flags.Args()); err != nil {
		if cmdErr, ok := err.(*commands.CmdError); ok {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], cmdErr.Message)
			os.Exit(cmdErr.Code)
		}
	}
}
