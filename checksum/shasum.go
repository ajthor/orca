package checksum

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "strings"

  "crypto/sha256"
  "path/filepath"

  log "github.com/gorobot-library/orca/logger"
)

func (c *Checksum) GetShasum(match string) (shasum string, err error) {
  if _, err := os.Stat(*c.ShaFile); os.IsNotExist(err) {
    return "", err
  }

  file, err := os.Open(*c.ShaFile)
  if err != nil {
    return "", err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    if contains := strings.Contains(scanner.Text(), match); contains {
      return scanner.Text(), nil
    }
  }

  if err := scanner.Err(); err != nil {
    return "", err
  }

  return "", fmt.Errorf("Shasum for %s not found.", match)
}

func (c *Checksum) GenerateShasum(file string) string {
  f, err := os.Open(file)
  if err != nil {
		log.Fatal(err)
	}

  defer f.Close()

  // Generate the hash.
  h := sha256.New()
  if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

  // The hash line is a combination of the hash, two spaces, and the filename.
  shasum := fmt.Sprintf("%x", h.Sum(nil)) + "  " + filepath.Base(file)

  return shasum
}

func (c *Checksum) CreateShasumFile(hashes []string) error {
  log.Debugf("---> %s", *c.ShaFile)

  // Check if the file exists, and remove it if it does.
  if _, err := os.Stat(*c.ShaFile); err == nil {
    res := log.Promptf(log.YESNO, "%s already exists. Overwrite?", *c.ShaFile)
    log.ShowInput(res)

    fres := log.FormatResponse(res)
    if fres != log.YES {
      return fmt.Errorf("Canceling...")
    }
    
    if err := os.Remove(*c.ShaFile); err != nil {
      return err
    }
  }

  // Then, we recreate the file.
  f, err := os.Create(*c.ShaFile)
  if err != nil {
    return err
  }

  defer f.Close()

  // And write the hashes to the file.
  w := bufio.NewWriter(f)

  for _, h := range hashes {
    _, err := w.WriteString(h + "\n")
    if err != nil {
      return err
    }
  }

  w.Flush()

  return nil
}
