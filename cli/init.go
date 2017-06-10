package cli
//
// import (
//   "os"
//
//   "github.com/gorobot-library/orca/initialize"
//
//   log "github.com/gorobot/robologger"
//   "github.com/spf13/cobra"
// )
//
// var initCmd = &cobra.Command{
// 	Use:   "init",
// 	Short: "Generates initial Docker files.",
//   Long:  `Generates initial Docker files.`,
//   Run: func(cmd *cobra.Command, args []string) {
//     // // Check if the configuration file exists. If it does, we need to ask
//     // // whether or not the user wants to continue. They will receive a prompt
//     // // here, and when the program tries to overwrite the config file.
//     // if _, err := os.Stat("./orca.toml"); err == nil {
//     //   log.Warn("Configuration file found.")
//     //   log.Warn("Config file will be overwritten. This cannot be undone.")
//     //   res := log.Prompt(log.YESNO, "Do you want to continue?")
//     //
//     //   fres, _ := log.ParseResponse(res)
//     //   if fres != log.YES {
//     //     return
//     //   }
//     // }
//     //
//     // // Run through a series of prompts to get the necessary information about
//     // // the project.
//     // data := initialize.GetTemplateData()
//     //
//     // log.Info("Initializing...")
//     //
//     // // Create the configuration file.
//     // initialize.GenerateConfigFile(data)
//     //
//     // // Create the project files, such as the Dockerfile, any scripts, and an
//     // // entrypoint (if any).
//     // initialize.GenerateProjectFiles(data)
//
//     log.Info("Done.")
//   },
// }
