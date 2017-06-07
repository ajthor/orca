package client

import (
  "errors"
  "os"

  "io/ioutil"
  "path/filepath"

  client "github.com/docker/docker/client"
  log "github.com/gorobot/robologger"
)

var (
  ErrClientOptionsNotDefined = errors.New("ClientOptions is not defined.")
)

type Client struct {
  ClientOptions

  // Docker client. The Orca client relies heavily on the Docker client.
  client *client.Client
}

type ClientOptions struct {
  // Points to the directory relative to the current directory where generated
  // files will be placed. Must be *string so that we can differentiate between
  // nil and an empty string.
  Directory *string
}

func NewClient(client *client.Client, opts *ClientOptions) *Client {
  log.Info("Initializing Orca client.")
  if opts == nil {
    log.Fatal(ErrClientOptionsNotDefined)
  }

  cli := &Client{
    ClientOptions: *opts,
    client: client,
  }

  // Here, we check default values and update the client to handle if any
  // defaults are not set.
  if err := cli.checkDirectory(); err != nil {
    log.Fatal(err)
  }

  return cli
}

// Close performs steps necessary to close the client.
func (c *Client) Close() error {
  return nil
}

// checkDirectory ensures that the Directory property passed to the client
// exists and is set. If it is not set, the directory is changed to the current
// working directory.
//
// Returns an error if the directory does not exist and cannot be created.
func (c *Client) checkDirectory() error {
  // If there is no directory supplied to the client, we default to the current
  // working directory.
  if c.Directory == nil {
    // If the file does not exist, we create it, along with any parent
    // directories that are required.
    dir, err := os.Getwd()
    if err != nil {
      return err
    }

    c.Directory = &dir
  }

  err := mkdir(*c.Directory)
  if err != nil {
    return err
  }

  return nil
}

// mkdir is a helper function that creates a directory at the specified path
// with the proper permissions. This function also ensures that the permissions
// are acceptable to modify files in the directory. Default is 0755.
func mkdir(path string) error {
  // Get the absolute path of the directory. If the directory specified is not
  // an absolute directory, the `Abs` function joins the path with the current
  // working directory to turn it into an absolute path.
  absPath, _ := filepath.Abs(path)

  // Make sure file exists and has correct permissions.
  if stat, err := os.Stat(absPath); os.IsExist(err) {
    if stat.Mode() < 0755 {
      if err := os.Chmod(absPath, 0755); err != nil {
        return err
      }
    }
  } else {
    if err := os.MkdirAll(absPath, 0755); err != nil {
      return err
    }
  }

  return nil
}

// tempdir is a helper function that creates a temporary directory and returns
// the directory name. Accepts a prefix, which will be added to the temporary
// directory name.
//
// Defaults to "/tmp/orcaXXXXXXXXX".
func tempdir(dir, prefix string) string {
  if prefix == "" {
    prefix = "orca"
  }

  // Create a temporary directory to hold all downloaded binaries.
  dir, err := ioutil.TempDir("", prefix)
  if err != nil {
    panic(err)
  }

  log.Debugf("---> %s", dir)
  return dir
}

// mustOpen is a helper function that wraps an os.Open function call and panics
// if the error is non-nil.
func mustOpen(name string) *os.File {
  file, err := os.Open(name)
  if err != nil {
    panic(err)
  }

  return file
}

// mustOpen is a helper function that wraps an os.Open function call and panics
// if the error is non-nil.
func mustCreate(name string) *os.File {
  file, err := os.Create(name)
  if err != nil {
    panic(err)
  }

  return file
}
