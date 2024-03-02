package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var client = resty.New().
	SetRetryCount(6).
	SetRetryWaitTime(2 * time.Second).
	SetRetryMaxWaitTime(60 * time.Second).
	AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() == http.StatusTooManyRequests
		},
	)

func IsOutdated(currentVersion string, repo string) bool {
	var tags []Tag
	var url = fmt.Sprintf("https://api.github.com/repos/%s/tags", repo)

	resp, err := client.R().SetResult(&tags).Get(url)
	if resp.StatusCode() != 200 || err != nil {
		return true
	}

	return currentVersion != tags[0].Name
}

func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return strings.Replace(path, "~", usr.HomeDir, 1), nil
	}
	return path, nil
}

func FetchUrl(url string) string {
	resp, err := client.R().Get(url)
	if err != nil {
		PrintError("Error fetching URL: %s", err)
	}

	return resp.String()
}

func CreateReport(directory string, downloads []Download) {
	filePath := filepath.Join(directory, "_report.md")
	file, err := os.Create(filePath)
	if err != nil {
		return
	}

	defer file.Close()

	// Filter the failed downloads
	failedDownloads := make([]Download, 0)
	for _, download := range downloads {
		if !download.IsSuccess {
			failedDownloads = append(failedDownloads, download)
		}
	}

	fileContent := "# Coomer - Download Report\n"
	fileContent += "## Failed Downloads\n"
	fileContent += fmt.Sprintf("- Total: %d\n", len(failedDownloads))

	for _, download := range failedDownloads {
		fileContent += fmt.Sprintf("### 🔗 Link: %s - ❌ **Failure**\n", download.Url)
		fileContent += "### 📝 Error:\n"
		fileContent += "```\n"
		fileContent += fmt.Sprintf("%s\n", download.Error)
		fileContent += "```\n"
		fileContent += "---\n"
	}

	_, _ = file.WriteString(fileContent)
}
