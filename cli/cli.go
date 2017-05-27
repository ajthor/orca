package cli

import (
  "github.com/spf13/cobra"
)

func SetupCLIRootCmd(rootCmd *cobra.Command)  {
  // Make sure that if the user attaches the '-h' flag, that we will always
  // display the help text.
  rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")

  // flags := cmd.Flags()
	// flags.BoolVarP(&opts.version, "version", "v", false, "Print version information and quit")

  addRootCommands(rootCmd)
}

func addRootCommands(rootCmd *cobra.Command)  {
  // Add all commands to the root command. Sub-commands (if any) will be stored
  // in the respective files.
  rootCmd.AddCommand(buildCmd)
  rootCmd.AddCommand(checksumCmd)
  // rootCmd.AddCommand(initCmd)
}
