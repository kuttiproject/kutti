package main

// On Windows, we will includer the following drivers:
//	- vbox

import (
	_ "github.com/kuttiproject/driver-vbox"
	"github.com/kuttiproject/kutti/internal/pkg/cli"
)

func init() {
	// The default driver will be:
	//	- vbox
	_, ok := cli.Default("driver")
	if !ok {
		cli.SetDefault("driver", "vbox")
	}
}
