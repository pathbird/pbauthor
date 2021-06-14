package auth

import (
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"
)

type Auth struct {
	ApiToken   string    `json:"apiToken"`
	Expiration time.Time `json:"expirationTime"`
}

func getApiTokenCacheFile() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "unable to determine auth file")
	}
	dir := path.Join(currentUser.HomeDir, ".pathbird")
	_ = os.MkdirAll(dir, 0700)
	return path.Join(currentUser.HomeDir, ".pathbird", "auth.json"), nil
}

var authCache = (*Auth)(nil)

func GetAuth() (*Auth, error) {
	if authCache != nil {
		return authCache, nil
	}

	file, err := getApiTokenCacheFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Debugf("authentication cache file does not exist")
		return nil, nil
	}

	fp, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open auth file (%s)", file)
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read auth file (%s)", file)
	}

	var auth Auth
	err = json.Unmarshal(data, &auth)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse auth file (%s)", file)
	}

	if time.Now().After(auth.Expiration) {
		log.Debugf("removing expired authentication cache")
		if err := os.Remove(file); err != nil {
			return nil, errors.Wrapf(err, "failed to remove expired auth file (%s)", file)
		}
	}

	authCache = &auth
	return &auth, nil
}

func GetAuthApiToken() (string, error) {
	auth, err := GetAuth()
	if auth == nil || err != nil {
		return "", err
	}
	return auth.ApiToken, nil
}

func SaveAuth(auth *Auth) error {
	authCache = auth
	file, err := getApiTokenCacheFile()
	if err != nil {
		return err
	}

	b, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, b, 0600)
	if err != nil {
		return errors.Wrapf(err, "failed to write auth file (%s)", file)
	}

	return nil
}
