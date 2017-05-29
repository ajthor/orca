package cli

import (
  "github.com/gorobot-library/orca/builder"
  "github.com/gorobot-library/orca/config"
  log "github.com/gorobot-library/orca/logger"

  "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images.",
  Long:  `Builds Docker images.`,
  Run: func(cmd *cobra.Command, args []string) {
    // Read the configuration file.
    cfg := config.NewConfig("orca")

    req := []string{
      "build",
      "build.version",
    }

    if err := config.HasRequired(cfg, req); err != nil {
      log.Fatal(err)
    }

    buildCfg := cfg.Sub("build")

    builder.Validate(buildCfg)
    builder.Build(buildCfg)
  },
}
