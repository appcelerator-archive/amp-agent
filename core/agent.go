package core

import (
    "fmt"
    "github.com/docker/engine-api/client"
    "github.com/docker/engine-api/types"
    "golang.org/x/net/context"
    "io"
)

type Agent struct {
  client *client.Client
  containers map[string]*ContainerData
  eventStreamReading bool
}

type ContainerData struct {
  labels map[string]string
  state string
  health string
  logsStream io.ReadCloser
  logsReadError bool
}

var agent Agent

func (self *Agent) start() {
  initAPI()
  initEvents()
  updateLogs()
}

//Connect to docker engine and get initial containers list
func AgentInit(version string) error {
  conf.init(version)
  initKafka()
  fmt.Println("Connecting to docker...")
  defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
  cli, err := client.NewClient(conf.dockerEngine, "v1.24", nil, defaultHeaders)
  if err != nil {
      return  err
  }
  agent.client=cli
  fmt.Println("Extracting containers list...")
  agent.containers=make(map[string]*ContainerData)
  ContainerListOptions := types.ContainerListOptions{All: false}
  containers, err := agent.client.ContainerList(context.Background(), ContainerListOptions)
  if err != nil {
      return err
  }
  for _, cont := range containers {
    agent.addContainer(cont.ID)
  }
  fmt.Println("done")
  agent.start()
  return nil
}

func (self *Agent) updateContainerMap(action string, containerId string) {
  if action=="start" {
    self.addContainer(containerId)
  } else if (action!="create") {
    self.removeContainer(containerId)
  }
}

func (self *Agent) addContainer(id string) {
  _, ok := self.containers[id]
  if (!ok) {
    inspect, err := self.client.ContainerInspect(context.Background(), id)
    if err == nil {
      data := ContainerData {
        labels: inspect.Config.Labels,
        state: inspect.State.Status,
        health: "",
        logsStream: nil,
        logsReadError: false,
      }
      if (inspect.State.Health != nil) {
        data.health = inspect.State.Health.Status
      }
      fmt.Println("add container", id)
      self.containers[id] = &data
      updateLogs()
    } else {
      fmt.Printf("Container inspect error: %v\n", err)
    }
  }
}

func (self *Agent) removeContainer(id string) {
  _, ok := agent.containers[id]
  if (ok) {
    fmt.Println("remove container", id)
    delete(self.containers, id)
    updateLogs()
  }
}
