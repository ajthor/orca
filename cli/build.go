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

    isValid, err := builder.Validate(buildCfg)
    if err != nil {
      log.Fatal(err)
    } else if !isValid {
      log.Fatal("Cannot build image. Invalid configuration.\n")
    }

    // err = builder.Build(buildCfg)
    // if err != nil {
    //   log.Std.Fatal(err)
    // }
  },
}
