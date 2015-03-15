package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
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

type submissionInfo struct {
	Type string `json:"type"`
}

// Submit submits the current challenge. It checks if it correctly configured and then
// sends a zip archive of the current directory to the specified endpoint.
func Submit(c *cli.Context) {
	config, err := readConfigFile()
	if err != nil || config.APIKey == "" {
		fmt.Println("Please configure")
		return
	}

	chal, err := readChallengeFile()
	if err != nil {
		fmt.Println("Challenge not found. Please run fetch")
		return
	}
	if chal.Status != "open" {
		fmt.Printf("Sorry, %s is no longer open\n", chal.Name)
		return
	}
	fmt.Printf("Submitting %s\n", chal.Name)

	if out, err := testsPass(chal.directory()); err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}
	fmt.Println("Tests pass")

	archive, err := createArchive(chal.directory())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Created %s\n", archive)

	err = uploadFile(archive, config.APIKey, &submissionInfo{Type: c.String("type")}, chal.ID)
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

func uploadFile(archive, apikey string, info *submissionInfo, id int) error {
	challengeID := strconv.Itoa(id)
	submissionURL := strings.Join([]string{apiURL, challengeID, "submissions"}, "/")

	req, err := newfileUploadRequest(submissionURL, archive, info)
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
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Status code: %d\nBody: %s", resp.StatusCode, string(body))
	}
	return nil
}

func newfileUploadRequest(url, path string, info *submissionInfo) (*http.Request, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", "application/json; charset=utf-8")
	p, err := w.CreatePart(h)
	if err != nil {
		return nil, err
	}
	enc := json.NewEncoder(p)
	if err := enc.Encode(*info); err != nil {
		return nil, err
	}

	h = make(textproto.MIMEHeader)
	h.Set("Content-Type", "application/zip")
	p, err = w.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(p, f)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	// set type & boundary
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", w.Boundary()))
	return req, nil
}
