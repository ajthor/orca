package checksum

import (
  "bytes"
  "fmt"
  "log"
  "os"

  "html/template"
  "io/ioutil"
  "net/url"
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

type fileInfo struct {
  fn        string
  uri       string
  shasum    string
}

func newFileInfo(fn string) *fileInfo {
  f := &fileInfo{}
  return f
}

func Generate() {
  dir := createTempdir()

  // Remove the temporary directory when we are finished.
  defer os.RemoveAll(dir)

  // Get variables from the config file.
  download_filename := viper.Get("download_filename")
  download_url := viper.Get("download_url")
  versions := viper.Get("checksum.versions")

  f := newFileInfo()

  // Create a slice to hold the hashes output from the generateShasum function.
  files := make([]string, len(versions))

  // Iterate over the versions and generate and hashes for each file.
  for i, v := range versions {
    generateFilename(f, v)
    generateURI(f)

    fn := filepath.Join(dir, f.fn)

    downloadFile(f, fn)
    // The hash line is a combination of the hash, two spaces, and the filename.
    hashes[i] = generateShasum(fn) + "  " + fn
  }

  // Generate the shasums file.
  generateShasumFile(filepath.Join(dir, "SHASUMS256.txt"), hashes)

}
