package client

import (
  "io"
  // "os"

  // "io/ioutil"
  "path/filepath"
  "text/template"

  "github.com/gorobot-library/orca/manifest"

  "github.com/docker/docker/pkg/archive"
)

type Context struct {
  ContextOptions

  // contextDirectory holds the temporary directory where we copy the files to
  // create the context. Any files that need to be included in the context are
  // specified in `Files`.
  contextDirectory *string

  Dockerfile string
  Files []string
}

type ContextOptions struct {
  // Directory is the current directory we will copy the context from. It is
  // generally the directory that contains the Dockerfile, and any included
  // context files.
  Directory *string

  // Data holds any variables which are available to the templates when we
  // build the context.
  Data interface{}
}

func NewContext(m manifest.Manifest, opts *ContextOptions) *Context {
  ctx := &Context{
    ContextOptions: *opts,
    Dockerfile: m.Dockerfile,
    Files: m.Files,
  }

  // If no data is passed, we default to using the manifest as the data that is
  // available to the template.
  if opts.Data == nil {
    ctx.Data = m
  }

  return ctx
}

// Tar templates the Dockerfile into a temporary directory created for the
// context, and it copies any files that were specified as includes into the
// same directory. Then it uses Docker's archive tool to create a tarfile to
// serve as the context.
//
// Returns an io.ReadCloser that is the context tarfile. It is up to the caller
// to close the reader.
func (c *Context) Tar() (io.ReadCloser, error) {
  // Create a temporary directory for the context.
  if c.contextDirectory == nil {
    dir := tempdir("", "context")
    c.contextDirectory = &dir
  }

  // Template the Dockerfile into the temporary directory.
  err := c.AddFile(c.Dockerfile)
  if err != nil {
    return nil, err
  }

  // Copy in the include files.
  for _, file := range c.Files {
    err := c.AddFile(file)
    if err != nil {
      return nil, err
    }
  }

  // Create the build context tar file.
  return archive.Tar(*c.contextDirectory, archive.Gzip)
}

// AddFile templates the file from the source directory into the context
// directory.
//
// The Dockerfile in the current directory likely contains golang templating
// variables, which means the Dockerfile as-is would not be buildable. We use
// the variables from the manifest to template the Dockerfile so that it
// becomes buildable.
//
// Other include files can also use variables, but the primary use of the
// templating feature is to modify the Dockerfile.
//
// Returns an error if templating fails.
func (c *Context) AddFile(file string) error {
  // Open the source file.
  srcPath, _ := filepath.Abs(file)

  // Create a template using the source file.
  t, err := template.ParseFiles(srcPath)
  if err != nil {
		return err
	}

  // Create the destination file.
  destPath := filepath.Join(*c.contextDirectory, file)
  dest := mustCreate(destPath)
  defer dest.Close()

  // The variables that are currently set in the c.Data object are used to
  // template the Dockerfile. Any options that are passed during the creation
  // of the Context are available during templating.
  err = t.Execute(dest, c.Data)
  if err != nil {
    return err
  }

  return nil
}
