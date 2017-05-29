package builder

import (
  "context"
  "fmt"
  "os"

  "path/filepath"

  log "github.com/gorobot-library/orca/logger"

  "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

  "github.com/Masterminds/semver"

  "github.com/spf13/viper"
)

func Validate(cfg *viper.Viper) (isValid bool, err error) {
  // fmt.Print(log.Std)
  log.Info("Checking configuration...")

  c, err := client.NewEnvClient()
	if err != nil {
		log.Debug("\n")
    log.Fatal(err)
	}

  log.Info("Has compatible docker version...")
  err = checkDockerVersion(c)
  if ok := log.Done(err); !ok {
    log.Fatal(err)
  }

  // if err != nil {
  //   log.Debug("\n")
  //   log.Fatal(err)
  // } else if !isValidVersion {
  //   isValid = false
  //   log.Modifyf(-1, log.ERROR, "Has compatible docker version... %s\n", "no")
  //   log.Fatal("Orca needs Docker >= 17.05 to run.\n")
  //   return
  // }
  //
  // log.Modifyf(-1, log.INFO, "Has compatible docker version... %s\n", "ok")

  log.Info("Has base image...")
  err = checkBaseImageExists(cfg, c)
  if ok := log.Done(err); !ok {
    closestMatch, _ := getClosestMatch(cfg, c)
    log.Warnf("Did you mean: %s\n", closestMatch)
    //
    // NOTE: Prompt the user for confirmation. If yes, use the matched image.
    // NOTE: Perhaps try and download the image from Docker Hub?
    //
    log.Error(err)
  }

  // if err != nil {
  //   log.Modify(-1, log.ERROR, "Has base image... no\n")
  //   log.Fatal(err)
  // } else if !isValidBaseImage {
  //   isValid = false
  //   log.Modify(-1, log.ERROR, "Has base image... no\n")
  //
  //   log.Fatal("Invalid base image.\n")
  //   return
  // }
  //
  // log.Modify(-1, log.INFO, "Has base image... ok\n")

  log.Info("Has valid Dockerfile...")
  err = checkDockerfileExists(cfg)
  if ok := log.Done(err); !ok {
    log.Fatal(err)
  }

  // if err != nil {
  //   log.Modify(-1, log.ERROR, "Has valid Dockerfile... no\n")
  //   log.Fatal(err)
  // } else if !isDockerfileExists {
  //   isValid = false
  //   log.Modify(-1, log.ERROR, "Has valid Dockerfile... no\n")
  //   log.Fatal("Dockerfile not found.\n")
  //   return
  // }
  //
  // log.Modify(-1, log.INFO, "Has valid Dockerfile... ok\n")

  log.Info("Has valid includes...")
  err = checkIncludeFilesExist(cfg)
  if ok := log.Done(err); !ok {
    log.Fatal(err)
  }

  // if err != nil {
  //   log.Modify(-1, log.ERROR, "Has valid includes... no\n")
  //   log.Fatal(err)
  // } else if !isIncludeFilesExist {
  //   isValid = false
  //   log.Modify(-1, log.ERROR, "Has valid includes... no\n")
  //   log.Fatal("Include files not found.\n")
  //   return
  // }
  //
  // log.Modify(-1, log.INFO, "Has valid includes... ok\n")

  // check := checkIncludeFilesExist(cfg)
  // isValid = check("Has valid includes...", check)

  isValid = true
  return
}

// func check(msg string, vf func(args ...interface{}) error) (valid bool) {
//   m := logger.NewLogger(logger.LoggerOptions{})
//
  // log.Info("%s...", msg)
//   vret, err := vf(args...)
//   if err != nil {
    // log.Error("%s... no\n", msg)
    // log.Fatal(err)
//   } else if !vret {
    // log.Error("%s... no\n", msg)
    // log.Fatal("Include files not found.\n")
//   }
//
  // log.Info("%s... ok\n", msg)
//
//   valid = true
//
//   return
// }

func getServerVersion(c *client.Client) (ver string, err error) {
  v, err := c.ServerVersion(context.Background())
  ver = v.Version
  return
}

func checkDockerVersion(c *client.Client) (err error) {
  minVer, err := semver.NewConstraint(">= 17.02.*-ce")
  if err != nil {
    return
  }

  currentVer, err := getServerVersion(c)
  if err != nil {
    return
  }

  v, _ := semver.NewVersion(currentVer)
  ok, _ := minVer.Validate(v)
  if !ok {
    err = fmt.Errorf("Invalid Docker version.\n")
  }

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

func checkBaseImageExists(cfg *viper.Viper, c *client.Client) (err error) {
  base := cfg.GetString("base")
  if base == "scratch" {
    return
  }

  imageTags, err := getImages(c)
  for _, tag := range imageTags {
    if base == tag {
      return
    }
  }

  err = fmt.Errorf("Invalid base image.\n")

  return
}

func getClosestMatch(cfg *viper.Viper, c *client.Client) (match string, err error) {
  base := cfg.GetString("base")

  imageTags, err := getImages(c)

  // Start asking the program to generate the levenshtein distance between two tags.
  ch := make(chan int)
  dists := make([]int, len(imageTags))

  // Use a goroutine to calculate the levenshtein distance.
  for i, tag := range imageTags {
    go levDist(tag, base, ch)
    dists[i] = <- ch
  }

  // Find the minimum distance.
  min := dists[0]
  minIndex := 0

  for i, v := range dists {
    if v < min {
        min = v
        minIndex = i
    }
  }

  match = imageTags[minIndex]

  return
}

func levDist(str1, str2 string, ch chan int) {
	matrix := make([][]int, len(str1)+1)

  for i := range matrix {
		matrix[i] = make([]int, len(str2)+1)
	}

  for i := range matrix {
		matrix[i][0] = i
	}

  for j := range matrix[0] {
		matrix[0][j] = j
	}

  for i := 1; i <= len(str1); i++ {

		for j := 1; j <= len(str2); j++ {

      if str1[i-1] == str2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				min := matrix[i-1][j]

        if matrix[i][j-1] < min {
					min = matrix[i][j-1]
				}

        if matrix[i-1][j-1] < min {
					min = matrix[i-1][j-1]
				}

        matrix[i][j] = min + 1
			}
		}

	}

	ch <- matrix[len(str1)][len(str2)]
}

func checkDockerfileExists(cfg *viper.Viper) (err error) {
  df := cfg.GetString("dockerfile")

  absPath, _ := filepath.Abs(df)

  _, err = os.Stat(absPath)

  return
}

func checkIncludeFilesExist(cfg *viper.Viper) (err error) {
  includeFiles := cfg.GetStringSlice("include")

  for _, file := range includeFiles {
    absPath, _ := filepath.Abs(file)

    if _, err = os.Stat(absPath); err != nil {
      return
    }
  }

  return
}
