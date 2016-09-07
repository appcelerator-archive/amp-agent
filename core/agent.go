package core

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

//Agent data
type Agent struct {
	dockerClient       *client.Client
	containers         map[string]*ContainerData
	eventsStream       io.ReadCloser
	eventStreamReading bool
	lastUpdate         time.Time
	kafkaClient        Kafka
	kafkaProducer      sarama.AsyncProducer
}

//ContainerData data
type ContainerData struct {
	name          string
	labels        map[string]string
	state         string
	health        string
	logsStream    io.ReadCloser
	logsReadError bool
}

var agent Agent

//AgentInit Connect to docker engine, get initial containers list and start the agent
func AgentInit(version string, build string) error {
	runtime.GOMAXPROCS(50)
	agent.trapSignal()
	conf.init(version, build)
	err := agent.kafkaClient.Connect(conf.kafkaHost)
	if err != nil {
		return err
	}
	log.Println("Connected to Kafka")
	agent.kafkaProducer, err = agent.kafkaClient.NewAsyncProducer()
	if err != nil {
		return err
	}
	log.Println("Kafka producer successfuly created")

	topicFound := false

WaitForTopic:
	for i := 0; i < 60; i++ {
		topics, err := agent.kafkaClient.Topics()
		if err != nil {
			return err
		}
		for _, topic := range topics {
			if topic == kafkaLogsTopic {
				log.Println("Kafka logs topic available: ", kafkaLogsTopic)
				topicFound = true
				break WaitForTopic
			}

		}
		time.Sleep(1 * time.Second)
	}

	if !topicFound {
		log.Fatalln("Kafka topic not available.")
	}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(conf.dockerEngine, "v1.24", nil, defaultHeaders)
	if err != nil {
		agent.kafkaClient.Close()
		return err
	}
	agent.dockerClient = cli

	fmt.Println("Connected to Docker-engine")

	fmt.Println("Extracting containers list...")
	agent.containers = make(map[string]*ContainerData)
	ContainerListOptions := types.ContainerListOptions{All: false}
	containers, err := agent.dockerClient.ContainerList(context.Background(), ContainerListOptions)
	if err != nil {
		agent.kafkaClient.Close()
		return err
	}
	for _, cont := range containers {
		agent.addContainer(cont.ID)
	}
	agent.lastUpdate = time.Now()
	fmt.Println("done")
	agent.start()
	return nil
}

//Main agent loop, verify if events and logs stream are started if not start them
func (agt *Agent) start() {
	initAPI()
	for {
		updateLogsStream()
		updateEventsStream()
		time.Sleep(time.Duration(conf.period) * time.Second)
	}
}

//Update containers list concidering event action and event containerId
func (agt *Agent) updateContainerMap(action string, containerID string) {
	if action == "start" {
		agt.addContainer(containerID)
	} else if action == "destroy" || action == "die" || action == "kill" || action == "stop" {
		go func() {
			time.Sleep(5 * time.Second)
			agt.removeContainer(containerID)
		}()
	}
}

//add a container to the main container map and retrieve some container information
func (agt *Agent) addContainer(ID string) {
	_, ok := agt.containers[ID]
	if !ok {
		inspect, err := agt.dockerClient.ContainerInspect(context.Background(), ID)
		if err == nil {
			data := ContainerData{
				name:          inspect.Name,
				labels:        inspect.Config.Labels,
				state:         inspect.State.Status,
				health:        "",
				logsStream:    nil,
				logsReadError: false,
			}
			if inspect.State.Health != nil {
				data.health = inspect.State.Health.Status
			}
			fmt.Println("add container", data.name)
			agt.containers[ID] = &data
		} else {
			fmt.Printf("Container inspect error: %v\n", err)
		}
	}
}

//Suppress a container from the main container map
func (agt *Agent) removeContainer(ID string) {
	data, ok := agent.containers[ID]
	if ok {
		fmt.Println("remove container", data.name)
		delete(agt.containers, ID)
	}
}

func (agt *Agent) updateContainer(ID string) {
	data, ok := agt.containers[ID]
	if ok {
		inspect, err := agt.dockerClient.ContainerInspect(context.Background(), ID)
		if err == nil {
			data.labels = inspect.Config.Labels
			data.state = inspect.State.Status
			data.health = ""
			if inspect.State.Health != nil {
				data.health = inspect.State.Health.Status
			}
			fmt.Println("update container", data.name)
		} else {
			fmt.Printf("Container %s inspect error: %v\n", data.name, err)
		}
	}
}

// Launch a routine to catch SIGTERM Signal
func (agt *Agent) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("\namp-agent received SIGTERM signal")
		agt.eventsStream.Close()
		closeLogsStreams()
		//kafka.close()
		agt.kafkaClient.Close()
		os.Exit(1)
	}()
}
