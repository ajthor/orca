package checksum

import (
  "bytes"
  "io"
  "log"
  "os"
  "path"

  "crypto/sha256"
  "html/template"
  "net/url"
)

func generateFilename(fn string, version string) string {
  // Create the template that we will use.
  t := template.New("")
  t, _ = t.Parse(fn)

  var tpl bytes.Buffer

  if err := t.Execute(&tpl, version); err != nil {
    log.Fatal(err)
  }

  // Template the filename using the version.
  return tpl.String()
}

func generateURI(fn string, uri string) string {
  // Generate the URL for the download.
  u, err := url.Parse(uri)
  if err != nil {
		log.Fatal(err)
	}
  
  u.Path = path.Join(u.Path, fn)

  return u.String()
}

func generateShasum(fn string) []byte {
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
  return append(h.Sum(nil), ("  " + fn)...)
}

func generateShasumFile(fn string, hashes [][]byte) {
  f, err := os.Create(fn)
  if err != nil {
    log.Fatal(err)
  }

  defer f.Close()

  for _, h := range hashes {
    r := bytes.NewReader(h)
    if _, err := io.Copy(f, r); err != nil {
      log.Fatal(err)
    }
  }
}
