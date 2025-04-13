package jenkins

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Structs for fetching jobs
type Job struct {
	Name string
}

type JobsResponse struct {
	Jobs []Job
}

// Struct for checking exist jobs
type ExistJob struct {
	Name  string
	Exist bool
}

// FetchAllJobs method
func (c *Client) FetchAllJobs() (*JobsResponse, error) {
	req, err := c.newRequest("GET", "/api/json")
	if err != nil {
		return nil, err
	}

	client := &http.Client{} // Create pointer via instance of http.Client struct with default values
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch jobs, status code: %d", resp.StatusCode)
	}

	var jobsResponse JobsResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobsResponse); err != nil {
		return nil, err
	}
	return &jobsResponse, nil
}

/* CheckExistJobs method
 * Input: a single job string or a file path of jobs
 */
func (c *Client) CheckExistJobs(input string) ([]ExistJob, error) {
	var jobArr []string

	if file, err := os.Open(input); err == nil {
		// CASE: Input is a file
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				jobArr = append(jobArr, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	} else {
		// CASE: Input is a single job
		jobArr = append(jobArr, input)
	}

	var results []ExistJob
	for _, jobName := range jobArr {
		endpoint := fmt.Sprintf("/job/%s/api/json", jobName)
		req, err := c.newRequest("GET", endpoint)
		if err != nil {
			results = append(results, ExistJob{Name: jobName, Exist: false})
			continue
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			results = append(results, ExistJob{Name: jobName, Exist: false})
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			results = append(results, ExistJob{Name: jobName, Exist: true})
		} else {
			results = append(results, ExistJob{Name: jobName, Exist: false})
		}
	}

	return results, nil
}
