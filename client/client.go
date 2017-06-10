package client

import (
  "errors"
  "os"

  "io/ioutil"
  "path/filepath"

  dockerClient "github.com/docker/docker/client"
  log "github.com/gorobot/robologger"
)

var (
  ErrClientOptionsNotDefined = errors.New("ClientOptions is not defined.")
)

type Client struct {
  ClientOptions

  // Docker client. The Orca client relies heavily on the Docker client.
  dockerClient *dockerClient.Client
}

type ClientOptions struct {
  // Points to the directory relative to the current directory where generated
  // files will be placed. Must be *string so that we can differentiate between
  // nil and an empty string.
  Directory *string
}

// NewClient creates a new Orca client. It takes a Docker client and a
// ClientOptions struct as arguments. The Docker client is required for any
// functions that interact with Docker directly.
func NewClient(dockerClient *dockerClient.Client, opts *ClientOptions) *Client {
  if opts == nil {
    log.Fatal(ErrClientOptionsNotDefined)
  }

  cli := &Client{
    ClientOptions: *opts,
    dockerClient: dockerClient,
  }

  // Here, we check default values and update the client to handle if any
  // defaults are not set.
  if err := cli.checkDockerClient(); err != nil {
    log.Fatal(err)
  }

  if err := cli.checkDirectory(); err != nil {
    log.Fatal(err)
  }

  return cli
}

// Close performs steps necessary to close the client.
func (c *Client) Close() error {
  return nil
}

func (c *Client) checkDockerClient() error {
  if c.dockerClient != nil {
    return nil
  }

  // Set up the Docker client. We use default, expected values in order to
  // set up the client.
  host := "unix:///var/run/docker.sock"
  version := "v1.29"

  headers := map[string]string{
    "User-Agent": "engine-api-cli-1.0",
  }

  // Create a new connection to the Docker server.
  cli, err := dockerClient.NewClient(host, version, nil, headers)
  if err != nil {
    return err
  }

  c.dockerClient = cli

  return nil
}

// checkDirectory ensures that the Directory property passed to the client
// exists and is set. If it is not set, the directory is changed to the current
// working directory.
//
// Returns an error if the directory does not exist and cannot be created.
func (c *Client) checkDirectory() error {
  // If the file does not exist, we create it, along with any parent
  // directories that are required.
  // err := mkdir(*c.Directory)
  // if err != nil {
  //   return err
  // }

  if c.Directory != nil {
    return nil
  }
  // If there is no directory supplied to the client, we default to the current
  // working directory.

  dir, err := os.Getwd()
  if err != nil {
    return err
  }

  c.Directory = &dir

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
