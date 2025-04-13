package jenkins

import (
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

// Client Init struct
type Client struct {
	BaseURL  string `yaml:"base_url"`
	Username string `yaml:"username"`
	Token    string `yaml:"api_token"`
}

// LoadConfig method -- Load jenkins config from file
func LoadConfig(configPath string) (*Client, error) {
	// Read & validate file from path
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close() // defer: make sure the file will be closed when the function ends

	// Decode file content to Client struct
	decoder := yaml.NewDecoder(file)
	var client Client
	if err := decoder.Decode(&client); err != nil {
		return nil, err
	}
	return &client, nil
}

/* newRequest method -- INPUT list
 * 1. FetchAllJob: method: GET -- endpoint: /api/json
 * 2. CheckJobExist: method: GET -- endpoint: /job/{job_name}/api/json
 */
func (c *Client) newRequest(method, endpoint string) (*http.Request, error) {
	var url string = c.BaseURL + endpoint
	req, err := http.NewRequest(method, url, nil) // Empty body
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Token)
	return req, nil
}
