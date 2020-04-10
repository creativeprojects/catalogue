package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/creativeprojects/catalogue/cmd"
)

func main() {
	log.SetHandler(cli.Default)
	cmd.Execute()
}
