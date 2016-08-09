package core

import (
  "fmt"
  "io"
  "encoding/json"
  "github.com/docker/engine-api/types"
  "github.com/docker/engine-api/types/events"
  "github.com/docker/engine-api/types/filters"
  "golang.org/x/net/context"
  "time"
)

func initEvents() {
  agent.eventStreamReading=false
  for {
    if (!agent.eventStreamReading) {
      fmt.Println("Opening docker events stream...")
      args := filters.NewArgs()
      args.Add("type", "container")
      args.Add("event", "die")
      args.Add("event", "stop")
      args.Add("event", "destroy")
      args.Add("event", "kill")
      args.Add("event", "create")
      args.Add("event", "start")
      eventsOptions := types.EventsOptions{ Filters: args }
      stream, err := agent.client.Events(context.Background(), eventsOptions)
      if (err!=nil) {
        fmt.Printf("docker event stream open error: %v\n", err)
      } else {
        startEventStream(stream)
      }
    }
    time.Sleep(3 * time.Second)
  }
}

func startEventStream(stream io.ReadCloser) {
  dec := json.NewDecoder(stream)
  agent.eventStreamReading=true
  fmt.Println("start events stream reader")
  go func() {
    for {
      var event events.Message
      err := dec.Decode(&event)
      if err != nil {
        fmt.Println(err)
        agent.eventStreamReading=true
        stream.Close()
        return;
      }
      //fmt.Printf("event status:%s type:%s, action:%s Actor:%s\n", event.Status, event.Type, event.Action, event.Actor)
      fmt.Printf("action=%s containerId=%s\n", event.Action, event.Actor.ID)
      agent.updateContainerMap(event.Action, event.Actor.ID)
    }
  }()
}
