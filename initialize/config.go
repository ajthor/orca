package initialize

import (
  "os"

  "path/filepath"

  log "github.com/gorobot-library/orca/logger"
)

func GenerateConfigFile() {
  name := getName()
  log.ShowInput(name)

  base := getBase()
  log.ShowInput(base)

  version := getVersion()
  log.ShowInput(version)

  data := &TemplateData{
    Name: name,
    Base: base,
    Version: version,
  }

  err := templateFile("./initialize/templates/orca.toml", "orca.toml", data)
  if err != nil {
    log.Fatal(err)
  }
}

func getName() string {
  cwd, err := os.Getwd()
  if err != nil {
      log.Fatal(err)
  }

  cwdBase := filepath.Base(cwd)

  name := log.Promptf(log.DEFAULT, "Name (%s):", cwdBase)
  if name == "" {
    name = cwdBase
  }

  return name
}

func getBase() string {
  defaultBase := "scratch"
  base := log.Promptf(log.DEFAULT, "Base (%s):", defaultBase)
  if base == "" {
    base = defaultBase
  }

  return base
}

func getVersion() string {
  defaultVersion := "0.0.1"
  version := log.Promptf(log.DEFAULT, "Version (%s):", defaultVersion)
  if version == "" {
    version = defaultVersion
  }

  return version
}
