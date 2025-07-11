package main

// On Mac OS on Apple silicon, we will include the following drivers:
//	- vbox
//  - lima

import (
	_ "github.com/kuttiproject/driver-lima"
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
