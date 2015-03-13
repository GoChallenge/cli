package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
)

var (
	zipFile  = "gochallenge.zip" // the username is added to the filename
)

func submit(c *cli.Context) {
	config, err := getConfig()
	if err != nil || config.ApiKey == "" {
		fmt.Println("Please configure")
		return
	}
	
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := testsPass(cwd)
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}

	archive, err := createArchive(cwd)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Created " + archive)

	// TODO
	uploadFile(archive)	
}

func testsPass(testDir string) (string, error) {
	cmd := exec.Command("go", "test")
	cmd.Dir = testDir
	out, err := cmd.Output()
	return string(out), err
}

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
	// TODO rename file with current challenge
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

func uploadFile(archive string) {

}
