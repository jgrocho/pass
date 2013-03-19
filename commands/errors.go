package commands

type CmdError struct {
	Code    int
	Message string
}

func (ce *CmdError) Error() string {
	return ce.Message
}

var ErrNoCommand = &CmdError{1, "command required"}

func ErrUnknownCommand(command string) *CmdError {
	return &CmdError{2, "unknown command: " + command}
}

func ErrPrefixNotExist(prefix string) *CmdError {
	return &CmdError{3, "prefix directory " + prefix + " does not exist"}
}

func ErrPrefixInaccessible(prefix string) *CmdError {
	return &CmdError{3, "prefix directory " + prefix + " is not accessible"}
}

func ErrPrefixNotDir(prefix string) *CmdError {
	return &CmdError{3, "prefix path " + prefix + " is not a directory"}
}
