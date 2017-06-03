package builder

import (
  "io"
  "os"

  "path/filepath"
  "text/template"
  
  log "github.com/gorobot/robologger"
)

func (b *Builder) addContextFiles() error {
  for _, file := range b.Config.Includes {
    absPath, _ := filepath.Abs(file)
    src, err := os.Open(absPath)
    defer src.Close()
    if err != nil {
      return err
    }

    newPath := filepath.Join(*b.Directory, file)
    dest, err := os.Create(newPath)
    defer dest.Close()
    if err != nil {
      return err
    }

    if _, err = io.Copy(dest, src); err != nil {
      return err
    }

    if err = dest.Sync(); err != nil {
      return err
    }

    log.Debugf("---> %s\n", newPath)
  }

  return nil
}

func (b *Builder) generateDockerfile() error {

  // Populate the struct to pass to the template.
  // data := struct {
  //   Base, Version, Hash, File, Mirror, URL string
  // }{
  //   Base: base,
  //   Version: version,
  //   Hash: hash,
  //   File: fn,
  //   Mirror: mirror,
  //   URL: uri,
  // }

  src, _ := filepath.Abs(b.Config.Dockerfile)

  path := filepath.Join(*b.Directory, "Dockerfile")
  dest, err := os.Create(path)
  if err != nil {
    return err
  }

  defer dest.Close()

  t, err := template.ParseFiles(src)
  if err != nil {
		return err
	}

  err = t.Execute(dest, b.Config)
  if err != nil {
    return err
  }

  log.Debugf("---> %s\n", path)

  return nil
}
