package cli

import (
  "errors"

  log "github.com/gorobot/robologger"
  "github.com/spf13/cobra"
)

var (
  ErrNoArgs = errors.New("Argument required.")
  ErrNoFlag = errors.New("Flag required.")
  ErrInvalidCommand = errors.New("Invalid command.")
  ErrInvalidArgument = errors.New("Invalid argument.")
  ErrInvalidFlag = errors.New("Invalid flag.")
  ErrMissingFlag = errors.New("Missing flag.")
)

// CmdError is a wrapper for the log.Fatal function. It displays an error if
// the command was entered incorrectly and prints the command usage
func CmdError(cmd *cobra.Command, err error) {
  cmd.Usage()
  // Print a newline.
  log.Print("")
  log.Fatal(err)
}

// NoArgs ensures that there are no args passed to the command.
func NoArgs(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
    CmdError(cmd, ErrInvalidArgument)
	}
}

// HasArgs ensures the correct number of args exist.
func HasArgs(cmd *cobra.Command, args []string, num int) {
	if len(args) != num {
    CmdError(cmd, ErrNoArgs)
	}
}

// FlagChanged ensures that the specified flag is set by the user.
func FlagChanged(cmd *cobra.Command, name string) {
  v := cmd.Flags().Lookup(name)
  if !v.Changed {
    CmdError(cmd, ErrMissingFlag)
  }
}
