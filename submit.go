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
	zipFile  = "gochal-%s.zip" // the username is added to the filename
	errMixed = errors.New("Files and Folders both received")
)

func submit(c *cli.Context) {
	logdetails, err := getLoginDetails()
	if err != nil || logdetails.GithubUsername == "" || logdetails.ApiKey == "" {
		fmt.Println("Please log in")
		return
	}

	if !c.Args().Present() {
		fmt.Println("No arguments supplied")
		return
	}

	testDir := path.Dir(c.Args().First())
	fmt.Println("testDir")
	out, err := testsPass(testDir)
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}

	archive, err := createArchive(c.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Created " + archive)

	// TODO upload to server
}

func testsPass(testDir string) (string, error) {
	cmd := exec.Command("go", "test")
	cmd.Dir = testDir
	out, err:= cmd.Output()
	return string(out), err
}

// Creates an uncompressed .zip file containing .go files
func createArchive(args cli.Args) (string, error) {
	fi, err := os.Stat(args.First())
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		err = archiveDir(args.First())
	} else {
		// they are files
		err = archiveFiles(args)
	}

	if err != nil {
		return "", err
	}

	return archiveName()
}

func archiveDir(root string) error {
	w, err := newArchWriter()
	if err != nil {
		return err
	}
	filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		return writeToZip(w, fpath)
	})

	return w.Close()
}

func archiveFiles(args cli.Args) error {
	w, err := newArchWriter()
	if err != nil {
		return err
	}
	for _, filename := range args {
		if !strings.HasSuffix(filename, ".go") {
			continue
		}

		info, err := os.Stat(filename)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return errMixed
		}
		err = writeToZip(w, filename)
		if err != nil {
			return err
		}
	}

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

func newArchWriter() (*zip.Writer, error) {
	filename, err := archiveName()
	if err != nil {
		return nil, err
	}
	w, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return zip.NewWriter(w), nil
}

func archiveName() (string, error) {
	logdetails, err := getLoginDetails()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(zipFile, logdetails.GithubUsername), nil
}
