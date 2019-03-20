package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the structure of the configuration file.
type Config struct {
	Matrix MatrixConfig `yaml:"matrix"`
}

// MatrixConfig represents the structure of the `matrix` part of the
// configuration file.
type MatrixConfig struct {
	HSURL       string `yaml:"hs_url"`
	AccessToken string `yaml:"access_token"`
	UserID      string `yaml:"user_id"`
}

// Load loads the content of a configuration file into a Config instance.
// Returns an error if reading the file or parsing the YAML data failed.
func Load(path string) (cfg Config, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, &cfg)
	// TODO: Check if all required fields are set.
	return
}
