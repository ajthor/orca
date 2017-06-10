package main

import (
  "github.com/gorobot-library/orca/cli"

  log "github.com/gorobot/robologger"
  "github.com/spf13/cobra"
)

func newCLICommand() *cobra.Command {
  // Initialize the "root command".
  // By default, running the `orca` command will display the usage message.
  cmd := &cobra.Command{
      Use:   "orca",
      Short: "A simple Docker image build tool.",
      Long:  rootCmdDesc,
    }

  // Perform initialization steps, such as attaching commands to the root
  // command. Commands and initialization are handled by the `cli/cli.go` file.
  cli.SetupCLIRootCmd(cmd)

  return cmd
}

func main() {
  cmd := newCLICommand()
  // Run the command.
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Command description for help message.
var rootCmdDesc = `Orca 0.0.1
A simple image build tool.

Orca is a Docker image build tool that uses templates to create the build
context. It provides tools for building Docker images and generating shasums.

Orca is designed to be run either form the command line or inside a Docker
container. For example:

$ docker run -it --rm -v "$PWD:/" orca:0.0.1 <command>`
