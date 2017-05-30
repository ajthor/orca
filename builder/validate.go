package builder

import (
  "fmt"
  "os"

  "path/filepath"

  log "github.com/gorobot-library/orca/logger"

	"github.com/docker/docker/client"
  "github.com/Masterminds/semver"
  "github.com/spf13/viper"
)

func Validate(cfg *viper.Viper) {
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
    log.Warn("Orca needs Docker >= 17.05 to build images.")
    log.Fatal(err)
  }

  log.Info("Has base image...")
  err = checkBaseImageExists(cfg, c)
  if ok := log.Done(err); !ok {
    closestMatch, _ := getClosestMatch(cfg, c)
    log.Warnf("Did you mean %s?", closestMatch)
    //
    // NOTE: Prompt the user for confirmation. If yes, use the matched image.
    // NOTE: Perhaps try and download the image from Docker Hub?
    //
    log.Fatal(err)
  }

  log.Info("Has valid Dockerfile...")
  err = checkDockerfileExists(cfg)
  if ok := log.Done(err); !ok {
    log.Fatal(err)
  }

  log.Info("Has valid includes...")
  err = checkIncludeFilesExist(cfg)
  if ok := log.Done(err); !ok {
    log.Fatal(err)
  }
}

func checkDockerVersion(c *client.Client) (err error) {
  minVer, err := semver.NewConstraint(">= 17.05.*-ce")
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
