package client

import (
  "context"

  "github.com/docker/docker/api/types"
)

// ImageList wraps the Docker client.ImageList function and returns a string
// slice containing all of the image tags currently on the host system.
func (c *Client) ImageList(ctx context.Context) ([]string, error) {
  imageTags := []string{}

  // ImageList returns image structures, but we are interested only in the
  // RepoTags variable, which is a string slice. We append those tags to the
  // return variable.
	images, err := c.dockerClient.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return imageTags, err
	}

	for _, image := range images {
    for _, tag := range image.RepoTags {
      imageTags = append(imageTags, tag)
    }
	}

  return imageTags, nil
}
