package main

import (
	"github.com/nyudlts/fstat/cmd"
)

func main() {
	cmd.Execute()
	cmd.Walk()
	cmd.ShutDown()
}


