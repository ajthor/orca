package checksum

import (
  "bufio"
  "os"
  "strings"

  "io/ioutil"
  "path/filepath"

  log "github.com/gorobot-library/orca/logger"

  "github.com/spf13/viper"
)

func createTempdir() string {
  // Create a temporary directory. Default is /tmp/orca.checksumXXXXXXXXX
  dir, err := ioutil.TempDir("", "checksum")
	if err != nil {
		log.Fatal(err)
	}

  log.Debugf("---> %s", dir)

  return dir
}

func GenerateChecksums(cfg *viper.Viper, versions []string) {
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
    fn := GetFilename(cfg, ver)
    log.Debugf("Generated filename: %s", fn)
    uri := GetURI(cfg, fn)
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

}

func GetChecksum(file string) (shasum string) {
  cwd, err := os.Getwd()
  if err != nil {
      log.Fatal(err)
  }

  shasumFilePath := filepath.Join(cwd, "SHASUMS256.txt")

  if _, err := os.Stat(shasumFilePath); os.IsNotExist(err) {
    return ""
  }

  sf, err := os.Open(shasumFilePath)
  if err != nil {
      log.Fatal(err)
  }
  defer sf.Close()

  scanner := bufio.NewScanner(sf)
  for scanner.Scan() {
    if contains := strings.Contains(scanner.Text(), file); contains {
      shasum = scanner.Text()
    }
  }

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  return
}
