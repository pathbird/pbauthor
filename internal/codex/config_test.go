package codex

import (
	"testing"
)

var testConfigSrc = []byte(`
[upload]
codex_category = "foo"
name = "Intro to Foo-ology"
`)

func TestUnmarshallConfig(t *testing.T) {
	config := &Config{}
	err := config.Unmarshal(testConfigSrc)
	if err != nil {
		t.Error(err)
		return
	}

	if config.Upload.Name != "Intro to Foo-ology" {
		t.Errorf("unexpected value for Upload.Name: %s", config.Upload.Name)
	}

	if config.Upload.CodexCategory != "foo" {
		t.Errorf("unexpected value for Upload.CodexCategory: %s", config.Upload.CodexCategory)
	}
}
