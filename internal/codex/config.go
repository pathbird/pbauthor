package codex

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"io/ioutil"
)

type Config struct {
	Upload UploadConfig `toml:"upload"`
}

type UploadConfig struct {
	CodexCategory string `toml:"codex_category"`
	Name          string `toml:"name"`
}

func UnmarshallConfig(data []byte) (*Config, error) {
	config := &Config{}
	err := toml.Unmarshal(data, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse codex config")
	}

	if err := config.Upload.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (u *UploadConfig) validate() error {
	if u.CodexCategory == "" {
		return errors.New("upload.codex_category must be specified")
	}
	return nil
}

func ReadConfigFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return UnmarshallConfig(data)
}
