package manifest

import (
  "errors"
  "reflect"

  "encoding/json"

  "github.com/spf13/viper"
)

var (
  ErrManifestNotFound = errors.New("Manifest not found.")
)

// Generated using: https://mholt.github.io/json-to-go/
type ManifestSchema struct {
  SchemaVersion int `json:"schemaVersion"`
  Manifests []Manifest `json:"manifests"`
}

// NewSchema creates a new ManifestSchema that can be unmarshaled from a config
// file.
func NewSchema() *ManifestSchema {
  return &ManifestSchema{}
}

// Marshal wraps a json.Marshal call on the manifest.
func (s ManifestSchema) Marshal() ([]byte, error) {
  return json.Marshal(s)
}

// Unmarshal takes a viper config file and unmarshals the data into the
// ManifestSchema.
func (s *ManifestSchema) Unmarshal(cfg *viper.Viper) error {
  err := cfg.Unmarshal(s)
  if err != nil {
  	return err
  }

  return nil
}

// Get returns a manifest that has a name matching the supplied name. Returns
// a nil pointer if no manifest matching the provided name can be found.
func (s ManifestSchema) Get(name string) *Manifest {
  for i, m := range s.Manifests {
    if name == m.Name {
      return &s.Manifests[i]
    }
  }

  return nil
}

// Filter filters the manifests in the schema and returns the matches as a
// slice. If no matches are found, the function returns an empty slice.
func (s ManifestSchema) Filter(f func(Manifest) bool) []*Manifest {
  matches := make([]*Manifest, 0)
  for i, m := range s.Manifests {
    if f(m) {
      matches = append(matches, &s.Manifests[i])
    }
  }

  return matches
}

// Match returns the first matching manifest based on the constraints specified
// by the map[string]interface{} argument. It returns the first manifest found
// or nil if no manifests matching the constraints are found.
func (s ManifestSchema) Match(args map[string]interface{}) *Manifest {
  // Iterate through all manifests.
  for i, m := range s.Manifests {
    val := reflect.ValueOf(m)

    // Go through the constraints, one by one.
    for k, v := range args {
      field := val.FieldByName(k)
      if !field.IsValid() {
        break
      }

      if field.Kind() == reflect.Slice {
        for j := 0; j < field.Len(); j++ {
          if reflect.DeepEqual(v, field.Index(j).Interface()) {
            return &s.Manifests[i]
          }
        }
        break
      }

      if reflect.DeepEqual(v, field.Interface()) {
        return &s.Manifests[i]
      }
    }
  }

  return nil
}

// MatchAll returns the first matching manifest based on the constraints
// specified by the map[string]interface{} argument. It returns all manifests
// found.
func (s ManifestSchema) MatchAll(args map[string]interface{}) []*Manifest {
  matches := make([]*Manifest, 0)
  // Iterate through all images.
  for i, m := range s.Manifests {
    val := reflect.ValueOf(m)

    // Go through the constraints, one by one.
    for k, v := range args {
      field := val.FieldByName(k)
      if !field.IsValid() {
        break
      }

      if field.Kind() == reflect.Slice {
        for j := 0; j < field.Len(); j++ {
          if reflect.DeepEqual(v, field.Index(j).Interface()) {
            matches = append(matches, &s.Manifests[i])
          }
        }
        break
      }

      if reflect.DeepEqual(v, field.Interface()) {
        matches = append(matches, &s.Manifests[i])
      }
    }
  }

  return matches
}
