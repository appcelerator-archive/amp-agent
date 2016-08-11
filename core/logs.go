package core

import (
  "fmt"
  "io"
  "github.com/docker/engine-api/types"
  "golang.org/x/net/context"
  "bufio"
  "strings"
  "time"
  "net/http"
  "bytes"
  "io/ioutil"
)

const elasticSearchTimeIdQuery=`{"query":{"match":{"container_id":"[container_id]"}},"sort":{"time_id":{"order":"desc"}},"from":0,"size":1}`

func updateLogsStream() {
  if (kafka.kafkaReady) {
    for id, data := range agent.containers {
      if (data.logsStream==nil || data.logsReadError) {
        lastTimeId := getLastTimeId(id)
        if (lastTimeId == "") {
          fmt.Printf("open logs stream from the begining on container %s\n", id)
        } else {
          fmt.Printf("open logs stream from time_id=%s on container %s\n", lastTimeId, id)
        }
        stream, err := openLogsStream(id, lastTimeId)
        if (err!=nil) {
          fmt.Printf("Error opening logs stream on container: %s\n", id)
        } else {
          data.logsStream = stream
          startReadingLogs(id, data)
        }
      }
    }
  }
}

func openLogsStream(id string, lastTimeId string) (io.ReadCloser, error) {
  containerLogsOptions := types.ContainerLogsOptions {
    ShowStdout: true,
    ShowStderr: true,
    Follow: true,
    Timestamps: true,
  }
  if (lastTimeId != "") {
    containerLogsOptions.Since = lastTimeId
  }
  return agent.client.ContainerLogs(context.Background(), id, containerLogsOptions)
}

//Use elasticsearch REST API directly
func getLastTimeId(id string) string {
  request := strings.Replace(elasticSearchTimeIdQuery, "[container_id]", id, 1)
  req, err := http.NewRequest("POST", "http://"+conf.elasticsearchUrl, bytes.NewBuffer([]byte(request)))
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      fmt.Println("request error: ", err)
      return ""
  }
  defer resp.Body.Close()

  body, _ := ioutil.ReadAll(resp.Body)
  return extractTimeId(string(body))
}

//Extract time_id from body without unmarchall the body
func extractTimeId(body string) string {
  ll := strings.Index(body, "time_id")
  if (ll<0) {
    return ""
  }
  delim1 := strings.IndexByte(body[ll+8:], '"')
  if (delim1<0) {
    return ""
  }
  delim2 := strings.IndexByte(body[ll+8+delim1+1:], '"')
  if (delim2<0) {
    return ""
  }
  return body[ll+delim1+9:ll+delim1+9+delim2]
}

func startReadingLogs(id string, data *ContainerData) {
  go func() {
    stream := data.logsStream
    serviceName := data.labels["com.docker.swarm.service.name"]
    serviceId := data.labels["com.docker.swarm.service.id"]
    nodeId := data.labels["com.docker.swarm.node.id"]
    reader := bufio.NewReader(stream)
    fmt.Printf("start reading logs on container: %s\n", id)
    for {
      line, err :=reader.ReadString('\n')
      if (err!=nil) {
        fmt.Printf("stop reading log on container %s: %v\n", id, err)
        data.logsReadError = true
        stream.Close()
        return
      }
      var slog string
      if (len(line)>39) {
        slog = strings.TrimSuffix(line[39:], "\n")
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
          if (kafka.kafkaReady) {
            kafka.sendLog(mes)
          } else {
            fmt.Printf("Kafka not ready anymore, stop reading log on container %s\n", id)
            data.logsReadError = true
            stream.Close()
            return
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
    if (data.logsStream!=nil) {
      data.logsStream.Close()
    }
  }
}