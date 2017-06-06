package config

import (
  "fmt"
  "strings"

  "path/filepath"

  "github.com/spf13/viper"
)

// Create a new configuration object.
func New(p string) (*viper.Viper, error) {
  // Create a new viper config.
  cfg := viper.New()

  // The name of the default configuration file, without extensions.
  cfg.SetConfigName("orca")

  // If a path is specified, we need to change the configuration name and add
  // the new directory to the list of search paths.
  if p != "" {
    dir, base := filepath.Split(p)

    // If the included path has a directory, returned by the `Split` function
    // above, we include that directory in the list of directories we search
    // for our file in.
    if dir != "" {
      cfg.AddConfigPath(dir)
    }

    // Trim the extension off of the base name.
    base = strings.TrimSuffix(base, filepath.Ext(base))

    // Change the config name.
    cfg.SetConfigName(base)
  }

  // Set default paths to look for the orca config file.
  cfg.AddConfigPath("/")
  cfg.AddConfigPath("$HOME/.orca")
  cfg.AddConfigPath(".")

  // Actually read in the configuration.
  err := cfg.ReadInConfig()
  if err != nil {
    return cfg, err
  }

  return cfg, nil
}

func HasRequired(cfg *viper.Viper, requiredFields []string) error {
  for _, v := range requiredFields {
    if set := cfg.IsSet(v); !set {
      return fmt.Errorf("Required field: %s", v)
    }
  }

  return nil
}
