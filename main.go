package main

import "astrocyte/cmd"

var (
	version = "development"
	debug   bool
)

func main() {
	cmd.Execute()
}
