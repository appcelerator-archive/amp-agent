package core

import (
  "fmt"
  "io"
  "github.com/docker/engine-api/types"
  "golang.org/x/net/context"
  "bufio"
  "strings"
  "time"
)

func updateLogs() {
  for id, data := range agent.containers {
    if (data.logsStream==nil || data.logsReadError) {
      stream, err := openLogsStream(id)
      if (err!=nil) {
        fmt.Printf("Error opening logs stream on container: %s\n", id)
      } else {
        data.logsStream = stream
        startReadingLogs(id, data)
      }
    }
  }
}

func openLogsStream(id string) (io.ReadCloser, error) {
  containerLogsOptions := types.ContainerLogsOptions {
    ShowStdout: true,
    ShowStderr: true,
    Follow: true,
    Timestamps: true,
  }
  return agent.client.ContainerLogs(context.Background(), id, containerLogsOptions)
}

func startReadingLogs(id string, data *ContainerData) {
  fmt.Printf("start reading logs on container: %s\n", id)
  go func() {
    stream := data.logsStream
    serviceName := data.labels["com.docker.swarm.service.name"]
    serviceId := data.labels["com.docker.swarm.service.id"]
    nodeId := data.labels["com.docker.swarm.node.id"]
    reader := bufio.NewReader(stream)
    for {
      line, err :=reader.ReadString('\n')
      if (err!=nil) {
        fmt.Printf("stop reading log on container %s: %v\n", id, err)
        data.logsReadError = true
        return
      }
      slog := strings.TrimSuffix(line[39:], "\n")
      ntime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", line[8:38])
      if (conf.kafka!="") {
        mes := logMessage {
          Service_name : serviceName,
          Service_uuid: serviceName,
          Service_id : serviceId,
          Node_id : nodeId,
          Container_id : id,
          Message : slog,
          Timestamp : ntime,
          Time_id : line[8:38],
        }
        kafka.sendLog(mes)
      }
    }
  }()
}
