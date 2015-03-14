package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
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

func Submit(c *cli.Context) {
	config, err := getConfig()
	if err != nil || config.APIKey == "" {
		fmt.Println("Please configure")
		return
	}

	chal, err := readChallengeFile()
	if err != nil {
		fmt.Println("Challenge not found. Please fetch")
		return
	}
	fmt.Println(fmt.Sprintf("Submitting challenge %d", chal.ID))

	if out, err := testsPass(chal.directory()); err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}

	archive, err := createArchive(chal.directory())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Created " + archive)

	err = uploadFile(archive, config.APIKey, c.String("type"), chal.ID)
	if err != nil {
		fmt.Println("Upload failed - ", err)
		return
	}
	fmt.Println("Successfully submitted")
}

func testsPass(testDir string) (string, error) {
	cmd := exec.Command("go", "test")
	cmd.Dir = testDir
	out, err := cmd.Output()
	return string(out), err
}

func uploadFile(archive, apikey, subtype string, id int) error {
	reqbody, err := getReqBody(archive, subtype)
	if err != nil {
		return err
	}

	challengeID := strconv.Itoa(id)
	submissionURL := strings.Join([]string{apiURL, challengeID, "submissions"}, "/")
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

func getReqBody(archive, subtype string) ([]byte, error) {
	d, err := ioutil.ReadFile(archive)
	if err != nil {
		return nil, err
	}
	data := base64.StdEncoding.EncodeToString(d)

	hash := md5.Sum(d)
	sub := submission{hash, subtype, data}

	return json.Marshal(sub)
}
