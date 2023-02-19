package main

import (
	"github.com/riley-martine/sundial/cmd"
)

// Set by goreleaser:
// https://goreleaser.com/cookbooks/using-main.version/?h=version
var version = "dev"

func main() {
	cmd.Execute(version)
}
