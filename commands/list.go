package commands

import (
	"flag"
	"fmt"
	"github.com/jgrocho/pass/options"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type list struct {
	flags   flag.FlagSet
	all     bool
	color   bool

	colors  map[string]string
	regexes map[string]*regexp.Regexp
}

func init() {
	cmd := list{
		colors: map[string]string{"dir": "01;34", "normal": "0", "default": "01;34", "off": "0"},
		regexes: map[string]*regexp.Regexp{
			"dir":    regexp.MustCompile(":?di=([^:]+):?"),
			"normal": regexp.MustCompile(":?no=([^:]+):?"),
		},
	}

	cmd.flags.BoolVar(&cmd.all, "all", false, "Show hidden passwords")
	cmd.flags.BoolVar(&cmd.color, "color", true, "Colorize output")

	Register("list", &cmd)
	Register("ls", &cmd)
}

func (cmd *list) setColors() {
	lsColors := os.Getenv("LS_COLORS")
	if lsColors == "" {
		return
	}

	dirMatch := cmd.regexes["dir"].FindStringSubmatch(lsColors)
	if dirMatch != nil {
		cmd.colors["dir"] = dirMatch[1]
	}

	normalMatch := cmd.regexes["normal"].FindStringSubmatch(lsColors)
	if normalMatch != nil {
		cmd.colors["normal"] = normalMatch[1]
	}
}

func (cmd *list) walker(prefix string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == prefix {
			return nil
		}

		path = path[len(prefix+string(filepath.Separator)):]
		parts := strings.Split(path, string(filepath.Separator))
		name := parts[len(parts)-1]
		if name == ".git" {
			return filepath.SkipDir
		}
		if !strings.HasPrefix(name, ".") || cmd.all {
			for i := 0; i < len(parts)-2; i++ {
				fmt.Print("    ")
			}
			if len(parts) > 1 {
				fmt.Print(" \u2514\u2500 ")
			}
			if info.IsDir() {
				if cmd.color {
					fmt.Printf("\x1b[%sm%s\x1b[%sm\n", cmd.colors["dir"], name, cmd.colors["off"])
				} else {
					fmt.Printf("%s\n", name)
				}
			} else {
				if strings.HasSuffix(name, ".gpg") {
					name = name[:len(name)-4]
					colorHere := cmd.colors["normal"]
					if name == "" {
						name = "(default)"
						colorHere = cmd.colors["default"]
					}
					if cmd.color {
						fmt.Printf("\x1b[%sm%s\x1b[%sm\n", colorHere, name, cmd.colors["off"])
					} else {
						fmt.Printf("%s\n", name)
					}
				}
			}
		}
		return nil
	}
}

func (cmd *list) Usage() {
	// Print usage message for list here.
}

func (cmd *list) Flags() flag.FlagSet {
	return cmd.flags
}

func (cmd *list) Run(globals options.Options, args []string) error {
	if cmd.color {
		cmd.setColors()
	}

	prefix := globals.Prefix.String()
	prefixInfo, err := os.Stat(prefix)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrPrefixNotExist(prefix)
		} else {
			return ErrPrefixInaccessible(prefix)
		}
	}
	if !prefixInfo.IsDir() {
		return ErrPrefixNotDir(prefix)
	}

	filepath.Walk(prefix, cmd.walker(prefix))

	return nil
}
