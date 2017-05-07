package cli

import (
  "github.com/gorobot-library/orca/builder"

  "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images.",
  Long:  `Builds Docker images.`,
  Run: func(cmd *cobra.Command, args []string) error {
    client := builder.NewClient()

    return builder.Build(client)
  },
}
