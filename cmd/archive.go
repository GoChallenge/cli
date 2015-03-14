package cmd

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	zipFile = "gochallenge.zip"
)

// Creates an uncompressed .zip file containing .go files
func createArchive(cwd string) (string, error) {
	archive, err := archiveName(path.Dir(cwd))
	if err != nil {
		return "", err
	}
	w, err := newArchWriter(archive)
	if err != nil {
		return "", err
	}

	err = archiveDir(w, cwd)
	return archive, err
}

func archiveName(dir string) (string, error) {
	return path.Join(dir, zipFile), nil
}

func newArchWriter(archive string) (*zip.Writer, error) {
	w, err := os.Create(archive)
	if err != nil {
		return nil, err
	}
	return zip.NewWriter(w), nil
}

func archiveDir(w *zip.Writer, root string) error {
	filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		return writeToZip(w, fpath)
	})

	return w.Close()
}

func writeToZip(w *zip.Writer, fpath string) error {
	f, err := w.Create(path.Base(fpath))
	if err != nil {
		return err
	}

	code, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	_, err = f.Write(code)
	return err
}
