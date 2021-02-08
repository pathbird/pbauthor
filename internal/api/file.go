package api

import (
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"os"
)

type FileRef struct {
	// The name of the file that will be uploaded.
	// This is usually not the path to the file on the disk, but rather,
	// the path relative to the directory that's being uploaded.
	// For example, if we have a directory structure like this:
	//   $HOME/
	//     foo/
	//	     bar.txt
	//       spam/
	//         eggs.txt
	// then if we're uploading $HOME/foo, our files would have names
	// "bar.txt" and "spam/eggs.txt".
	Name string

	// The path to the file on disk (that will be read as part of the upload)
	FsPath string
}

func (f *FileRef) addToWriter(fieldname string, w *multipart.Writer) error {
	part, err := w.CreateFormFile(fieldname, f.Name)
	if err != nil {
		return err
	}
	file, err := os.Open(f.FsPath)
	if err != nil {
		return errors.Wrapf(err, "couldn't open file (%s)", f.FsPath)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return errors.Wrapf(err, "failed to copy file (%s)", f.FsPath)
	}
	return nil
}
