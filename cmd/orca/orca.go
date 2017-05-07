package main

import (
  "fmt"
  "os"

  "github.com/gorobot-library/orca/cli"
  "github.com/spf13/cobra"
)

func newCLICommand() *cobra.Command {
  cmd := &cobra.Command{
      Use:   "orca",
      Short: "Orca is a simple Docker image build tool.",
      Long:  `Orca is a simple Docker image build tool.`,
      Run: func(cmd *cobra.Command, args []string) {
        // Do Stuff Here
      },
    }

  cli.SetupCLIRootCmd(cmd)
  
  // flags := cmd.Flags()
	// flags.BoolVarP(&opts.version, "version", "v", false, "Print version information and quit")

  return cmd
}

func main() {
  cmd := newCLICommand()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
