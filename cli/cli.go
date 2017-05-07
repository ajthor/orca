package cli

import (
  "fmt"

  "github.com/spf13/cobra"
)

func SetupCLIRootCmd(rootCmd *cobra.Command)  {
  rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")

  addRootCommands(rootCmd)
}

func addRootCommands(rootCmd *cobra.Command)  {
  rootCmd.AddCommand(buildCmd)
  rootCmd.AddCommand(checksumCmd)
  rootCmd.AddCommand(createCmd)
}
