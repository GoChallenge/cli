package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
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

	archive, err := createArchive(c.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Created " + archive)

	// TODO upload to server
}

// For now this creates a gzipped tar file. This can be changed
func createArchive(args cli.Args) (string, error) {
	fi, err := os.Stat(args.First())
	if err != nil {
		return "", err
	}

	if fi.IsDir() == true {
		return archiveDir(args.First())
	}
	// they are files
	tw, err := newArchWriter("gochallenge")
	for _, file := range args {
		info, _ := os.Stat(file)
		header, _ := tar.FileInfoHeader(info, "")
		header.Name = info.Name()
		tw.WriteHeader(header)
		data, _ := ioutil.ReadFile(file)
		tw.Write(data)
		tw.Flush()
	}
	tw.Close()
	return "gochallenge.tar", nil
}

func archiveDir(root string) (string, error) {
	fmt.Println("Creating archive of", root)
	dir := filepath.Dir(root)
	tw, err := newArchWriter(root)
	if err != nil {
		return "", err
	}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		header, _ := tar.FileInfoHeader(info, "")
		header.Name = path[len(dir):]
		tw.WriteHeader(header)
		if info.IsDir() == false {
			data, _ := ioutil.ReadFile(path)
			tw.Write(data)
			tw.Flush()
		}
		return nil
	})
	tw.Close()
	fmt.Println("Created", root+".tar")
	return root + ".tar", nil
}

func newArchWriter(name string) (*tar.Writer, error) {
	w, err := os.Create(name + ".tar")
	if err != nil {
		return new(tar.Writer), err
	}
	cw := gzip.NewWriter(w)
	return tar.NewWriter(cw), nil
}
