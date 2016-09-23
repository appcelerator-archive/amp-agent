package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/docker/docker/api/types"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const elasticSearchTimeIDQuery = `{"query":{"match":{"container_id":"[container_id]"}},"sort":{"time_id":{"order":"desc"}},"from":0,"size":1}`

func updateLogsStream() {
	for ID, data := range agent.containers {
		if data.logsStream == nil || data.logsReadError {
			lastTimeID := getLastTimeID(ID)
			if lastTimeID == "" {
				fmt.Printf("open logs stream from the begining on container %s\n", data.name)
			} else {
				fmt.Printf("open logs stream from time_id=%s on container %s\n", lastTimeID, data.name)
			}
			stream, err := openLogsStream(ID, lastTimeID)
			if err != nil {
				fmt.Printf("Error opening logs stream on container: %s\n", data.name)
			} else {
				data.logsStream = stream
				startReadingLogs(ID, data)
			}
		}
	}
}

func openLogsStream(ID string, lastTimeID string) (io.ReadCloser, error) {
	containerLogsOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	}
	if lastTimeID != "" {
		containerLogsOptions.Since = lastTimeID
	}
	return agent.dockerClient.ContainerLogs(context.Background(), ID, containerLogsOptions)
}

//Use elasticsearch REST API directly
func getLastTimeID(ID string) string {
	request := strings.Replace(elasticSearchTimeIDQuery, "[container_id]", ID, 1)
	req, err := http.NewRequest("POST", "http://"+conf.elasticsearchURL, bytes.NewBuffer([]byte(request)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request error: ", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return extractTimeID(string(body))
}

//Extract time_id from body without unmarchall the body
func extractTimeID(body string) string {
	ll := strings.Index(body, "time_id")
	if ll < 0 {
		return ""
	}
	delim1 := strings.IndexByte(body[ll+8:], '"')
	if delim1 < 0 {
		return ""
	}
	delim2 := strings.IndexByte(body[ll+8+delim1+1:], '"')
	if delim2 < 0 {
		return ""
	}
	ret := body[ll+delim1+9 : ll+delim1+9+delim2]
	if len(ret) < 25 {
		return ""
	}
	return ret
}

func startReadingLogs(ID string, data *ContainerData) {
	go func() {
		stream := data.logsStream
		serviceName := data.labels["com.docker.swarm.service.name"]
		serviceID := data.labels["com.docker.swarm.service.id"]
		taskName := data.labels["com.docker.swarm.task.name"]
		taskID := data.labels["com.docker.swarm.task.id"]
		nodeID := data.labels["com.docker.swarm.node.id"]
		stackID := data.labels["io.amp.stack.id"]
		stackName := data.labels["io.amp.stack.name"]
		reader := bufio.NewReader(stream)
		fmt.Printf("start reading logs on container: %s\n", data.name)
		nbErr := 0
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("stop reading log on container %s: %v\n", data.name, err)
				data.logsReadError = true
				stream.Close()
				agent.removeContainer(ID)
				return
			}
			var slog string
			if len(line) <= 39 {
				fmt.Printf("invalid log: [%s]\n", line)
				continue
			}

			slog = strings.TrimSuffix(line[39:], "\n")
			ntime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", line[8:38])

			logEntry := logs.LogEntry{
				ServiceName: serviceName,
				ServiceId:   serviceID,
				TaskName:    taskName,
				TaskId:      taskID,
				StackId:     stackID,
				StackName:   stackName,
				NodeId:      nodeID,
				ContainerId: ID,
				Message:     slog,
				Timestamp:   ntime.Format(time.RFC3339Nano),
				TimeId:      line[8:38], //TODO:have a true timeId
			}

			encoded, err := proto.Marshal(&logEntry)
			if err != nil {
				fmt.Printf("error marshalling log entry: %v", err)
			}

			select {
			case agent.kafkaProducer.Input() <- &sarama.ProducerMessage{Topic: kafkaLogsTopic, Value: sarama.ByteEncoder(encoded)}:
				nbErr = 0
			case err := <-agent.kafkaProducer.Errors():
				fmt.Println("Failed to produce message", err)
				nbErr++
				if nbErr > 20 {
					fmt.Println("Kafka not ready anymore: exit")
					agent.eventsStream.Close()
					closeLogsStreams()
					agent.kafkaClient.Close()
					os.Exit(1)
				}
			}
		}
	}()
}

func closeLogsStreams() {
	for _, data := range agent.containers {
		if data.logsStream != nil {
			data.logsStream.Close()
		}
	}
}
