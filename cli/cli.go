package cli

import (
  "strings"

  "github.com/spf13/cobra"
)

func SetupCLIRootCmd(rootCmd *cobra.Command)  {
  rootCmd.InitDefaultHelpCmd()
  rootCmd.InitDefaultHelpFlag()
  // rootCmd.PersistentFlags().MarkHidden("help")

  rootCmd.AddCommand(buildCmd)
  buildCmd.Flags().StringP("manifest", "", "", "Path to manifest file.")
  buildCmd.Flags().StringP("image", "", "", "Name of image to build.")
  buildCmd.Flags().String("sha-file", "", "Path to shasum file.")
  buildCmd.Flags().StringSliceP("tag", "t", []string{}, "Tag(s) for build.")
  buildCmd.Flags().StringP("version", "v", "", "Version to build.")

  rootCmd.AddCommand(shasumCmd)
  shasumCmd.Flags().StringSliceP("version", "v", []string{}, "Version(s) to create shasums for.")
  // shasumCmd.Flags().StringP("out", "o", "", "Output file.")
  // shasumCmd.Flags().BoolP("force", "f", false, "Force download of all files.")

  // rootCmd.AddCommand(initCmd)
  // shasumCmd.Flags().StringP("path", "C", "", "Relative output path.")
}

// GetNamesFromArgs returns the name from the arguments.
func GetNamesFromArgs(cmd *cobra.Command, args []string) (string, string) {
  var manifest, image string
  manifest = args[0]
  image, _ = cmd.Flags().GetString("image")

  if contains := strings.Contains(manifest, "/"); contains {
    str := strings.Split(manifest, "/")
    manifest = str[0]
    image = str[1]
  }

  return manifest, image
}
