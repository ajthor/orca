package digest

import (
  "encoding/json"

  "github.com/gorobot-library/orca/config"

  "github.com/spf13/viper"
)

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
