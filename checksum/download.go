package checksum

import (
  "bytes"
  "io"
  "os"
  "path"


  "net/http"
  "net/url"
  "html/template"

  log "github.com/gorobot-library/orca/logger"

  "github.com/spf13/viper"
)

func downloadFile(uri string, dest string) (err error) {
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

func GetFilename(cfg *viper.Viper, version string) string {
  file := cfg.GetString("file")

  // Create the template that we will use.
  // t := template.New("")
  type Version struct {
    Version string
  }

  // funcMap := template.FuncMap{
  //   "version": version,
  // }
  t, err := template.New("").Parse(file)
  if err != nil {
		log.Fatal(err)
	}

  var tpl bytes.Buffer

  // Template the filename using the version.
  if err := t.Execute(&tpl, Version{version}); err != nil {
    log.Fatal(err)
  }

  return tpl.String()

}

func GetURI(cfg *viper.Viper, fn string) string {
  mirror := cfg.GetString("mirror")

  // Generate the URL for the download.
  u, err := url.Parse(mirror)
  if err != nil {
		log.Fatal(err)
	}

  u.Path = path.Join(u.Path, fn)

  return u.String()
}
