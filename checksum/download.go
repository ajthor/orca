package checksum

import (
  "io"
  "os"

  "net/http"
)

func (c *Checksum) downloadFile(uri string, dest string) (err error) {
  // Create a dummy file to copy the bytes to.
  out, err := os.Create(dest)
  if err != nil {
    return
  }

  defer out.Close()

  // Get the file from the download URL.
  r, err := http.Get(uri)
  if err != nil {
    return
  }

  defer r.Body.Close()

  // _, err = ioutil.ReadAll(r.Body)
	// r.Body.Close()
	// if err != nil {
	// 	return
	// }

  // Copy the file, byte by byte to the temporary file. This allows us to
  // download large files and not eat up memory.
  _, err = io.Copy(out, r.Body)
  if err != nil {
		return
	}

  return
}
