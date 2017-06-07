package client

import (
  // "bufio"
  "errors"
  "fmt"
  "io"
  "os"
  "strings"

  "crypto/sha256"
  "net/http"
  "path/filepath"

  // "github.com/gorobot-library/orca/manifest"

  log "github.com/gorobot/robologger"
)

var (
  ErrInvalidHttpResponse = errors.New("Invalid http response.")
)

// GenerateShasums takes a manifest and generates shasums for the files and
// versions specified in the file.
func (c *Client) GenerateShasums(files, urls []string) ([]*Shasum, error) {
  // Create an empty slice to hold our generated shasums.
  hashes := make([]*Shasum, len(files))

  dir := tempdir("", "downloads")

  for i, dlurl := range urls {
    dlfile := filepath.Join(dir, files[i])

    // Create an empty file to copy the bytes to.
    dest := mustCreate(dlfile)
    defer dest.Close()

    // Download the file to the temporary directory.
    err := downloadFile(dlurl, dest)
    if err != nil {
      return hashes, err
    }

    // Create a shasum for the downloaded file.
    hashes[i], err = NewShasum(dlfile)
    if err != nil {
      return hashes, err
    }
  }

  return hashes, nil
}

// Shasum is an interface for all shasum types used in this file. Shasums in this context are at least sha256.
type shasum interface {
  fmt.Stringer
  File()
  Hash()
  Generate()
}

type Shasum struct {
  raw string
}

// NewShasum returns a new Shasum. If the string is a filepath, the
// function attempts to generate a shasum.
func NewShasum(s string) (*Shasum, error) {
  if ok := ValidateShasum(s); !ok {
    return nil, ErrInvalidHash
  }

  absPath, _ := filepath.Abs(s)

  // If the string is a file and it exists, we generate a hash for the file.
  if _, err := os.Stat(absPath); os.IsExist(err) {
    sum, err := GenerateShasum(s)
    if err != nil {
      return nil, err
    }

    return sum, nil
  }

  return &Shasum{
    raw: s,
  }, nil
}

// GenerateShasum reads in the file specified as a parameter and returns a
// pointer to a new Shasum.
func GenerateShasum(file string) (*Shasum, error) {
  f := mustOpen(file)
  defer f.Close()

  // CreateShasum the hash.
  h := sha256.New()
  if _, err := io.Copy(h, f); err != nil {
    return nil, err
  }

  // The hash line is a combination of the hash, two spaces, and the filename.
  shasum := fmt.Sprintf("%x", h.Sum(nil)) + "  " + filepath.Base(file)

  sum := &Shasum{
    raw: shasum,
  }

  return sum, nil
}

// ValidateShasum does rudimentary checking to ensure that the hash read from a
// file is in the correct format.
func ValidateShasum(h string) bool {
  // Makes sure that the string contains a double space.
  if contains := strings.Contains(h, "  "); !contains {
    return false
  }

  // Make sure that the length of a string split around a double space is two.
  s := strings.Split(h, "  ")
  if stringLen := len(s); stringLen != 2 {
    return false
  }

  // Checks to make sure the hash is at least 64 bytes long.
  if hashLen := len(s[0]); hashLen >= 64 {
    return false
  }

  return true
}

// String returns a string representation of the full hash, i.e. (hash + "  " +
// file).
func (s *Shasum) String() string {
  return s.raw
}

// Hash returns the hash only.
func (s *Shasum) Hash() string {
  return strings.Split(s.raw, "  ")[0]
}

// File returns the filename only.
func (s *Shasum) File() string {
  return strings.Split(s.raw, "  ")[1]
}

// fileDownloader wraps the io.ReadCloser returned from the http.Get function
// so that we can display a progress bar alongside the download. To accomplish
// this, we incorporate our own 'Read' and 'Close' methods for the ReadCloser
// interface.
type fileDownloader struct {
  io.ReadCloser
  Total int64
  resp *http.Response
  Update func(int, interface{})
}

func (dl *fileDownloader) Read(p []byte) (int, error) {
  n, err := dl.ReadCloser.Read(p)
  dl.Total += int64(n)

  if err == nil {
    dl.Update(int(100 * (float32(dl.Total) / float32(dl.resp.ContentLength))), "downloading...")
  }

  return n, err
}

// Close closes the fileDownloader reader.
func (dl *fileDownloader) Close() error {
  return dl.ReadCloser.Close()
}

// downloadFile fetches the file from `url` and saves it to the file specified by `name`.
func downloadFile(url string, file *os.File) error {
  m := log.Infof("Downloading %s...", url)

  // Get the file from the download URL.
  r, err := http.Get(url)
  if err != nil {
    return err
  }

  defer r.Body.Close()

  // If the response from the server is not 200 (OK), we have an error in the
  // file download and return an error.
  if err := checkHttpResponse(r.StatusCode); err != nil {
    return err
  }

  Update := log.Progress()

  dl := &fileDownloader{
    ReadCloser: r.Body,
    resp: r,
    Update: Update,
  }

  // Copy the file, byte by byte to the temporary file. This allows us to
  // download large files and not eat up memory.
  _, err = io.Copy(file, dl)
  if err != nil {
    return err
  }

  Update(100, "done")
  log.Status(log.OK, m)

  return nil
}

// checkHttpResponse uses the response received from the http.Get function and
// determines if it is alright to continue the download.
func checkHttpResponse(code int) error {
  switch {
  case code == 200:
    return nil
  // Place all invalid responses here with a fallthrough to the error.
  default:
    fallthrough
  case code >= 400:
    return ErrInvalidHttpResponse
  }
}

//
// // GetSha256 takes the filename specified as a parameter and searches the
// // SHASUM256.txt file for a matching filename. If one is found, it returns the
// // entire line in the file.
// func (c *Client) GetSha256(file string, match string) (string, error) {
//   // Get all hashes from the shasum file.
//   hashes, err := c.ReadShasumFile(file)
//   if err != nil {
//     return "", err
//   }
//
//   // Return a match if there is one.
//   for _, h := range hashes {
//     if contains := strings.Contains(h, match); contains {
//       return h, nil
//     }
//   }
//
//   return "", fmt.Errorf("Shasum not found.")
// }
