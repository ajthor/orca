package cli

import (
  "context"
  "errors"

  "github.com/gorobot-library/orca/client"
  // "github.com/gorobot-library/orca/config"
  // "github.com/gorobot-library/orca/manifest"

  "github.com/docker/docker/api/types"
  log "github.com/gorobot/robologger"
  "github.com/spf13/cobra"
)

var (
  ErrNameNotSpecified = errors.New("No name specified.")
  ErrVersionNotSpecified = errors.New("No version specified.")
  ErrTagNotSpecified = errors.New("No tag specified.")
  ErrManifestNotSpecified = errors.New("No manifest specified.")
  ErrVersionNotFound = errors.New("Version not found in manifest.")
  ErrShasumFileNotFound = errors.New("Shasum file not found.")
)

// buildCmd is a cobra command that calls the Build function when used.
var buildCmd = &cobra.Command{
	Use:   "build [name]",
	Short: "Build Docker images.",
  Long:  buildCmdDesc,
  PreRun: func(cmd *cobra.Command, args []string) {
    HasArgs(cmd, args, 1)
    FlagChanged(cmd, "tag")
  },
  Run: Build,
}

// Build builds an image using the options specified. It requires the user to
// specify a version on the command line and optionally a tag for the image.
func Build(cmd *cobra.Command, args []string) {
  manifest, image := GetNamesFromArgs(cmd, args)
  version, _ := cmd.Flags().GetString("version")
  tags, _ := cmd.Flags().GetStringSlice("tag")

  // Get the manifest.
  schemaFile := GetManifestFileName(cmd)
  s, err := OpenManifestSchema(schemaFile)
  if err != nil {
    log.Fatal(err)
  }

  m, err := GetManifest(manifest, s)
  if err != nil {
    log.Fatal(err)
  }

  img, err := GetImage(image, m)
  if err != nil {
    log.Fatal(err)
  }

  // Create a new client.
  msg := log.Info("Initializing client...")
  cli := client.NewClient(nil, &client.ClientOptions{})
  defer cli.Close()
  log.Status(log.OK, msg)

  // Make sure the version is in the manifest, or if no version is specified,
  // default to the latest version.
  switch version {
  case "latest", "":
    version = img.GetLatestVersion()
  default:
    if ok := img.HasVersion(version); !ok {
      log.Fatal(ErrVersionNotFound)
    }
  }

  // Get the file name.
  file := img.GetRemoteFile(version)
  // Get the url.
  url := img.GetRemoteURL(version)

  // Get the path to the shasum file.
  shasumFilePath := GetShasumFilePath(cmd, *cli.Directory)

  var hash client.Shasum
  // Get hash from shasum file.
  sf, err := client.OpenShasumFile(shasumFilePath)
  if err != nil {
    log.Warn(ErrShasumFileNotFound)
  } else {
    shasum := sf.Find(file)
    if shasum != nil {
      log.Fatal(err)
    }
  }

  data := make(map[string]interface{})
  data["Name"] = manifest
  data["Base"] = img.BaseImage
  data["Version"] = version
  data["Hash"] = hash.String()
  data["File"] = file
  data["Mirror"] = img.Remote.Mirror
  data["URL"] = url

  // Create the build context.
  buildCtx := client.NewContext(img, &client.ContextOptions{
    Directory: cli.Directory,
    Data: data,
  })

  // Generate build tags.
  // tags := []string{}
  // for _, t := range m.GetTags() {
  //   if contains := strings.Contains(t, version); contains {
  //     tags = append(tags, t)
  //   }
  // }
  // if len(tags) == 0 {
  //   tags = append(tags, strings.Join(m.Repo, m.Name, ":", version, "-"))
  // }

  log.Info("Building...")
  for _, t := range tags {
    log.Debugf("%s", t)
  }

  err = cli.ImageBuild(context.Background(), buildCtx, types.ImageBuildOptions{
    Tags: tags,
  })
  if err != nil {
    log.Fatal(err)
  }

  log.Info("done")
}

// CreateBuildContext creates a build context associated with the image
// specified by the arguments.
// func CreateBuildContext(m string, name string, version string) (*client.Context, error) {
//
// }

// Command description for help message.
var buildCmdDesc = `Build Docker images.`
