package checksum

import (
  "bytes"
  "fmt"
  "log"
  "os"

  "crypto/sha256"
  "html/template"
  "io/ioutil"
  "net/url"
  "path/filepath"
)

func generateFilename(f *fileInfo, version string) {
  // Create the template that we will use.
  t := template.New("")
  t, _ = t.Parse(f.fn)

  if err := t.Execute(&tpl, version); err != nil {
    return err
  }

  // Template the filename using the version.
  f.fn = &bytes.Buffer
}

func generateURI(f *fileInfo) {
  // Generate the URL for the download.
  u, err := url.Parse(uri)
  u.Path = path.Join(u.Path, f.fn)

  f.uri = u.String()
}

func generateShasum(f *fileInfo) {
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
  return h.Sum(nil)
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
