package config

import (
  "errors"
  "fmt"
  "log"

  "github.com/spf13/viper"
)

type BuildInfo struct {
  base            string
  version         string
  externalFiles   []string
}

type RemoteInfo struct {
  url             string
  fileName        string
}

type ChecksumInfo struct {
  versions        []string
}

type Config struct {
  build           BuildInfo       `mapstructure:"build"`
  remote          RemoteInfo      `mapstructure:"remote"`
  checksum        ChecksumInfo    `mapstructure:"checksum"`
}

// Read in the configuration file. Will usually be specified as orca.toml, but
// viper supports JSON and YAML files, as well.
func NewConfig(configName string) *viper.Viper {
  // config := &Config{}

  cfg := viper.New()

  // The name of the configuration file, without extensions.
  cfg.SetConfigName(configName)
  // Set paths to look for the orca config file.
  cfg.AddConfigPath("/")
  cfg.AddConfigPath("$HOME/.orca")
  cfg.AddConfigPath(".")

  // Actually read in the configuration. Once we have read the configuration,
  // we can access members of the config using viper.Get()
  err := cfg.ReadInConfig()
  if err != nil {
    log.Fatal("Configuration file not found.")
  }

  // err := cfg.Unmarshal(&config)
  // if err != nil {
  //   fmt.Println("Could not unmarshal configuration.")
  // 	log.Fatal("Syntax error in configuration file.")
  // }

  return cfg
}

func HasRequired(cfg *viper.Viper, requiredFields []string) error {
  for _, v := range requiredFields {
    if set := cfg.IsSet(v); set != true {
      return errors.New(fmt.Sprintf("Required field: %s", v))
    }
  }

  return nil
}

func SetConfigWatch(cfg *viper.Viper) {

  // Call a function which handles configuration changes during the program
  // execution. Likely, this will not be an issue, but if a person leaves the
  // orca container running, the configuration file may change.
  // viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	fmt.Println("Config file changed:", e.Name)
	// })
}
