package main

import (
  "fmt"
  "os"

  "github.com/gorobot-library/orca/cli"
  "github.com/spf13/cobra"
)

func newCLICommand() *cobra.Command {
  // Initialize the "root command".
  // By default, simply running the `orca` command will not perform an action.
  cmd := &cobra.Command{
      Use:   "orca",
      Short: "Orca is a simple Docker image build tool.",
      Long:  `Orca is a simple Docker image build tool.`,
    }

  // Perform initialization steps, such as attaching commands to the root
  // command. Commands and initialization are handled by the cli/cli.go file.
  cli.SetupCLIRootCmd(cmd)

  return cmd
}

func main() {
  cmd := newCLICommand()
  // Run the command.
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
