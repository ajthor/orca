package initialize

import (
  "os"

  "path/filepath"
  "text/template"

  log "github.com/gorobot-library/orca/logger"

  "github.com/spf13/viper"
)

func Initialize(cfg *viper.Viper) {
  log.Info("Initializing...")

  name := cfg.GetString("name")
  base := cfg.GetString("build.base")
  version := cfg.GetString("build.version")

  data := &TemplateData{
    Name: name,
    Base: base,
    Version: version,
  }

  err := templateFile("./initialize/templates/Dockerfile", "Dockerfile", data)
  if err != nil {
    log.Fatal(err)
  }
  // log.Prompt(log.YESNO, "Does your Dockerfile use a base language?")
}

type TemplateData struct {
  Name, Base, Version string
}

func templateFile(srcPath, destPath string, data *TemplateData) (err error) {
  src, _ := filepath.Abs(srcPath)

  dir, _ := os.Getwd()
  // relPath, err := filepath.Rel("template", )

  path := filepath.Join(dir, destPath)
  log.Infof("---> %s", path)

  if _, err = os.Stat(path); err == nil {
    res := log.Prompt(log.YESNO, "File already exists. Overwrite?")
    log.ShowInput(res)

    fres := log.FormatResponse(res)
    if fres != log.YES {
      return
    }
  } else {
    return
  }

  dest, err := os.Create(path)
  if err != nil {
    return
  }

  defer dest.Close()

  t, err := template.ParseFiles(src)
  if err != nil {
    return
  }

  err = t.Execute(dest, data)
  if err != nil {
    return
  }

  return
}
