package cli

import (
  "path/filepath"
  // "net/url"

  "github.com/gorobot-library/orca/client"
  "github.com/gorobot-library/orca/manifest"

  log "github.com/gorobot/robologger"
  "github.com/spf13/cobra"
)

var shasumCmd = &cobra.Command{
	Use:   "shasum [name|all]",
	Short: "Generate the shasum file.",
  Long:  shasumCmdDesc,
  PreRun: func(cmd *cobra.Command, args []string) {
    HasArgs(cmd, args, 1)
  },
  Run: func(cmd *cobra.Command, args []string) {
    manifest, _ := GetNamesFromArgs(cmd, args)

    // Get the manifest.
    schemaFile := GetManifestFileName(cmd)
    s, err := OpenManifestSchema(schemaFile)
    if err != nil {
      log.Fatal(err)
    }

    switch manifest {
    case "all":
      for _, m := range s.Manifests {
        err := GenerateShasums(cmd, args, &m)
        if err != nil {
          log.Fatal(err)
        }
      }
    default:
      // Get the manifest from the schema.
      m, err := GetManifest(manifest, s)
      if err != nil {
        log.Fatal(err)
      }

      err = GenerateShasums(cmd, args, m)
      if err != nil {
        log.Fatal(err)
      }
    }

    log.Info("done")
  },
}
//     //
//     //
//     // // // Get the versions from the version flag.
//     // // if len(versions) == 0 {
//     // //   log.Fatal(ErrVersionNotSpecified)
//     // // }
//     // //
//     // // for _, v := range versions {
//     // //   if ok := HasVersion(v, m); !ok {
//     // //     log.Fatal(ErrVersionNotFound)
//     // //   }
//     // // }
//     // //
//     //
//
//
//
//     // cfg, err := config.New("manifest.json")
//     // if err != nil {
//     //   log.Fatal(err)
//     // }
//     //
//     // req := []string{
//     //   "version",
//     //   "manifests",
//     // }
//     // if err := config.HasRequired(cfg, req); err != nil {
//     //   log.Fatal(err)
//     // }
//     //
//     // // looks for the manifest file, which can be specified in the options
//     // // passed to the Manifest structure or be left as the default. It then
//     // // looks for a file named manifest.xyz, where the type of the file can be
//     // // any config compatible with Viper.
//     // schema := manifest.New()
//     //
//     // schema.Unmarshal(cfg)
//     //
//     // c := client.NewClient(&client.ClientOptions{})
//     //
//     // for _, m := range schema.Manifests {
//     //   // Generate the filenames and urls to be used to download/save the binaries.
//     //   files := m.GenerateFilenames()
//     //   urls := m.GenerateURLs(files)
//     //
//     //   hashes, err := c.GenerateShasums(files, urls)
//     //   if err != nil {
//     //     log.Fatal(err)
//     //   }
//     //
//     //   shaFile := filepath.Join(*c.Directory, m.Name, "SHASUM256.txt")
//     //
//     //   err = c.AppendShasumFile(shaFile, hashes)
//     //   if err != nil {
//     //     log.Fatal(err)
//     //   }
//     // }
//     //
//     // c.Close()
//
//     log.Info("Done.")
//   },
// }
//
// GenerateShasums creates shasums for a particular manifest, named by the argument.
func GenerateShasums(cmd *cobra.Command, args []string, m *manifest.Manifest) error {
  _, image := GetNamesFromArgs(cmd, args)
  versions, _ := cmd.Flags().GetStringSlice("version")

  img, err := GetImage(image, m)
  if err != nil {
    log.Fatal(err)
  }

  // Create a new client.
  msg := log.Info("Initializing client...")

  cli := client.NewClient(nil, &client.ClientOptions{
    Directory: &m.Name,
  })
  defer cli.Close()

  log.Status(log.OK, msg)

//   // Try to fetch the shasum file, if one exists.
//   u, _ := url.Parse(m.Remote.Mirror)
//   u.Path = filepath.Join(u.Path, "SHASUM256.txt")
//   dlurl := u.String()
//
//   log.Info("Attempting to fetch shasums from remote.")
//
//   shasumFile, err := client.FetchShasumFile(dlurl)
//   if err != nil {
//     log.Warn("File could not be downloaded.")
//     log.Warn(err)
//   }
//   if err == nil {
//     // Check the file to ensure the hash for the specified versions exist in
//     // the downloaded file.
//     files := []string{}
//     for _, v := range m.Versions {
//       fn := m.GetRemoteFilename(v)
//       log.Infof("Checking for %s", fn)
//
//       files = append(files, fn)
//     }
//
//     for _, f := range files {
//       shasumFile.Find(f)
//     }
//   }
//
//   // If no shasum file exists, download the files and generate shasums manually.
//   // m, err := GetManifest(name, s)
//   // if err != nil {
//   //   return err
//   // }

  files := make([]string, 0)
  urls := make([]string, 0)

  if len(versions) == 0 {
    for _, v := range img.Versions {
      files = append(files, img.GetRemoteFile(v))
      urls = append(urls, img.GetRemoteURL(v))
    }
  }

  shasumFilePath := GetShasumFilePath(cmd, *cli.Directory)
  shasumFile, err := client.CreateShasumFile(shasumFilePath)
  if err != nil {
    log.Fatal(err)
  }

  shasums, err := cli.GenerateShasums(files, urls)
  if err != nil {
    log.Fatal(err)
  }

  err = shasumFile.Write(shasums)
  if err != nil {
    log.Fatal(err)
  }

  return nil
}

// GetShasumFilePath returns the path to the shasum file from the flags. If the
// name is an empty string, it returns the default: name/SHASUM256.txt
func GetShasumFilePath(cmd *cobra.Command, name string) (path string) {
  path, _ = cmd.Flags().GetString("sha-file")

  if path == "" {
    path = filepath.Join(name, "SHASUM256.txt")
  }

  return
}

// Command description for help message.
var shasumCmdDesc = `Generate the shasum file.`
