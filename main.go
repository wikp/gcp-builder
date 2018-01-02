package main

import (
	"fmt"
	"github.com/wendigo/gcp-builder/cli"
	"github.com/wendigo/gcp-builder/config"
	"os"
)

func main() {

	cfg, err := config.Get()
	if err != nil {
		exit(err)
	}

	client, err := cli.New(cfg)
	if err != nil {
		exit(err)
	}

	if err := client.Run(); err != nil {
		exit(err)
	}
}

func exit(reason error) {
	fmt.Printf("Command failed due to: %s\n", reason)
	os.Exit(1)
}
