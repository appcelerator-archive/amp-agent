package main

import (
	"github.com/appcelerator/amp-agent/core"
	"log"
)

// build vars
var (
        // Version is set with a linker flag (see Makefile)
        Version string

        // Build is set with a linker flag (see Makefile)
        Build string
)

func main() {
	err := core.AgentInit(Version,  Build)
	if err != nil {
		log.Fatal(err)
	}
}
