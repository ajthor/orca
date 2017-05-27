package cli

import (
  create "github.com/gorobot-library/orca/init"

  "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates initial Docker files.",
  Long:  `Generates initial Docker files.`,
  Run: func(cmd *cobra.Command, args []string) {
    create.Initialize()
  },
}
