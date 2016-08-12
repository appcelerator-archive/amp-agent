package main

import (
	"fmt"
	"github.com/appcelerator/amp-agent/core"
)

const version string = "1.0.0-4"

func main() {
	err := core.AgentInit(version)
	if err != nil {
		fmt.Println("Agent init error: ", err)
	}
}
