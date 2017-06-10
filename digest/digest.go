package digest

import (
  "encoding/json"

  "github.com/gorobot-library/orca/config"
  log "github.com/gorobot/robologger"

  "github.com/spf13/viper"
)

// Generated using: https://mholt.github.io/json-to-go/
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

// NewDigest creates a new, empty digest.
func NewDigest() *Digest {
  return &Digest{}
}
