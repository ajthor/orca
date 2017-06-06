package cli

import (
  "path/filepath"

  "github.com/gorobot-library/orca/manifest"
  "github.com/gorobot-library/orca/client"

  log "github.com/gorobot/robologger"
  "github.com/spf13/cobra"
)

var shasumCmd = &cobra.Command{
	Use:   "shasum",
	Short: "Generates shasum files.",
  Long:  `Generates shasum files using the remote file and versions specified in the manifest.`,
  Run: func(cmd *cobra.Command, args []string) {

    schema := manifest.Load("manifest.json")

    c := client.New(&client.ClientOptions{})

    for _, m := range schema.Manifests {
      // Generate the filenames and urls to be used to download/save the binaries.
      files := m.GenerateFilenames()
      urls := m.GenerateURLs(files)

      hashes, err := c.GenerateShasums(files, urls)
      if err != nil {
        log.Fatal(err)
      }

      shaFile := filepath.Join(*c.Directory, m.Name, "SHASUM256.txt")

      err = c.AppendShasumFile(shaFile, hashes)
      if err != nil {
        log.Fatal(err)
      }
    }

    c.Close()

    log.Info("Done.")
  },
}
