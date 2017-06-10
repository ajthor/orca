package client

import (
  "context"
  "os"

  // "github.com/gorobot-library/orca/manifest"

  "github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

// ImageBuild wraps the docker client.ImageBuild function. The Orca function
// automatically creates the build context based on the manifest and outputs
// the build progress to stdout.
//
// It is important to set the Tags variable inside the ImageBuildOptions passed
// to the function. The function does not generate them for you, but there is a
// method in the `manifest` package to assist in generating the tags.
func (c *Client) ImageBuild(ctx context.Context, buildContext *Context, options types.ImageBuildOptions) error {
  // Create the build context.
  tarFile, err := buildContext.Tar()
	if err != nil {
		return err
	}

  // Make sure to close the tarfile when we're done.
  defer tarFile.Close()

  // Ensure that NoCache is set to true in the options.
  options.NoCache = true

  // Once we have the build context, we can go through the process of building
  // the image. This is handles by the `ImageBuild` function.
  resp, err := c.dockerClient.ImageBuild(ctx, tarFile, options)
  if err != nil {
		return err
	}

  defer resp.Body.Close()

  // Print the Docker build info to the stream.
  outFd, isTerminalOut := term.GetFdInfo(os.Stdout)
  err = jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stdout, outFd, isTerminalOut, nil)
	if err != nil {
    return err
	}

  return nil
}
