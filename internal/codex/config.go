package codex

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Upload UploadConfig `toml:"upload"`

	configFile string
}

type UploadConfig struct {
	CodexCategory string `toml:"codex_category,omitempty"`
	Name          string `toml:"name,omitempty"`

	// The id of the codex (if being re-uploaded)
	CodexId string `toml:"codex_id,omitempty"`
}

func (c *Config) Unmarshal(data []byte) error {
	err := toml.Unmarshal(data, c)
	if err != nil {
		return errors.Wrap(err, "failed to parse codex config")
	}

	if err := c.Upload.validate(); err != nil {
		return errors.Wrap(err, "failed to validate codex config")
	}
	return nil
}

func (c *Config) UnmarshalFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "failed to read codex config file (%s)", filename)
	}
	err = c.Unmarshal(data)
	if err != nil {
		return err
	}
	c.configFile = filename

	log.Debugf("read config file: %#v", c)
	return nil
}

func (c *Config) Save() error {
	if c.configFile == "" {
		return errors.New("cannot save codex config: no file selected")
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "failed to marshal codex config")
	}

	err = ioutil.WriteFile(c.configFile, data, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to save codex config file (%s)", c.configFile)
	}

	return nil
}

func (u *UploadConfig) validate() error {
	if u.CodexCategory == "" {
		return errors.New("upload.codex_category must be specified")
	}
	return nil
}

func GetOrInitCodexConfig(dirname string) (*Config, error) {
	configFilePath := filepath.Join(dirname, "codex.toml")
	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			return initConfig(configFilePath)
		}
		return nil, errors.Wrap(err, "unable to stat codex.toml")
	}
	config := &Config{}
	if err := config.UnmarshalFromFile(configFilePath); err != nil {
		return nil, err
	}
	return config, nil
}
