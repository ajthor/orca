package digest

import (
  //"encoding/json"

  "github.com/gorobot-library/orca/config"
  log "github.com/gorobot/robologger"
)

type Digest struct {
  Tag string `json:"tag"`
  Commit string `json:"commit"`
  Dockerfile string `json:"dockerfile"`
  Hash string `json:"hash"`
  Platform struct {
    Architecture string `json:"architecture"`
    Os string `json:"os"`
  } `json:"platform"`
}

// Generated using: https://mholt.github.io/json-to-go/
type DigestSchema struct {
	SchemaVersion int `json:"schemaVersion"`
	Digests []Digest `json:"digests"`
}

type DigestFileOptions struct {
  Path string
}

func LoadDigest(opts *DigestFileOptions) *DigestSchema {
  cfg, err := config.New(opts.Path)
  if err != nil {
    log.Fatal(err)
  }

  req := []string{
    "version",
    "manifests",
  }
  if err := config.HasRequired(cfg, req); err != nil {
    log.Fatal(err)
  }

  s := &DigestSchema{}

  err = cfg.Unmarshal(s)
  if err != nil {
  	log.Fatalf("Unable to parse %s, %v", opts.Path, err)
  }

  return s
}

// func (d *Digest) Marshal() ([]byte, error) {
//   return json.Marshal(d), nil
// }
//
// func (d *Digest) Unmarshal()  {
//
// }
