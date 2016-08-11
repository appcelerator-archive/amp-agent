package core

import (
    "fmt"
    "net/http"
    "encoding/json"

)

type APIContainer struct {
  ContainerId string
  ServiceName string
  ServiceId string
  State string
  Health string
}

const baseUrl="/api/v1"

//Start API server
func initAPI() {
    fmt.Println("Start API server on port "+conf.apiPort)
    go func() {
      http.HandleFunc(baseUrl+"/health", agentHealth)
      http.HandleFunc(baseUrl+"/containers", getHandledContainers)
      http.ListenAndServe(":"+conf.apiPort, nil)
    }()
}

//for HEALTHCHECK Dockerfile instruction
func agentHealth(resp http.ResponseWriter, req *http.Request) {
  if (agent.eventStreamReading) {
    resp.WriteHeader(200)
  } else {
    resp.WriteHeader(400)
  }
}

//return the running container list with their paremeter including health
func getHandledContainers(resp http.ResponseWriter, req *http.Request) {
  containers := make([]APIContainer, len(agent.containers))
  var nn int=0
  for key, data := range agent.containers {
    containers[nn]= APIContainer {
      ContainerId: key,
      ServiceName: data.labels["com.docker.swarm.service.name"],
      ServiceId: data.labels["com.docker.swarm.service.id"],
      State: data.state,
      Health: data.health,
    }
    nn++
  }
  json.NewEncoder(resp).Encode(containers)
}
