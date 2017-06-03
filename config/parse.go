package config

import (
  "bytes"

  "net/url"
  "path/filepath"
  "text/template"
)

// func GetFilenames(mirror, file string, versions []string) ([]string, error) {
//   files := make([]string, len(versions))
//
//   for i, ver := range versions {
//     fn, err := ParseFilename(file, ver)
//     if err != nil {
//       return files, err
//     }
//
//     files[i], err = ParseURL(mirror, fn)
//     if err != nil {
//       return files, err
//     }
//   }
//
//   return files, nil
// }

// Template the filename with the version.
func ParseFilename(file string, version string) (string, error) {
  data := struct {
    Version string
  }{
    Version: version,
  }

  t, err := template.New("").Parse(file)
  if err != nil {
		return "", err
	}

  var tpl bytes.Buffer

  // Template the filename using the version.
  if err := t.Execute(&tpl, data); err != nil {
    return "", err
  }

  return tpl.String(), nil
}

// Generate the full URL.
func ParseURL(mirror string, file string) (string, error) {
  u, err := url.Parse(mirror)
  if err != nil {
		return "", err
	}

  u.Path = filepath.Join(u.Path, file)

  return u.String(), nil
}
