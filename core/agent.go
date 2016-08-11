package core

import (
    "fmt"
    "github.com/docker/engine-api/client"
    "github.com/docker/engine-api/types"
    "golang.org/x/net/context"
    "io"
    "time"
    "runtime"
    "os"
    "os/signal"
    "syscall"

)

type Agent struct {
  client *client.Client
  containers map[string]*ContainerData
  eventsStream io.ReadCloser
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


//Connect to docker engine, get initial containers list and start the agent
func AgentInit(version string) error {
  runtime.GOMAXPROCS(50)
  agent.trapSignal()
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

//Main agent loop, verify if events and logs stream are started if not start them
func (self *Agent) start() {
  initAPI()
  for {
    updateLogsStream()
    updateEventsStream()
    time.Sleep(time.Duration(conf.period) * time.Second)
  }
}

//Update containers list concidering event action and event containerId
func (self *Agent) updateContainerMap(action string, containerId string) {
  if action=="start" {
    self.addContainer(containerId)
  } else if (action!="create") {
    self.removeContainer(containerId)
  }
}

//add a container to the main container map and retrieve some container information
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
    } else {
      fmt.Printf("Container inspect error: %v\n", err)
    }
  }
}

//Suppress a container from the main container map
func (self *Agent) removeContainer(id string) {
  _, ok := agent.containers[id]
  if (ok) {
    fmt.Println("remove container", id)
    delete(self.containers, id)
  }
}

//Launch a routine to catch SIGTERM Signal
func (self *Agent) trapSignal() {
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, os.Interrupt)
    signal.Notify(ch, syscall.SIGTERM)
    go func() {
        <-ch
        fmt.Println("\namp-agent received SIGTERM signal")
        self.eventsStream.Close()
        closeLogsStreams()
        kafka.close()
        os.Exit(1)
    }()
}