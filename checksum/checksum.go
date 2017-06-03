package checksum

import (
  "os"

  "io/ioutil"
  "path/filepath"

  "github.com/gorobot-library/orca/config"

  log "github.com/gorobot/robologger"
)

var ShasumDefaultFile string = "SHASUMS256.txt"

type checksum interface {
  CreateShasumFile()
  GenerateHashes()
  GetShasum()
  GenerateShasum()
}

type Checksum struct {
  Config *ChecksumOptions
  // This holds the name of the temporary directory we use to download the
  // binaries into. Binaries and the temporary directory will be deleted after
  // they are used.
  Directory *string
  // ShaFile holds the name of the shasum file we are generating. By default,
  // this will be set as "SHASUMS256.txt".
  ShaFile *string
}

type ChecksumOptions struct {
  Mirror string
  File string
  Versions []string
}

func New(file string) *Checksum {
  if file == "" {
    ex, _ := os.Executable()
    dir := filepath.Dir(ex)

    file = filepath.Join(dir, ShasumDefaultFile)
  }

  return &Checksum{
    ShaFile: &file,
  }
}

func (c *Checksum) GenerateHashes(opts *ChecksumOptions) ([]string, error) {

  c.Config = opts

  log.Info("Generating shasums...")

  // If there is no directory supplied to the checksum generator, we need to
  // create a temporary directory to hold all downloaded binaries.
  if c.Directory == nil {
    // Create a temporary directory. Default is /tmp/orca.checksumXXXXXXXXX
    dir, err := ioutil.TempDir("", "checksum")
    if err != nil {
      log.Fatal(err)
    }

    log.Debugf("---> %s", dir)
    c.Directory = &dir

    // Remove the temporary directory when we are finished.
    defer os.RemoveAll(dir)
  }

  // Create a slice to hold the hashes output from the generateShasum function.
  hashes := make([]string, len(c.Config.Versions))

  // Iterate over the versions and generate and hashes for each file.
  for i, ver := range c.Config.Versions {
    // Download the file.
    file, err := config.ParseFilename(c.Config.File, ver)
    if err != nil {
      return hashes, err
    }

    uri, err := config.ParseURL(c.Config.Mirror, file)
    if err != nil {
      return hashes, err
    }

    log.Infof("Downloading %s...", uri)

    dlFile := filepath.Join(*c.Directory, file)

    err = c.downloadFile(uri, dlFile)
    if err != nil {
      return hashes, err
    }

    log.Debugf("---> %s", dlFile)

    defer os.Remove(dlFile)

    // Generate the shasum.
    log.Info("Generating hash...")
    hashes[i] = c.GenerateShasum(dlFile)
    if err != nil {
      return hashes, err
    }
  }

  return hashes, nil
}
