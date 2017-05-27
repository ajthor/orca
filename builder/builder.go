package builder

// "context"
// "io/ioutil"
// "os"
//
// "github.com/moby/moby/api/types"

import (
  "fmt"

	"github.com/moby/moby/client"
)

func newClient() *client.Client {
  client, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

  return client
}

func Build()  {
  fmt.Println("Building image...")
  // client := newClient()

  // client.ImageBuild(context.Background(), , types.ImageBuildOptions{})

  // func (cli *Client) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
}

// func getImages(client *client.Client)  {
//   containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{})
//   if err != nil {
//   	panic(err)
//   }
//
//   for _, container := range containers {
//   	fmt.Printf("%s %s\n", container.ID[:10], container.Image)
//   }
// }
