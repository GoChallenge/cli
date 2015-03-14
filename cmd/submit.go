package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	challengeID := strconv.Itoa(id)
	submissionURL := strings.Join([]string{apiURL, challengeID, "submissions"}, "/")

	req, err := newfileUploadRequest(submissionURL, archive, subtype)
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
		return fmt.Errorf("Upload failed. Status code: %d\nBody: %s", resp.StatusCode, string(body))
	}
	return nil
}

func newfileUploadRequest(url, path, subtype string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	writer.WriteField("type", subtype)

	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", url, body)
}
