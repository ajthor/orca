package main

import (
  "fmt"
  "os"

  "github.com/gorobot-library/orca/cli"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

func newCLICommand() *cobra.Command {
  // Initialize the "root command".
  // By default, simply running the `orca` command will not perform an action.
  cmd := &cobra.Command{
      Use:   "orca",
      Short: "Orca is a simple Docker image build tool.",
      Long:  `Orca is a simple Docker image build tool.`,
    }

  // Perform initialization steps, such as attaching commands to the root
  // command. Commands and initialization are handled by the cli/cli.go file.
  cli.SetupCLIRootCmd(cmd)

  return cmd
}

// Read in the configuration file. Will usually be specified as orca.yaml, but
// viper supports TOML and JSON files, as well.
func readConfig() {
  // The name of the configuration file, without extensions.
  viper.SetConfigName("orca")
  // Set paths to look for the orca config file.
  viper.AddConfigPath("/")
  viper.AddConfigPath("$HOME/.orca")
  viper.AddConfigPath(".")

  // Actually read in the configuration.
  err := viper.ReadInConfig()
  // Once we have read the configuration, we can access members of the config
  // using viper.Get()
  if err != nil {
  	panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }

  viper.WatchConfig()
  // Call a function which handles configuration changes during the program
  // execution. Likely, this will not be an issue, but if a person leaves the
  // orca container running, the configuration file may change.
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

func main() {
  // Read the configuration file.
  if err := readConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
  }

  cmd := newCLICommand()
  // Run the command.
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
