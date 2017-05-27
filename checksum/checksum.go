package checksum

import (
  "log"
  "os"

  "io/ioutil"
  "path/filepath"

  "github.com/spf13/viper"
)

func createTempdir() string {
  // Create a temporary directory. Default is /tmp/checksumXXXXXXXXX
  dir, err := ioutil.TempDir("", "checksum")
	if err != nil {
		log.Fatal(err)
	}

  return dir
}

func GenerateChecksums(r *viper.Viper, versions []string) error {

  dir := createTempdir()

  // Remove the temporary directory when we are finished.
  defer os.RemoveAll(dir)

  cwd, err := os.Getwd()
  if err != nil {
      log.Fatal(err)
  }

  shasumFile := filepath.Join(cwd, "SHASUMS256.txt")

  if err := removeShasumFile(shasumFile); err != nil {
    log.Fatal(err)
  }

  // Create a slice to hold the hashes output from the generateShasum function.
  hashes := make([]string, len(versions))

  // Iterate over the versions and generate and hashes for each file.
  for i, ver := range versions {
    fn := getFilename(r, ver)
    uri := getURI(r, fn)

    dlFile := filepath.Join(dir, fn)

    err := downloadFile(uri, dlFile)
    if err != nil {
      log.Fatal(err)
    }

    defer os.Remove(dlFile)

    // Generate the shasums.
    hashes[i] = generateShasum(dlFile)
    if err != nil {
      log.Fatal(err)
    }
  }

  // Write the hashes to the shasum file.
  err = createShasumFile(shasumFile, hashes)
  if err != nil {
    log.Fatal(err)
  }

  return nil
}
