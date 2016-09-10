package core

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

//Verify if the event stream is working, if not start it
func updateEventsStream() {
	if !agent.eventStreamReading {
		fmt.Println("Opening docker events stream...")
		args := filters.NewArgs()
		args.Add("type", "container")
		args.Add("event", "die")
		args.Add("event", "stop")
		args.Add("event", "destroy")
		args.Add("event", "kill")
		args.Add("event", "create")
		args.Add("event", "start")
		eventsOptions := types.EventsOptions{Filters: args}
		stream, err := agent.dockerClient.Events(context.Background(), eventsOptions)
		agent.eventsStream = stream
		if err != nil {
			fmt.Printf("docker event stream open error: %v\n", err)
		} else {
			startEventStream(stream)
		}
	}
}

//Start and read the docker event stream, send events to kafka and update container list accordingly
func startEventStream(stream io.ReadCloser) {
	dec := json.NewDecoder(stream)
	agent.eventStreamReading = true
	fmt.Println("start events stream reader")
	go func() {
		for {
			var event events.Message
			err := dec.Decode(&event)
			if err != nil {
				fmt.Println(err)
				agent.eventStreamReading = false
				stream.Close()
				return
			}
			fmt.Printf("Docker event: action=%s containerId=%s\n", event.Action, event.Actor.ID)
			agent.updateContainerMap(event.Action, event.Actor.ID)
			//if conf.kafka != "" && kafka.kafkaReady {
			//	kafka.sendEvent(event)
			//}
			// TODO: send to NATS
		}
	}()
}
