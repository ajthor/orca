package cli

import (
  "github.com/gorobot-library/orca/builder"
  "github.com/gorobot-library/orca/checksum"
  "github.com/gorobot-library/orca/config"

  log "github.com/gorobot/robologger"
  "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images.",
  Long:  `Builds Docker images.`,
  Run: func(cmd *cobra.Command, args []string) {
    var (
      base string
      version string
      dockerfile string
      mirror string
      file string
      uri string
      hash string
    )

    // Read the configuration file.
    cfg, err := config.New("orca")
    if err != nil {
      log.Fatal(err)
    }

    req := []string{
      "build",
      "build.version",
    }

    if err := config.HasRequired(cfg, req); err != nil {
      log.Fatal(err)
    }

    // Here, we get all information necessary for the build from the config
    // file. We need to retreive the values to pass into the builder.

    buildCfg := cfg.Sub("build")
    base = buildCfg.GetString("base")
    version = buildCfg.GetString("version")
    dockerfile = buildCfg.GetString("dockerfile")
    includes := buildCfg.GetStringSlice("includes")
    tags := buildCfg.GetStringSlice("tags")

    remoteCfg := cfg.Sub("remote")
    mirror = remoteCfg.GetString("mirror")
    file = remoteCfg.GetString("file")

    if file != "" {
      file, err = config.ParseFilename(file, version)
      if err != nil {
        log.Fatal(err)
      }

      uri, err = config.ParseURL(mirror, file)
      if err != nil {
        log.Fatal(err)
      }

      c := checksum.New("")
      hash, _ = c.GetShasum(file)
    }

    opts := &builder.BuildOptions{
      Base: base,
      Version: version,
      Dockerfile: dockerfile,
      Includes: includes,
      Tags: tags,
      Mirror: mirror,
      File: file,
      URL: uri,
      Hash: hash,
    }

    b := builder.New(builder.ClientOptions{
      Host: "unix:///var/run/docker.sock",
      Version: "v1.29",
    })

    b.Build(opts)

    log.Info("Done.")
  },
}
