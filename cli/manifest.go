package cli
//
import (
  // "bufio"
  // "context"
  "errors"

  // "github.com/gorobot-library/orca/client"
  "github.com/gorobot-library/orca/config"
  "github.com/gorobot-library/orca/manifest"

  // "github.com/docker/docker/api/types"
  log "github.com/gorobot/robologger"
  "github.com/spf13/cobra"
)

var (
  ErrConfigFileNotFound = errors.New("Config file not found.")
  ErrManifestNotFound = errors.New("Manifest not found.")
  ErrImageNotFound = errors.New("Image not found.")
)

// var manifestCmd = &cobra.Command{
// 	Use:   "manifest",
// 	Short: "Command for working with the build manifest.",
//   Long:  `Command for working with the build manifest.`,
// }
//
// var manifestAddCmd = &cobra.Command{
// 	Use:   "add [version]",
// 	Short: "Adds a version to the manifest.",
//   Long:  `Adds a version to the manifest.
//   eg: orca manifest add 0.1.2`,
//   Run: func(cmd *cobra.Command, args []string) {
//     cfgPath := "manifest.json"
//     cfg, err := config.New(&cfgPath)
//     if err != nil {
//       log.Fatal(err)
//     }
//
//     m := manifest.Unmarshal(cfg)
//
//     log.Info("Done.")
//   },
// }

// GetManifestFileName returns the manifest file path, either from the flags or
// returns the default.
func GetManifestFileName(cmd *cobra.Command) string {
  var file string
  file, _ = cmd.Flags().GetString("manifest")

  if file == "" {
    file = "manifest.json"
  }

  return file
}

// OpenManifestSchema fetches the manifest from the manifest config file.
//
// It returns an ErrConfigFileNotFound if the config file cannot be found.
func OpenManifestSchema(name string) (*manifest.ManifestSchema, error) {
  cfg, err := config.NewConfig(name)
  if err != nil {
    return nil, ErrConfigFileNotFound
  }

  s := manifest.NewSchema()
  s.Unmarshal(cfg)

  log.Printf("Schema Version: v%d", s.SchemaVersion)
  return s, nil
}

// GetManifest fetches the manifest from the manifest schema.
//
// It returns an ErrManifestNotFound if the manifest cannot be found in the
// config file.
func GetManifest(name string, s *manifest.ManifestSchema) (*manifest.Manifest, error) {
  msg := log.Info("Loading manifest...")

  m := s.Get(name)
  if m == nil {
    log.Status(log.ERR, msg)
    return nil, ErrManifestNotFound
  }

  log.Status(log.OK, msg)
  return m, nil
}

// GetImage fetches the image from the manifest.
//
// It returns an ErrImageNotFound if the image cannot be found in the
// config file.
func GetImage(name string, m *manifest.Manifest) (*manifest.Image, error) {
  msg := log.Info("Loading image...")

  img := m.Get(name)
  if m == nil {
    log.Status(log.ERR, msg)
    return nil, ErrImageNotFound
  }

  log.Status(log.OK, msg)
  return img, nil
}
