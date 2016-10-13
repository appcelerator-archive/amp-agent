package main

import (
	"github.com/appcelerator/amp-agent/core"
	"log"
	"net/http"
	"os"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "healthcheck" {
		if !healthcheck() {
			os.Exit(1)
		} 
		os.Exit(0)
	}
	err := core.AgentInit(Version, Build)
	if err != nil {
		log.Fatal(err)
	}
}

func healthcheck() bool {
 	response, err := http.Get("http://127.0.0.1:3000/api/v1/health")
        if err != nil {
               return false
        } 
	if response.StatusCode == 200 {
		return true
	} 
        return false
}
