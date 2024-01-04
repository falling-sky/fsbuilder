package crowdinio

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

// config contains configuration options
type config struct {
	Token     string `json:"token" yaml:"token"`
	ProjectID int    `json:"project_id" yaml:"project_id"`
}

// load a config file, return it after adjusting for defaults
func load(filename string) (*config, error) {
	r := &config{}

	// If a filename is specified, load it.

	b, e := ioutil.ReadFile(filename)
	if e != nil {
		return r, e
	}
	e = json.Unmarshal(b, r)
	if e != nil {
		return r, e
	}
	return r, nil
}

func (r *config) String() string {
	b, e := json.MarshalIndent(r, "", "\t")
	if e != nil {
		log.Fatal(e)
	}

	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

	return string(b)
}

// Return a sample config with defaults
func Example() string {
	r := &config{}
	return r.String()
}
