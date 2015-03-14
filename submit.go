package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
)

const (
	apikeyHeader = "Auth-ApiKey"
)

var (
	errAuthFail = errors.New("Authkey rejected. Please configure")
)

// TODO get the challenge ID and submission type before submission
type submission struct {
	Hash [16]byte `json:"hash"`
	Type string   `json:"type"`
	Data string   `json:"data"`
}

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
	err = uploadFile(config.ApiKey, archive)
	if err != nil {
		fmt.Println("Upload failed - ", err)
		return
	}
	fmt.Println("Successfully uploaded")
}

func testsPass(testDir string) (string, error) {
	cmd := exec.Command("go", "test")
	cmd.Dir = testDir
	out, err := cmd.Output()
	return string(out), err
}

func uploadFile(apikey, archive string) error {
	reqbody, err := getReqBody(archive)
	if err != nil {
		return err
	}

	challengeID := "1"
	submissionURL := strings.Join([]string{apiUrl, challengeID, "submissions"}, "/")
	req, err := http.NewRequest("POST", submissionURL, bytes.NewReader(reqbody))
	if err != nil {
		return err
	}
	req.Header.Add(apikeyHeader, apikey)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return errAuthFail
	}
	return nil
}

func getReqBody(archive string) ([]byte, error) {
	d, err := ioutil.ReadFile(archive)
	if err != nil {
		return nil, err
	}
	data := base64.StdEncoding.EncodeToString(d)

	hash := md5.Sum(d)
	subtype := "normal"
	sub := submission{hash, subtype, data}

	return json.Marshal(sub)
}
