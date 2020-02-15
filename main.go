package main

import "github.com/docker/machine/libmachine/drivers/plugin"

var (
	// Version of the driver
	Version string
)

func main() {
	plugin.RegisterDriver(NewDriver())
}
