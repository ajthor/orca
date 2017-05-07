package checksum

import (
  "io"
  "log"
  "os"

  "net/http"
)

func downloadFile(fn string, uri string) {
  // Create a dummy file to copy the bytes to.
  out, err := os.Create(fn)
  if err != nil {
    log.Fatal(err)
  }

  defer out.Close()

  // Get the file from the download URL.
  r, err := http.Get(uri)
  if err != nil {
    log.Fatal(err)
  }

  defer r.Body.Close()

  // Copy the file, byte by byte to the temporary file. This allows us to
  // download large files and not eat up memory.
  if _, err := io.Copy(out, r.Body); err != nil {
		log.Fatal(err)
	}
}
