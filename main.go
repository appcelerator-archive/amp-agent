package main

import (
	"github.com/appcelerator/amp-agent/core"
	"log"
)

const version string = "1.0.0-6"

func main() {
	err := core.AgentInit(version)
	if err != nil {
		log.Fatal(err)
	}
}
