package version

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var Version string = "<development>"

const upgradeUrl = "https://github.com/mynerva-io/author-cli/blob/main/docs/install.md#upgrade"

// TODO:
//		We should cache this result for a while (a day?) so that we don't hit the GitHub API
//		for every single invocation of the CLI
func getLatestVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx, "GET", "https://api.github.com/repos/mynerva-io/author-cli/releases/latest", nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create HTTP request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to query GitHub releases endpoint")
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}
	var release gitHubReleaseResponse
	if err := json.Unmarshal(data, &release); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal GitHub releases endpoint data")
	}
	return release.Name, nil
}

func CheckVersionAndPrintUpgradeNotice() {
	if Version == "<development>" {
		log.Debug("not checking for version upgrade in development")
		return
	}
	latestVersionString, err := getLatestVersion()
	if err != nil {
		log.Warnf("error occurred while trying to determine the latest version: %v", err)
		return
	}
	if latestVersionString == "" {
		log.Warnf("failed to determine latest version")
	}

	// We need to trim the "v" prefix so that semver doesn't choke
	current := *semver.New(strings.TrimLeft(Version, "v"))
	latest := *semver.New(strings.TrimLeft(latestVersionString, "v"))
	if current.LessThan(latest) {
		mag := color.New(color.Bold, color.FgMagenta)
		_, _ = fmt.Fprint(
			os.Stderr,
			mag.Sprint(" >> A new version of mynerva-author is available: "),
			color.RedString("%s", current),
			mag.Sprintf("%s", " => "),
			color.GreenString("%s", latest),
			mag.Sprintf("\n >> %s\n", upgradeUrl),
		)
	} else {
		log.Debugf("mynerva-author is up to date: current=%s, latest=%s", current, latest)
	}
}

type gitHubReleaseResponse struct {
	Name string `json:"name"`
}
