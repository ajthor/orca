package digest

import (
  "encoding/json"

  "github.com/gorobot-library/orca/config"
  log "github.com/gorobot/robologger"

  "github.com/spf13/viper"
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

func NewDigest() *DigestSchema {
  return &DigestSchema{}
}

// Marshal wraps a json.Marshal call on the digest.
func (s *DigestShema) Marshal() ([]byte, error) {
  return json.Marshal(s)
}

// Unmarshal takes a viper config file and unmarshals the data into the
// DigestSchema.
func (d *DigestSchema) Unmarshal(cfg *viper.Viper) error {
  err = cfg.Unmarshal(d)
  if err != nil {
  	return err
  }
}

// NewDigest creates a new, empty digest.
func NewDigest() *Digest {
  return &Digest{}
}
