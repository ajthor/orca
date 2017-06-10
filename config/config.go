package config

import (
  "errors"
  "strings"

  "path/filepath"

  "github.com/spf13/viper"
)

var (
  ErrConfigNotSpecified = errors.New("Configuration file not specified.")
)

// New creates a new configuration object. It takes a string, the name of a
// file which will be loaded as the configuration.
func NewConfig(name string) (*viper.Viper, error) {
  // Ensure the path variable is set.
  if name == "" {
    return nil, ErrConfigNotSpecified
  }

  // Create a new viper config.
  cfg := viper.New()

  // If the included path has a directory, we include that directory in the
  // list of paths we search for our file in.
  dir, base := filepath.Split(name)
  if dir != "" {
    cfg.AddConfigPath(dir)
  }

  // Trim the extension off of the base name.
  base = strings.TrimSuffix(base, filepath.Ext(base))

  // Change the config name.
  cfg.SetConfigName(base)

  // Set default paths to look for the orca config file.
  cfg.AddConfigPath("/")
  cfg.AddConfigPath("$HOME/.orca")
  cfg.AddConfigPath(".")

  // Actually read in the configuration.
  err := cfg.ReadInConfig()
  if err != nil {
    return nil, err
  }

  return cfg, nil
}

func HasRequired(cfg *viper.Viper, fields []string) bool {
  for _, v := range fields {
    if set := cfg.IsSet(v); !set {
      return false
    }
  }
  return true
}
