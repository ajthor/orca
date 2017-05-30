package builder

import (
  "context"
  "io"
  "os"

  "archive/tar"
  "compress/gzip"
  "io/ioutil"
  "path/filepath"

  log "github.com/gorobot-library/orca/logger"

  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/spf13/viper"
)

func Build(cfg *viper.Viper) {
  // "unix:///var/run/docker.sock"
  c, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

  log.Info("Creating build context...")
  buildContext, err := createBuildContext(cfg)
	if err != nil {
		log.Fatal(err)
	}
  // Make sure to close the tarfile when we're done.
  defer buildContext.Close()

  // Create the build options.
  log.Info("Getting build information...")
  buildOptions, err := createBuildOptions(cfg)
  if ok := log.Done(err); !ok {
    log.Fatal(err)
  }

  // buildOptions.Context = buildContext

  log.Info("Building image...")
  resp, err := c.ImageBuild(context.Background(), buildContext, buildOptions)
  if ok := log.Done(err); !ok {
    log.Info(buildOptions)
    log.Info(buildContext)
    log.Info(resp)
    log.Fatal(err)
  }
  defer resp.Body.Close()

  return
}

func createBuildContext(cfg *viper.Viper) (tarFile *os.File, err error) {
  // We create a temporary directory for the build context.
  dir := createTempdir()

  includes := cfg.GetStringSlice("includes")

  // Create the tar file.
  tarFile, err = ioutil.TempFile(dir, "build.tar.gz")
  if err != nil {
    return
  }

  log.Debugf("chmod: %s [0644]", tarFile.Name())
  if err := os.Chmod(tarFile.Name(), 0644); err != nil {
		log.Fatal(err)
	}

  gw := gzip.NewWriter(tarFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

  // Add the dockerfile to the includes.
  // includes = append(includes, df)
  fp, err := generateDockerfile(cfg, dir)
  log.Infof("---> %s\n", fp)
  err = addTarFile(tw, fp)
  if err != nil {
    return tarFile, err
  }

  // Copy in the include files.
  for _, file := range includes {
    fp, err := copyFile(file, dir)
    if err != nil {
      return tarFile, err
    }

    log.Infof("---> %s\n", fp)
    err = addTarFile(tw, fp)
    if err != nil {
      return tarFile, err
    }
  }

  return
}

func copyFile(f string, dir string) (newPath string, err error) {
  absPath, _ := filepath.Abs(f)
  src, err := os.Open(absPath)
  if err != nil {
    return
  }

  defer src.Close()

  newPath = filepath.Join(dir, f)
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

func addTarFile(tw *tar.Writer, file string) (err error) {
  out, err := os.Open(file)
	if err != nil {
		return
	}
	defer out.Close()

  stat, err := out.Stat()
  if err != nil {
    return
  }

	hdr := &tar.Header{
		Name: filepath.Base(file),
	  Size: stat.Size(),
		Mode: int64(stat.Mode()),
	}

  err = tw.WriteHeader(hdr)
	if err != nil {
		return
	}

  _, err = io.Copy(tw, out)
  if err != nil {
		return
	}

  return
}

func createBuildOptions(cfg *viper.Viper) (buildOptions types.ImageBuildOptions, err error) {
  // Populate the build options from the config.
  buildOptions.Tags = cfg.GetStringSlice("tags")
  buildOptions.Dockerfile = "Dockerfile"

  return
}
