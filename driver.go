package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

// Default variables for nl-ams-su1 region
const (
	defaultImage         = "285bcbf1-738b-4bcf-a9a9-b9940587a026" // Ubuntu 18.04
	defaultRegion        = "nl-ams-su1"                           // Serverius/ Amsterdam-Netherlands
	defaultServerFlavor  = "ar-2-1-15"                            // Smallest flavor
	defaultNetwork       = "fe9645fc-2234-4865-895b-e3bb4bb0eb7b" // Public1
	defaultSecurityGroup = "771874f3-541e-4693-97ad-d585e78999ef"
	defaultSSHUser       = "ubuntu"
)

// Driver struct
type Driver struct {
	*drivers.BaseDriver

	APIToken      string
	Image         string
	Region        string
	ServerID      string
	ServerFlavor  string
	Network       string
	SecurityGroup string
	SSHKeyID      string
}

// NewDriver returns a Driver struct
func NewDriver() *Driver {
	return &Driver{
		Image:         defaultImage,
		Region:        defaultRegion,
		ServerFlavor:  defaultServerFlavor,
		Network:       defaultNetwork,
		SecurityGroup: defaultSecurityGroup,

		BaseDriver: &drivers.BaseDriver{
			SSHUser: defaultSSHUser,
			SSHPort: drivers.DefaultSSHPort,
		},
	}
}

// getClient returns a new API client
func (d *Driver) getClient() *Client {
	return NewClient(d.APIToken, d.Region)
}

func (d *Driver) createServer() error {

	server := ServerRequest{
		Image:      d.Image,
		Flavor:     d.ServerFlavor,
		SSHKeyName: d.SSHKeyID,
		Name:       d.GetMachineName(),
		Network:    d.Network,
		SecurityGroups: []ServerSecurityGroup{
			{Name: d.SecurityGroup},
		},
		SSHKey: true,
		Count:  1,
	}

	Server, err := d.getClient().CreateServer(&server)
	if err != nil {
		return err
	}

	d.ServerID = Server.Data.ID

	for {
		serverState, err := d.GetState()
		if err != nil {
			return err
		}

		if serverState == state.Running {
			break
		}

		time.Sleep(1 * time.Second)
	}

	ServerData, err := d.getServer()
	if err != nil {
		return err
	}

	d.IPAddress = ServerData.IPAddress

	return nil
}

// Create Server
func (d *Driver) Create() error {

	// Create SSH key
	if err := d.createSSHKey(); err != nil {
		return err
	}

	// Create Server
	err := d.createServer()

	return err
}

func (d *Driver) createSSHKey() error {
	// Generate new SSH Key pair
	err := ssh.GenerateSSHKey(d.GetSSHKeyPath())
	if err != nil {
		return err
	}

	// Read public SSH Key
	publicKey, err := ioutil.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return err
	}

	sshkey := SSHKey{
		Name:      d.GetMachineName(),
		PublicKey: string(publicKey),
	}

	err = d.getClient().UploadSSHKey(&sshkey)

	d.SSHKeyID = d.GetMachineName()

	return err
}

// GetSSHUsername ...
func (d *Driver) GetSSHUsername() string {
	if d.SSHUser == "" {
		d.SSHUser = defaultSSHUser
	}

	return d.SSHUser
}

func (d *Driver) getServer() (*Server, error) {
	server, err := d.getClient().GetServer(d.ServerID)

	return server, err
}

// Start Server
func (d *Driver) Start() error {
	err := d.getClient().StartServer(d.ServerID)
	return err
}

// Stop Server
func (d *Driver) Stop() error {
	err := d.getClient().StopServer(d.ServerID)
	return err
}

// Restart Server
func (d *Driver) Restart() error {
	err := d.getClient().RestartServer(d.ServerID)
	return err
}

// Remove Server
func (d *Driver) Remove() error {
	err := d.getClient().RemoveServer(d.ServerID)
	if err != nil {
		return err
	}

	err = d.getClient().RemoveSSHKey(d.SSHKeyID)

	return err
}

// Kill Server
func (d *Driver) Kill() error {
	err := d.getClient().StopServer(d.ServerID)
	return err
}

// GetState ...
func (d *Driver) GetState() (state.State, error) {
	server, err := d.getServer()
	if err != nil {
		return state.Error, err
	}

	switch server.Status {
	case "build":
		return state.Starting, nil
	case "active":
		return state.Running, nil
	case "stop":
		return state.Stopped, nil
	}

	return state.None, nil
}

// GetSSHHostname ...
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetURL ...
func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

// DriverName ...
func (d *Driver) DriverName() string {
	return "arvan"
}

// GetCreateFlags ...
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "ARVAN_API_TOKEN",
			Name:   "arvan-api-token",
			Usage:  "Api token",
		},
		mcnflag.StringFlag{
			EnvVar: "ARVAN_IMAGE",
			Name:   "arvan-image",
			Usage:  "Image",
			Value:  defaultImage,
		},
		mcnflag.StringFlag{
			EnvVar: "ARVAN_REGION",
			Name:   "arvan-region",
			Usage:  "Region",
			Value:  defaultRegion,
		},
		mcnflag.StringFlag{
			EnvVar: "ARVAN_SERVER_FLAVOR",
			Name:   "arvan-server-flavor",
			Usage:  "Server flavor",
			Value:  defaultServerFlavor,
		},
		mcnflag.StringFlag{
			EnvVar: "ARVAN_NETWORK",
			Name:   "arvan-network",
			Usage:  "Network",
			Value:  defaultNetwork,
		},
		mcnflag.StringFlag{
			EnvVar: "ARVAN_SECURITY_GROUP",
			Name:   "arvan-security-group",
			Usage:  "Security group",
			Value:  defaultSecurityGroup,
		},
		mcnflag.StringFlag{
			EnvVar: "ARVAN_SSH_USER",
			Name:   "arvan-ssh-user",
			Usage:  "SSH username",
			Value:  defaultSSHUser,
		},
	}
}

// SetConfigFromFlags ...
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.APIToken = flags.String("arvan-api-token")
	d.Image = flags.String("arvan-image")
	d.Region = flags.String("arvan-region")
	d.ServerFlavor = flags.String("arvan-server-flavor")
	d.Network = flags.String("arvan-network")
	d.SecurityGroup = flags.String("arvan-security-group")
	d.SSHUser = flags.String("arvan-ssh-user")

	d.SetSwarmConfigFromFlags(flags)

	if d.APIToken == "" {
		return fmt.Errorf("arvan driver requres the --arvan-api-token option")
	}

	return nil
}
