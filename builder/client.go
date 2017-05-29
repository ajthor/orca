package builder

import (
  "context"

  "io/ioutil"

  log "github.com/gorobot-library/orca/logger"

  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
)

func createTempdir() string {
  // Create a temporary directory. Default is /tmp/orca.buildXXXXXXXXX
  dir, err := ioutil.TempDir("", "orca.build")
	if err != nil {
		log.Fatal(err)
	}

  log.Debugf("temp: %s", dir)

  return dir
}

func getServerVersion(c *client.Client) (ver string, err error) {
  v, err := c.ServerVersion(context.Background())
  ver = v.Version
  return
}

func getImages(c *client.Client) (imageTags []string, err error) {
	images, err := c.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return
	}

	for _, image := range images {
    for _, tag := range image.RepoTags {
      imageTags = append(imageTags, tag)
    }
	}

  return
}
