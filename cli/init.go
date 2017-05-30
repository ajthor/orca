package cli

import (
  "os"

  "github.com/gorobot-library/orca/config"
  "github.com/gorobot-library/orca/initialize"
  log "github.com/gorobot-library/orca/logger"

  "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates initial Docker files.",
  Long:  `Generates initial Docker files.`,
  Run: func(cmd *cobra.Command, args []string) {
    // Read the configuration file.
    if _, err := os.Stat("./orca.toml"); os.IsNotExist(err) {
      log.Warn("No configuration file found.")
      initialize.GenerateConfigFile()
    }

    cfg, err := config.NewConfig("orca")
    if err != nil {
      log.Fatal("Configuration file not found.")
    }

    req := []string{
      "name",
      "build.base",
      "build.version",
    }

    if err := config.HasRequired(cfg, req); err != nil {
      log.Fatal(err)
    }

    initialize.Initialize(cfg)
  },
}
