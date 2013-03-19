package commands

import (
	"flag"
	"github.com/jgrocho/pass/options"
	"sync"
)

type Command interface {
	Flags() flag.FlagSet
	Usage()
	Run(options.Options, []string) error
}

var (
	commands map[string]Command
	mutex    sync.RWMutex
)

func Register(name string, cmd Command) {
	if commands == nil {
		commands = make(map[string]Command)
	}

	mutex.Lock()
	defer mutex.Unlock()

	if name == "" {
		panic("commands: register with empty name")
	}

	if cmd == nil {
		panic("commands: registering " + name + " with nil cmd")
	}

	if _, defined := commands[name]; defined {
		panic("commands: multiple registrations for " + name)
	}

	commands[name] = cmd
}

func Get(name string) Command {
	mutex.RLock()
	defer mutex.RUnlock()
	return commands[name]
}
