package client

import (
  "os"

  "io/ioutil"
  "path/filepath"

  // "github.com/gorobot-library/orca/manifest"

  log "github.com/gorobot/robologger"
)

type Context struct {

}

type Client struct {
  ClientOptions
}

type ClientOptions struct {
  // Points to the directory relative to the current directory where generated
  // files will be placed. Is a pointer to a string so that we can
  // differentiate between nil and an empty string.
  Directory *string
  // Holds the path to the directory where we copy files to create the build
  // context during `build` commands.
  ContextDirectory *string
}

func New(opts *ClientOptions) *Client {
  log.Info("Initializing Orca client.")

  if opts == nil {
    log.Fatal("Cannot initialize client.")
  }

  cli := &Client{
    *opts,
  }

  // Here, we check default values and update the client to handle if any
  // defaults are not set.
  if err := cli.checkDirectory(); err != nil {
    log.Fatal(err)
  }

  if cli.ContextDirectory == nil {
    cdir, err := cli.makeTempDirectory("context")
    if err != nil {
      log.Fatal(err)
    }
    
    cli.ContextDirectory = &cdir
  }

  return cli
}

// Close performs steps necessary to close the client.
func (c *Client) Close() error {
  defer os.RemoveAll(*c.ContextDirectory)

  return nil
}

// checkDirectory ensures that the Directory property passed to the client
// exists and is set. If it is not set, the directory is changed to the current
// working directory.
//
// Returns an error if the directory does not exist and cannot be created.
func (c *Client) checkDirectory() error {
  // If there is no directory supplied to the client, we need to create a
  // directory to hold all generatetd files. Default to the current directory.
  if c.Directory == nil {
    dir, _ := os.Getwd()
    c.Directory = &dir
  }

  // Get the absolute path of the directory. If the directory specified is not
  // an absolute directory, the `Abs` function joins the path with the current
  // working directory to turn it into an absolute path.
  absPath, _ := filepath.Abs(*c.Directory)

  // If the file does not exist, we create it, along with any parent
  // directories that are required.
  if _, err := os.Stat(*c.Directory); os.IsNotExist(err) {
    err := os.MkdirAll(*c.Directory, 0755)
    if err != nil { return err }
	}

  c.Directory = &absPath

  return nil
}

// makeDirectory creates a directory at the specified path with the proper
// permissions. This function is a helper function that also ensures that the
// permissions are acceptable to modify files in the directory. Default is
// 0755.
func (c *Client) makeDirectory(dir string) error {
  // Make sure file exists and has correct permissions.
  if stat, err := os.Stat(dir); os.IsExist(err) {
    if stat.Mode() < 0755 {
      if err := os.Chmod(dir, 0755); err != nil {
        return err
      }
    }
  } else {
    if err := os.MkdirAll(dir, 0755); err != nil {
      return err
    }
  }

  return nil
}

// makeTempDirectory creates a temporary directory and returns the directory
// name. Accepts a prefix, which will be added to the temporary directory name.
// Defaults to "/tmp/orcaXXXXXXXXX".
func (c *Client) makeTempDirectory(prefix string) (string, error) {
  if prefix == "" {
    prefix = "orca"
  }

  // Create a temporary directory to hold all downloaded binaries.
  dir, err := ioutil.TempDir("", prefix)

  log.Debugf("---> %s", dir)
  return dir, err
}
