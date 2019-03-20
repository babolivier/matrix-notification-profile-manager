package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Matrix MatrixConfig `yaml:"matrix"`
}

type MatrixConfig struct {
	HSURL       string `yaml:"hs_url"`
	AccessToken string `yaml:"access_token"`
	UserID      string `yaml:"user_id"`
}

func Load(path string) (cfg Config, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, &cfg)
	return
}
