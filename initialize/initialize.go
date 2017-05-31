package initialize

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go templates/...

import (
  "os"

  "path/filepath"
  "text/template"

  log "github.com/gorobot-library/orca/logger"
)

func GenerateConfigFile(data *TemplateData) {
  log.Info("Creating orca.toml...")

  err := templateFile("./initialize/templates/orca.toml.tmpl", "orca.toml", data)
  if err != nil {
    log.Fatal(err)
  }
}

func GenerateProjectFiles(data *TemplateData) {
  var tpl string

  log.Info("Creating project files...")

  tpl = "templates/Dockerfile.tmpl"
  err := templateFile(tpl, "Dockerfile", data)
  if err != nil {
    log.Fatal(err)
  }

  if data.HasEntrypoint {
    tpl = "./initialize/templates/docker-entrypoint.sh.tmpl"
    err := templateFile(tpl, "docker-entrypoint.sh", data)
    if err != nil {
      log.Fatal(err)
    }
  }
}

type TemplateData struct {
  Name, Base, Version, Tag, Mirror, File string
  HasEntrypoint bool
}

func GetTemplateData() *TemplateData {
  log.Info("Let's get some basic information about your project:")

  name := getName()
  log.ShowInput(name)

  base := getBase()
  log.ShowInput(base)

  version := getVersion()
  log.ShowInput(version)

  tag := getTag()
  log.ShowInput(tag)

  mirror := getRemoteMirror()
  log.ShowInput(mirror)

  file := getRemoteFile()
  log.ShowInput(file)

  hasEntrypoint := getHasEntrypoint()

  // res := log.Prompt(log.YESNO, "Does your Dockerfile use a base language?")

  return &TemplateData{
    Name: name,
    Base: base,
    Version: version,
    Tag: tag,
    Mirror: mirror,
    File: file,
    HasEntrypoint: hasEntrypoint,
  }
}

func getName() string {
  defaultName := "project"
  name := log.Promptf(log.DEFAULT, "Project name (%s):", defaultName)
  if name == "" {
    name = defaultName
  }

  return name
}

func getBase() string {
  defaultBase := "scratch"
  base := log.Promptf(log.DEFAULT, "Base image (%s):", defaultBase)
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

func getTag() string {
  return log.Prompt(log.DEFAULT, "Docker tag (e.g. repository/image:tag):")
}

func getRemoteMirror() string {
  return log.Prompt(log.DEFAULT, "Remote mirror (e.g. http://github.com/download/):")
}

func getRemoteFile() string {
  return log.Promptf(log.DEFAULT, "Remote file (e.g. sample-{{.Version}}.tar.gz):")
}

func getHasEntrypoint() bool {
  res := log.Prompt(log.YESNO, "Does your Dockerfile have an entrypoint?")
  log.ShowInput(res)

  fres := log.FormatResponse(res)
  if fres == log.YES {
    return true
  }

  return false
}

func templateFile(srcPath, destPath string, data interface{}) (err error) {
  // src, _ := filepath.Abs(srcPath)
  assetData, err := Asset(srcPath)
  if err != nil {
    return
  }

  ex, err := os.Executable()
  if err != nil {
    return
  }
  dir := filepath.Dir(ex)
  // relPath, err := filepath.Rel("template", )

  path := filepath.Join(dir, destPath)
  log.Infof("---> %s", path)

  if _, err = os.Stat(path); err == nil {
    res := log.Promptf(log.YESNO, "%s already exists. Overwrite?", path)
    log.ShowInput(res)

    fres := log.FormatResponse(res)
    if fres != log.YES {
      return
    }
  }

  dest, err := os.Create(path)
  if err != nil {
    return
  }

  defer dest.Close()

  t, err := template.New("").Parse(string(assetData))
  if err != nil {
    return
  }

  err = t.Execute(dest, data)
  if err != nil {
    return
  }

  return
}
