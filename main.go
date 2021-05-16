package main

import "github.com/ckotzbauer/chekr/cmd"

var (
	version = "master"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	cmd.Execute(version, commit, date, builtBy)
}
