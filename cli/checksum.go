package cli

import (
  "github.com/gorobot-library/orca/checksum"

  "github.com/spf13/cobra"
)

var checksumCmd = &cobra.Command{
	Use:   "checksum",
	Short: "Generates checksum files.",
  Long:  `Generates checksum files.`,
  Run: func(cmd *cobra.Command, args []string) error {
    return checksum.GenerateShasums()
  },
}
