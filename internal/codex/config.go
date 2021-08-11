package codex

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

const ConfigFileName = "codex.toml"

type Config struct {
	Upload UploadConfig `toml:"upload"`
	Kernel KernelConfig `toml:"kernel"`

	configFile string
}

type UploadConfig struct {
	// The ID of the codex category where the codex will be uploaded.
	CodexCategory string `toml:"codex_category"`
	// The name of the codex (as displayed in the Pathbird UI).
	Name string `toml:"name"`
	// The ID of the codex (if being re-uploaded).
	CodexId string `toml:"codex_id,omitempty"`
}

type KernelConfig struct {
	// The Docker image to use when running the kernel.
	// This overrides all other kernel config options.
	Image string `toml:"image,omitempty"`
	// An array of additional (usually Debian) packages to install
	// (e.g., {"texlive-latex-base"} if the `latex` command is required).
	SystemPackages []string `toml:"system_packages"`
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
	configFilePath := filepath.Join(dirname, ConfigFileName)
	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			return InitConfig(dirname)
		}
		return nil, errors.Wrap(err, "unable to stat codex config file")
	}
	config := &Config{}
	if err := config.UnmarshalFromFile(configFilePath); err != nil {
		return nil, err
	}
	return config, nil
}
