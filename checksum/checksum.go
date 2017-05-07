package checksum

import (
  "bytes"
  "fmt"
  "log"
  "os"

  "html/template"
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

func newFileInfo(download_filename string, download_url string) *fileInfo {
  f := &fileInfo{
    fn:     download_filename,
    uri:    download_url,
  }
  return f
}

func Generate() {
  dir := createTempdir()

  // Remove the temporary directory when we are finished.
  defer os.RemoveAll(dir)

  // Get variables from the config file.
  download_filename := viper.Get("download.filename")
  download_url := viper.Get("download.url")
  checksum_versions := viper.Get("checksum.versions")

  f := newFileInfo(download_filename, download_url)

  // Create a slice to hold the hashes output from the generateShasum function.
  hashes := make([]string, len(checksum_versions))

  // Iterate over the checksum_versions and generate and hashes for each file.
  for i, v := range checksum_versions {
    f.fn = generateFilename(f.fn, v)
    f.uri = generateURI(f.fn, f.uri)

    fn := filepath.Join(dir, f.fn)

    downloadFile(f.fn, f.uri)
    // The hash line is a combination of the hash, two spaces, and the filename.
    hashes[i] = generateShasum(f.fn)
  }

  // Generate the shasums file.
  generateShasumFile(filepath.Join(dir, "SHASUMS256.txt"), hashes)

}
