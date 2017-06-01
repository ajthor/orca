package config

// log "github.com/gorobot-library/orca/logger"
import (
  "errors"
  "fmt"

  "github.com/spf13/viper"
)

// Read in the configuration file. Will usually be specified as orca.toml, but
// viper supports JSON and YAML files, as well.
func New(configName string) (cfg *viper.Viper, err error) {
  cfg = viper.New()

  // The name of the configuration file, without extensions.
  cfg.SetConfigName(configName)
  // Set paths to look for the orca config file.
  cfg.AddConfigPath("/")
  cfg.AddConfigPath("$HOME/.orca")
  cfg.AddConfigPath(".")

  // Actually read in the configuration. Once we have read the configuration,
  // we can access members of the config using viper.Get()
  err = cfg.ReadInConfig()
  if err != nil {
    return
  }

  return
}

func HasRequired(cfg *viper.Viper, requiredFields []string) error {
  for _, v := range requiredFields {
    if set := cfg.IsSet(v); !set {
      return errors.New(fmt.Sprintf("Required field: %s", v))
    }
  }

  return nil
}
