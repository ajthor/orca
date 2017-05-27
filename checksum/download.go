package checksum

import (
  "bytes"
  "io"
  "log"
  "os"
  "path"

  "net/http"
  "net/url"
  "html/template"

  "github.com/spf13/viper"
)

func downloadFile(uri string, dest string) error {
  // Create a dummy file to copy the bytes to.
  out, err := os.Create(dest)
  if err != nil {
    return err
  }

  defer out.Close()

  // Get the file from the download URL.
  r, err := http.Get(uri)
  if err != nil {
    return err
  }

  defer r.Body.Close()

  // Copy the file, byte by byte to the temporary file. This allows us to
  // download large files and not eat up memory.
  if _, err := io.Copy(out, r.Body); err != nil {
		return err
	}

  return nil
}

func getFilename(r *viper.Viper, version string) string {
  file := r.GetString("file")

  // Create the template that we will use.
  t := template.New("")
  t, err := t.Parse(file)
  if err != nil {
		log.Fatal(err)
	}

  var tpl bytes.Buffer

  // Template the filename using the version.
  if err := t.Execute(&tpl, version); err != nil {
    log.Fatal(err)
  }

  return tpl.String()

}

func getURI(r *viper.Viper, fn string) string {
  mirror := r.GetString("mirror")

  // Generate the URL for the download.
  u, err := url.Parse(mirror)
  if err != nil {
		log.Fatal(err)
	}

  u.Path = path.Join(u.Path, fn)

  return u.String()
}
