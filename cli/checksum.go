package cli

import (
  "log"
  
  "github.com/gorobot-library/orca/checksum"
  "github.com/gorobot-library/orca/config"

  "github.com/spf13/cobra"
)

var checksumCmd = &cobra.Command{
	Use:   "checksum",
	Short: "Generates checksum files.",
  Long:  `Generates checksum files.`,
  Run: func(cmd *cobra.Command, args []string) {
    // Read the configuration file.
    cfg := config.NewConfig("orca")
    // if err := config.ReadConfig(cfg); err != nil {
    //   log.Fatal(err)
    // }

    req := []string{
      "remote",
      "remote.uri",
      "remote.file",
      "checksum",
      "checksum.versions",
    }

    if err := config.HasRequired(cfg, req); err != nil {
      log.Fatal(err)
    }

    versions := cfg.GetStringSlice("checksum.versions")
    remote := cfg.Sub("remote")

    checksum.GenerateChecksums(remote, versions)

  },
}
