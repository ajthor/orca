package cli

import (
  "github.com/gorobot-library/orca/checksum"
  "github.com/gorobot-library/orca/config"
  log "github.com/gorobot-library/orca/logger"

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
      "remote.mirror",
      "remote.file",
      "checksum",
      "checksum.versions",
    }

    if err := config.HasRequired(cfg, req); err != nil {
      log.Fatal(err)
    }

    versions := cfg.GetStringSlice("checksum.versions")
    remoteCfg := cfg.Sub("remote")

    checksum.GenerateChecksums(remoteCfg, versions)

  },
}
