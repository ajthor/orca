package checksum

import (
  "bytes"
  "fmt"
  "io"
  "log"
  "os"

  "crypto/sha256"
  "html/template"
  "net/url"
  "path/filepath"
)

func generateFilename(fn string, version string) string {
  // Create the template that we will use.
  t := template.New("")
  t, _ = t.Parse(fn)

  if err := t.Execute(&tpl, version); err != nil {
    return err
  }

  // Template the filename using the version.
  return &bytes.Buffer
}

func generateURI(fn string, uri string) string {
  // Generate the URL for the download.
  u, err := url.Parse(uri)
  u.Path = path.Join(u.Path, fn)

  return u.String()
}

func generateShasum(fn string) {
  // Open the file specified by the fn argument.
  f, err := os.Open(fn)
  if err != nil {
		log.Fatal(err)
	}

  // Make sure we close it afterward.
  defer f.Close()

  // Generate the hash.
  h := sha256.New()
  if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

  // Output the checksum to the file.
  return h.Sum(nil) + "  " + fn
}

func generateShasumFile(fn string, hashes []string) {
  f, err := os.Create(fn)
  if err != nil {
    log.Fatal(err)
  }

  defer f.Close()

  for _, h := range hashes {
    if _, err := io.Copy(f, h); err != nil {
      log.Fatal(err)
    }
  }
}
