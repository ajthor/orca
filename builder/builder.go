package builder

// "io/ioutil"
// "os"
//
// "context"
// "github.com/docker/docker/api/types"

import (
  log "github.com/gorobot-library/orca/logger"

  "github.com/docker/docker/client"

  "github.com/spf13/viper"
)

func Build(cfg *viper.Viper) (err error) {
  log.Info("Building image...\n")

  c, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

  imageTags, err := getImages(c)

  for _, tag := range imageTags {
    log.Debugf("%s\n", tag)
  }

  return
  // client.ImageBuild(context.Background(), , types.ImageBuildOptions{})

  // func (cli *Client) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
}
