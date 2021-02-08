package codex

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGetCodexFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if err := ioutil.WriteFile(filepath.Join(dir, "foo.txt"), []byte(`hello`), 0x644); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(dir, ".bar.txt"), []byte(`hello`), 0x644); err != nil {
		t.Fatal(err)
	}

	files, err := getCodexFiles(&Config{}, dir)
	if len(files) != 1 {
		t.Fatalf("expected one file, got: %d", len(files))
	}

	f := files[0]
	if f.Name != "foo.txt" {
		t.Errorf("unexpected file name: %s", f.Name)
	}
	if f.FsPath != filepath.Join(dir, "foo.txt") {
		t.Errorf("unexpected file path: %s", f.FsPath)
	}
}
