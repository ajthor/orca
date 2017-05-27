package cli

import (
  "log"
  
  "github.com/gorobot-library/orca/builder"
  "github.com/gorobot-library/orca/config"

  "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images.",
  Long:  `Builds Docker images.`,
  Run: func(cmd *cobra.Command, args []string) {
    // Read the configuration file.
    cfg := config.NewConfig("orca")
    // if err := config.ReadConfig(cfg); err != nil {
    //   log.Fatal(err)
    // }

    req := []string{
      "build",
      "build.version",
    }

    if err := config.HasRequired(cfg, req); err != nil {
      log.Fatal(err)
    }

    builder.Build()
  },
}
