package main

import (
	"github.com/kuttiproject/kutti/internal/pkg/cmd"
)

// The linker loader will assign the current version string if built
// with the makefile. The value here is a fallback in case "go install"
// is used.
var version = "v0.3.2-beta1-goinstall"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
