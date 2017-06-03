package cli

import (
  "github.com/gorobot-library/orca/checksum"
  "github.com/gorobot-library/orca/config"

  log "github.com/gorobot/robologger"
  "github.com/spf13/cobra"
)

var checksumCmd = &cobra.Command{
	Use:   "checksum",
	Short: "Generates checksum files.",
  Long:  `Generates checksum files.`,
  Run: func(cmd *cobra.Command, args []string) {
    var (
      mirror string
      file string
    )

    // Read the configuration file.
    cfg, err := config.New("orca")
    if err != nil {
      log.Fatal(err)
    }

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

    checksumCfg := cfg.Sub("checksum")
    versions := checksumCfg.GetStringSlice("versions")

    remoteCfg := cfg.Sub("remote")
    mirror = remoteCfg.GetString("mirror")
    file = remoteCfg.GetString("file")

    opts := &checksum.ChecksumOptions{
      Mirror: mirror,
      File: file,
      Versions: versions,
    }

    c := checksum.New("")

    hashes, err := c.GenerateHashes(opts)
    if err != nil {
      log.Fatal(err)
    }

    // Write the hashes to the shasum file.
    log.Info("Generating shasum file...")
    err = c.CreateShasumFile(hashes)
    if err != nil {
      log.Fatal(err)
    }

    log.Info("Done.")
  },
}
