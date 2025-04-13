package jenkins

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DeployItem method
func (c *Client) DeployItem(jobName string) error {
	url := fmt.Sprintf("/job/%s/build", jobName)
	req, err := c.newRequest("POST", url)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to trigger build for job %s, status code: %d", jobName, resp.StatusCode)
	}

	fmt.Printf("Build for job %s triggered successfully. Waiting for completion...\n", jobName)

	buildDetail, err := c.GetLastBuildNumber(jobName)
	if err != nil {
		return fmt.Errorf("failed to fetch build log for job %s, error: %s", jobName, err)
	}

	if strings.Contains(buildDetail, "Finished: FAILURE") || strings.Contains(buildDetail, "ERROR") || strings.Contains(buildDetail, "FAILED") {
		return fmt.Errorf("job %s failed, check logs for details:\n%s", jobName, buildDetail)
	}
	return nil
}

// GetLastBuildNumber method -- Get deployment process of a job
func (c *Client) GetLastBuildNumber(jobName string) (string, error) {
	url := fmt.Sprintf("/job/%s/lastBuild/logText/progressiveText", jobName)
	req, err := c.newRequest("GET", url)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get last build number for job %s, status code: %d", jobName, resp.StatusCode)
	}

	// Đọc toàn bộ log từ response body
	logData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(logData), nil
}
