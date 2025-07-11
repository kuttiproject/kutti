package main

// On Windows, we will include the following drivers:
//  - hyperv
//	- vbox

import (
	_ "github.com/kuttiproject/driver-hyperv"
	_ "github.com/kuttiproject/driver-vbox"
	"github.com/kuttiproject/kutti/internal/pkg/cli"
)

func init() {
	// The default driver will be:
	//	- hyperv
	_, ok := cli.Default("driver")
	if !ok {
		cli.SetDefault("driver", "hyperv")
	}
}
