package builder

import (
  "os"

  "path/filepath"
  "text/template"

  "github.com/gorobot-library/orca/checksum"

  "github.com/spf13/viper"
)

func generateDockerfile(cfg *viper.Viper, dir string) (path string, err error) {
  buildCfg := cfg.Sub("build")
  base := buildCfg.GetString("base")
  df := buildCfg.GetString("dockerfile")
  version := buildCfg.GetString("version")

  remoteCfg := cfg.Sub("remote")
  mirror := remoteCfg.GetString("mirror")

  fn := checksum.GetFilename(remoteCfg, version)
  uri := checksum.GetURI(remoteCfg, fn)

  hash := checksum.GetChecksum(fn)

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

  // absPath, err := filepath.Abs(df)
  // if err != nil {
  //   return
  // }
  // log.Debugf("Abs path: %s", absPath)
  // src, err := os.Open(absPath)
  // if err != nil {
  //   return
  // }
  //
  // defer src.Close()
  //
  // newPath := filepath.Join(dir, "Dockerfile")
  // log.Debugf("New path: %s", newPath)
  // dest, err := os.Create(newPath)
  // if err != nil {
  //   return
  // }
  //
  // defer dest.Close()
  //
  // if _, err = io.Copy(dest, src); err != nil {
  //   return
  // }
  //
  // if err = dest.Sync(); err != nil {
  //   return
  // }

  src, _ := filepath.Abs(df)
  // src, err := os.Open(absPath)
  // if err != nil {
  //   return
  // }
  //
  // defer src.Close()

  path = filepath.Join(dir, "Dockerfile")
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

func copyDockerfile(df string, dir string) (path string, dest *os.File, err error) {


  return
}
