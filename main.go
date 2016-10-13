package main

import (
	"github.com/appcelerator/amp-agent/core"
	"log"
	"net/http"
	"os"
	"fmt"
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
			fmt.Println("ko")
			os.Exit(1)
		} 
		fmt.Println("ok")
		os.Exit(0)
	}
	err := core.AgentInit(Version, Build)
	if err != nil {
		log.Fatal(err)
	}
}

func healthcheck() bool {
 	response, err := http.Get("http://localhost:3000/api/v1/health")
 	fmt.Println(err)
        if err != nil {
               return false
        } 
        fmt.Println(response.StatusCode)
	if response.StatusCode == 200 {
		return true
	} 
        return false
}
