package main

import "github.com/adrianrudnik/pushover-cli/cmd"

var (
	// goreleaser ldflags
	version = "dev"
)

func main() {
	cmd.Version = version
	cmd.Execute()
}
