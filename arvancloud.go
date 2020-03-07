package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

// ArvanBaseURL is the Base URL of ArvanCloud ECC Service
const ArvanBaseURL string = "https://napi.arvancloud.com/ecc/v1/regions"


const (
	GetServerPath = "/servers/%s"
	UploadSSHKeyPath = "/ssh-keys"
	RemoveSSHKeyPath = "/ssh-keys/%s"
	CreateServerPath = "/servers"
	StartServerPath = "/servers/%s/power-on"
	StopServerPath = "/servers/%s/power-off"
	RestartServerPath = "/servers/%s/reboot"
	RemoveServerPath = "/servers/%s"
)


// Client struct
type Client struct {
	APIToken string
	Region   string
	BaseURL  string
}

// SSHKey struct
type SSHKey struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

// ServerSecurityGroup struct
type ServerSecurityGroup struct {
	Name string `json:"name"`
}

// ServerRequest struct
type ServerRequest struct {
	Image          string                `json:"image_id"`
	Flavor         string                `json:"flavor_id"`
	SSHKeyName     string                `json:"key_name"`
	Name           string                `json:"name"`
	Network        string                `json:"network_id"`
	SecurityGroups []ServerSecurityGroup `json:"security_groups"`
	SSHKey         bool                  `json:"ssh_key"`
	Count          int                   `json:"count"`
}

// ServerResponseData struct
type ServerResponseData struct {
	ID string `json:"id"`
}

// ServerResponse struct
type ServerResponse struct {
	Data ServerResponseData `json:"data"`
}

// Server struct
type Server struct {
	ID        string `json:"id,omitempty"`
	IPAddress string `json:"addr,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
}

func (c *Client) getURL(format string, args ...interface{}) string {
	url := c.BaseURL + fmt.Sprintf(format, args...)
	return url
}

// NewClient returns an API client struct
func NewClient(apitoken, region string) *Client {
	return &Client{
		APIToken: apitoken,
		Region:   region,
		BaseURL:  fmt.Sprintf(ArvanBaseURL+"/%s", region),
	}
}

// DoRequest executes http requests
func (c *Client) DoRequest(req *http.Request, status int) ([]byte, error) {

	req.Header.Set("authorization", fmt.Sprintf("Apikey %s", c.APIToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if status != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

// GetServer returns information about an instance
func (c *Client) GetServer(serverId string) (*Server, error) {
	url := c.getURL(GetServerPath, serverId)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	data, err := c.DoRequest(req, 200)

	if err != nil {
		return nil, err
	}

	Server := Server{
		ID:        gjson.Get(string(data), "data.id").String(),
		IPAddress: gjson.Get(string(data), "data.addresses.public1.0.addr").String(),
		Name:      gjson.Get(string(data), "data.name").String(),
		Status:    gjson.Get(string(data), "data.status").String(),
	}

	return &Server, nil

}

// UploadSSHKey uploads an SSH key to ArvanCloud
func (c *Client) UploadSSHKey(sshkey *SSHKey) error {
	url := c.getURL(UploadSSHKeyPath)

	j, err := json.Marshal(sshkey)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))

	if err != nil {
		return err
	}
	_, err = c.DoRequest(req, 201)
	return err
}

// RemoveSSHKey removes an SSH key from ArvanCloud
func (c *Client) RemoveSSHKey(keyId string) error {
	url := c.getURL(RemoveSSHKeyPath, keyId)

	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}
	_, err = c.DoRequest(req, 200)
	return err
}

// CreateServer creates a server instance and returns the ID
func (c *Client) CreateServer(server *ServerRequest) (*ServerResponse, error) {
	url := c.getURL(CreateServerPath)

	j, err := json.Marshal(server)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))

	if err != nil {
		return nil, err
	}
	data, err := c.DoRequest(req, 201)

	if err != nil {
		return nil, err
	}

	var ServerResponse ServerResponse
	err = json.Unmarshal(data, &ServerResponse)
	if err != nil {
		return nil, err
	}

	return &ServerResponse, nil
}

// StartServer starts the instance
func (c *Client) StartServer(serverId string) error {
	url := c.getURL(StartServerPath, serverId)
	
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return err
	}
	_, err = c.DoRequest(req, 202)
	return err
}

// StopServer stops the instance
func (c *Client) StopServer(serverId string) error {
	url := c.getURL(StopServerPath, serverId)
	
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return err
	}
	_, err = c.DoRequest(req, 202)
	return err
}

// RestartServer reboots the instance
func (c *Client) RestartServer(serverId string) error {
	url := c.getURL(RestartServerPath, serverId)
	
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return err
	}
	_, err = c.DoRequest(req, 202)
	return err
}

//RemoveServer removes an instance
func (c *Client) RemoveServer(serverId string) error {
	url := c.getURL(RemoveServerPath, serverId)
	
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}
	_, err = c.DoRequest(req, 200)
	return err
}
