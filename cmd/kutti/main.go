package main

import (
	"github.com/kuttiproject/kutti/internal/pkg/cmd"
)

// The linker loader will assign the current version string.
var version string

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
