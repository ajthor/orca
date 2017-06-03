package builder

import (
  "fmt"
  "os"

  "path/filepath"


  "github.com/docker/docker/client"
  "github.com/Masterminds/semver"
  log "github.com/gorobot/robologger"
)

var MinDockerVersion string = ">= 17.05.*-ce"

func (b *Builder) Validate(opts *BuildOptions) error {
  log.Info("Checking configuration...")

  m1 := log.Info("Has compatible docker version...")
  err := checkDockerVersion(b.Client)
  if err != nil {
    log.Status(log.ERR, m1)
    log.Errorf("Orca needs Docker %s to build images.", MinDockerVersion)
    return err
  }
  log.Status(log.OK, m1)

  m2 := log.Info("Has base image...")
  err = checkBaseImageExists(opts, b.Client)
  if err != nil {
    closestMatch, _ := getClosestMatch(opts, b.Client)
    res := log.Promptf(log.YESNO, "Did you mean %s?", closestMatch)
    pres, _ := log.ParseResponse(res)
    if pres == log.YES {
      b.Config.Base = closestMatch
    }
    //
    // NOTE: Perhaps try and download the image from Docker Hub?
    //
    return err
  }
  log.Status(log.OK, m2)

  m3 := log.Info("Has valid Dockerfile...")
  err = checkDockerfileExists(opts)
  if err != nil {
    log.Status(log.ERR, m3)
    return err
  }
  log.Status(log.OK, m2)


  m4 := log.Info("Has valid includes...")
  err = checkIncludeFilesExist(opts)
  if err != nil {
    log.Status(log.ERR, m4)
    return err
  }
  log.Status(log.OK, m4)

  return nil
}

func checkDockerVersion(c *client.Client) error {
  minVer, err := semver.NewConstraint(MinDockerVersion)
  if err != nil {
    return err
  }

  currentVer, err := getServerVersion(c)
  if err != nil {
    return err
  }

  v, _ := semver.NewVersion(currentVer)
  ok, _ := minVer.Validate(v)
  if !ok {
    err = fmt.Errorf("Invalid Docker version.\n")
  }

  return nil
}

func checkBaseImageExists(opts *BuildOptions, c *client.Client) error {
  base := opts.Base
  if base == "scratch" {
    return nil
  }

  imageTags, err := getImages(c)
  for _, tag := range imageTags {
    if base == tag {
      return err
    }
  }

  err = fmt.Errorf("Invalid base image.\n")

  return nil
}

func getClosestMatch(opts *BuildOptions, c *client.Client) (match string, err error) {
  base := opts.Base

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

func checkDockerfileExists(opts *BuildOptions) error {
  absPath, _ := filepath.Abs(opts.Dockerfile)

  if _, err := os.Stat(absPath); os.IsNotExist(err) {
    return err
  }

  return nil
}

func checkIncludeFilesExist(opts *BuildOptions) error {
  for _, file := range opts.Includes {
    absPath, _ := filepath.Abs(file)

    if _, err := os.Stat(absPath); os.IsNotExist(err) {
      return err
    }
  }

  return nil
}
