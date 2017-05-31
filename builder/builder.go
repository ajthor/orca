package builder

// "archive/tar"
// "io/ioutil"
// "compress/gzip"
//
// "github.com/docker/docker/pkg/progress"
// "github.com/docker/docker/pkg/streamformatter"
// "github.com/docker/docker/pkg/term"
// "github.com/docker/libcompose/logger"
import (
  "context"
  "io"
  "os"

  "path/filepath"

  log "github.com/gorobot-library/orca/logger"

  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"

  "github.com/spf13/viper"
)

func Build(cfg *viper.Viper) {
  // "unix:///var/run/docker.sock"
  defaultHeaders := map[string]string{
    "User-Agent": "engine-api-cli-1.0",
  }

  c, err := client.NewClient("unix:///var/run/docker.sock", "v1.29", nil, defaultHeaders)
	if err != nil {
		log.Fatal(err)
	}

  log.Info("Creating build context...")
  // We create a temporary directory for the build context.
  dir := createTempdir()
  buildContext, err := createBuildContext(cfg, dir)
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

  log.Info("Building image...")
  resp, err := c.ImageBuild(context.Background(), buildContext, buildOptions)
  if err != nil {
		log.Fatal(err)
	}
  defer resp.Body.Close()

  outFd, isTerminalOut := term.GetFdInfo(os.Stdout)
  err = jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stdout, outFd, isTerminalOut, nil)
	if err != nil {
		// if jerr, ok := err.(*jsonmessage.JSONError); ok {
    //   log.Debugf("Status: %s, Code: %d", jerr.Message, jerr.Code)
    //   log.Debug(jerr.Error())
		// 	// If no error code is set, default to 1
		// 	// if jerr.Code == 0 {
		// 	// 	jerr.Code = 1
		// 	// }
		// 	// errBuff.Write([]byte(jerr.Error()))
		// 	// return fmt.Errorf("Status: %s, Code: %d", jerr.Message, jerr.Code)
		// }
    log.Fatal(err)
	}

  imgs, err := getImages(c)
  for _, img := range imgs {
    log.Debugf(" --- %s", img)
  }

}

func createBuildContext(cfg *viper.Viper, dir string) (ctx io.ReadCloser, err error) {
  includes := cfg.GetStringSlice("includes")

  // tarIncludes := make([]string, len(includes) + 1)

  // Create the tar file.
  // tarFile, err = ioutil.TempFile(dir, "tar")
  // if err != nil {
  //   return
  // }
  //
  // log.Debugf("chmod: %s [0644]", tarFile.Name())
  // if err := os.Chmod(tarFile.Name(), 0644); err != nil {
	// 	log.Fatal(err)
	// }

  // gw := gzip.NewWriter(tarFile)
	// defer gw.Close()
  // tw := tar.NewWriter(gw)
  // tw := tar.NewWriter(tarFile)
	// defer tw.Close()

  // Add the dockerfile to the includes.
  // includes = append(includes, df)
  fp, err := generateDockerfile(cfg, dir)
  if err != nil {
    return
  }
  // fp, err := generateDockerfile(cfg, dir)
  log.Debugf("---> %s\n", fp)
  // err = addTarFile(tw, fp)
  // if err != nil {
  //   return tarFile, err
  // }

  // Copy in the include files.
  for _, file := range includes {
    fp, _ := copyFile(file, dir)
    // if err != nil {
    //   return
    // }

    log.Debugf("---> %s\n", fp)
    // err = addTarFile(tw, fp)
    // if err != nil {
    //   return tarFile, err
    // }
  }

  // options := &archive.TarOptions{
	// 	Compression:     archive.Uncompressed,
	// 	IncludeFiles:    includes,
	// }

	ctx, err = archive.Tar(dir, archive.Gzip)

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

// func addTarFile(tw *tar.Writer, file string) (err error) {
//   out, err := os.Open(file)
// 	if err != nil {
// 		return
// 	}
// 	defer out.Close()
//
//   stat, err := out.Stat()
//   if err != nil {
//     return
//   }
//
// 	hdr := &tar.Header{
// 		Name: filepath.Base(file),
// 	  // Size: stat.Size(),
// 		// Mode: int64(stat.Mode()),
// 	}
//
//   err = tw.WriteHeader(hdr)
// 	if err != nil {
// 		return
// 	}
//
//   _, err = io.Copy(tw, out)
//   if err != nil {
// 		return
// 	}
//
//   return
// }

func createBuildOptions(cfg *viper.Viper) (buildOptions types.ImageBuildOptions, err error) {
  // Populate the build options from the config.
  buildOptions.Tags = cfg.GetStringSlice("build.tags")
  // buildOptions.Tags = []string{"test"}
  buildOptions.Dockerfile = "Dockerfile"

  return
}
