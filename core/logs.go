package core

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/docker/engine-api/types"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const elasticSearchTimeIDQuery = `{"query":{"match":{"container_id":"[container_id]"}},"sort":{"time_id":{"order":"desc"}},"from":0,"size":1}`

func updateLogsStream() {
	//if kafka.kafkaReady {
	for ID, data := range agent.containers {
		if data.logsStream == nil || data.logsReadError {
			lastTimeID := getLastTimeID(ID)
			if lastTimeID == "" {
				fmt.Printf("open logs stream from the begining on container %s\n", ID)
			} else {
				fmt.Printf("open logs stream from time_id=%s on container %s\n", lastTimeID, ID)
			}
			stream, err := openLogsStream(ID, lastTimeID)
			if err != nil {
				fmt.Printf("Error opening logs stream on container: %s\n", ID)
			} else {
				data.logsStream = stream
				startReadingLogs(ID, data)
			}
		}
	}
	//}
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
	return agent.client.ContainerLogs(context.Background(), ID, containerLogsOptions)
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
	return body[ll+delim1+9 : ll+delim1+9+delim2]
}

func startReadingLogs(ID string, data *ContainerData) {
	go func() {
		stream := data.logsStream
		serviceName := data.labels["com.docker.swarm.service.name"]
		serviceID := data.labels["com.docker.swarm.service.id"]
		nodeID := data.labels["com.docker.swarm.node.id"]
		reader := bufio.NewReader(stream)
		fmt.Printf("start reading logs on container: %s\n", ID)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("stop reading log on container %s: %v\n", ID, err)
				data.logsReadError = true
				stream.Close()
				return
			}
			var slog string
			if len(line) > 39 {
				slog = strings.TrimSuffix(line[39:], "\n")
				ntime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", line[8:38])
				if conf.kafka != "" {
					logEntry := logs.LogEntry{
						ServiceName: serviceName,
						ServiceId:   serviceID,
						NodeId:      nodeID,
						ContainerId: ID,
						Message:     slog,
						Timestamp:   ntime.String(),
						TimeId:      line[8:38],
					}
					//if kafka.kafkaReady {
					//	kafka.sendLog(mes)
					//} else {
					//	fmt.Printf("Kafka not ready anymore, stop reading log on container %s\n", ID)
					//	data.logsReadError = true
					//	stream.Close()
					//	return
					//}
					encoded, err := proto.Marshal(&logEntry)
					if err != nil {
						log.Printf("error marshalling log entry: %v", err)
					}
					_, err = nats.Publish("amp-logs", encoded)
					if err != nil {
						log.Printf("error sending log entry: %v", err)
					}
				}
			} else {
				fmt.Printf("invalid log: [%s]\n", line)
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
