package cli

import (
  "github.com/gorobot-library/orca/create"

  "github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generates initial Docker files.",
  Long:  `Generates initial Docker files.`,
  Run: func(cmd *cobra.Command, args []string) error {
    return create.Create()
  },
}
