package builder

import (
  "io"
  "os"

  "path/filepath"
  "text/template"

  "github.com/gorobot-library/orca/checksum"
  log "github.com/gorobot-library/orca/logger"

  "github.com/spf13/viper"
)

func generateDockerfile(cfg *viper.Viper, dir string) (fp string, err error) {
  df := cfg.GetString("dockerfile")
  // Copy Dockerfile.
  fp, err = copyDockerfile(df, dir)
  if err != nil {
    return
  }

  // Template the new Dockerfile.
  log.Info("Templating Dockerfile...")
  err = templateDockerfile(cfg, fp)
  if ok := log.Done(err); !ok {
    return
  }

  return
}

func templateDockerfile(cfg *viper.Viper, path string) (err error) {
  base := cfg.GetString("base")
  mirror := cfg.GetString("mirror")
  version := cfg.GetString("version")

  fn := checksum.GetFilename(cfg, version)
  uri := checksum.GetURI(cfg, fn)

  hash := checksum.GetChecksum(fn)
  log.Infof("Hash: %s", hash)

  data := struct {
    Base, Version, Hash, File, Mirror, URL string
  }{
    Base: base,
    Version: version,
    Hash: hash,
    File: fn,
    Mirror: mirror,
    URL: uri,
  }

  file, err := os.Open(path)
  if err != nil {
    return
  }

  defer file.Close()

  t, err := template.ParseFiles(path)
  if err != nil {
		return
	}

  err = t.Execute(file, data)
  if err != nil {
    return
  }

  return
}

func copyDockerfile(f string, dir string) (newPath string, err error) {
  absPath, _ := filepath.Abs(f)
  src, err := os.Open(absPath)
  if err != nil {
    return
  }

  defer src.Close()

  newPath = filepath.Join(dir, "Dockerfile")
  dest, err := os.Create(newPath)
  if err != nil {
    return
  }

  defer dest.Close()

  if _, err = io.Copy(dest, src); err != nil {
    return
  }

  if err = dest.Sync(); err != nil {
    return
  }

  return
}
