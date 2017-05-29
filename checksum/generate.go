package checksum

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "path"

  "crypto/sha256"

  log "github.com/gorobot-library/orca/logger"
)

func generateShasum(localFile string) string {
  file, err := os.Open(localFile)
  if err != nil {
		log.Fatal(err)
	}

  defer file.Close()

  // Generate the hash.
  h := sha256.New()
  if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

  fn := path.Base(localFile)

  // The hash line is a combination of the hash, two spaces, and the filename.
  // fullHash := append(h.Sum(nil), ("  " + fn)...)
  fullHash := fmt.Sprintf("%x", h.Sum(nil)) + "  " + fn
  // fmt.Println(fullHash)

  return fullHash
}

func removeShasumFile(shasumFile string) error {
  // Check if the file exists, and remove it if it does.
  if _, err := os.Stat(shasumFile); err == nil {
    if err := os.Remove(shasumFile); err != nil {
      return err
    }
  }

  return nil
}

func createShasumFile(fn string, hashes []string) error {
  f, err := os.Create(fn)
  if err != nil {
    return err
  }

  defer f.Close()

  w := bufio.NewWriter(f)

  for _, h := range hashes {
    _, err := w.WriteString(h + "\n")
    if err != nil {
      return err
    }

    // fmt.Printf("wrote %d bytes\n", n)
  }

  w.Flush()

  return nil
}
