package checksum

import (
  "os"

  "io/ioutil"
  "path/filepath"

  log "github.com/gorobot-library/orca/logger"

  "github.com/spf13/viper"
)

func createTempdir() string {
  // Create a temporary directory. Default is /tmp/checksumXXXXXXXXX
  dir, err := ioutil.TempDir("", "checksum")
	if err != nil {
		log.Fatal(err)
	}
  
  log.Debugf("Created temp directory: %s", dir)

  return dir
}

func GenerateChecksums(cfg *viper.Viper, versions []string) error {
  log.Info("Generating shasums...")

  dir := createTempdir()

  // Remove the temporary directory when we are finished.
  defer os.RemoveAll(dir)

  cwd, err := os.Getwd()
  if err != nil {
      log.Fatal(err)
  }

  shasumFile := filepath.Join(cwd, "SHASUMS256.txt")
  log.Debugf("Shasum file path: %s", shasumFile)

  if err := removeShasumFile(shasumFile); err != nil {
    log.Fatal(err)
  }

  // Create a slice to hold the hashes output from the generateShasum function.
  hashes := make([]string, len(versions))

  // Iterate over the versions and generate and hashes for each file.
  for i, ver := range versions {
    fn := getFilename(cfg, ver)
    log.Debugf("Generated filename: %s", fn)
    uri := getURI(cfg, fn)
    log.Debugf("Generated uri: %s", uri)

    dlFile := filepath.Join(dir, fn)

    log.Infof("Downloading %s...", dlFile)
    err := downloadFile(uri, dlFile)
    if ok := log.Done(err); !ok {
      log.Fatal(err)
    }

    defer os.Remove(dlFile)

    // Generate the shasums.
    log.Info("Generating shasum...")
    hashes[i] = generateShasum(dlFile)
    if ok := log.Done(err); !ok {
      log.Fatal(err)
    }
  }

  // Write the hashes to the shasum file.
  log.Info("Generating shasum file...")
  err = createShasumFile(shasumFile, hashes)
  if ok := log.Done(err); !ok {
    log.Fatal(err)
  }

  return nil
}
