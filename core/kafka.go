package core

import (
    "time"
    "fmt"
    "strings"
    "encoding/json"
    samara "github.com/Shopify/sarama"
    "github.com/docker/engine-api/types/events"
)

const maxBuffer = 10000

type Kafka struct {
    producer samara.AsyncProducer
    kafkaReady bool
    kafkaInit bool
    logBuffer []logMessage
    logBufferIndex int
    eventBuffer []events.Message
    eventBufferIndex int

}

type logMessage struct {
    Timestamp time.Time         `json:"timestamp"`
    Time_id string              `json:time_id`
    Service_id string           `json:"service_id"`
    Service_name string         `json:"service_name"`
    Message string              `json:"message"`
    Container_id string         `json:"container_id"`
    Node_id string              `json:"node_id"`
}

var kafka Kafka

//init Kafka struct
func initKafka() {
    fmt.Println("Init Kafka")
    kafka.kafkaReady = false
    kafka.kafkaInit = true
    kafka.logBuffer = make([]logMessage, maxBuffer)
    kafka.logBufferIndex = 0
    kafka.eventBuffer = make([]events.Message, maxBuffer)
    kafka.eventBufferIndex = 0    
    kafka.startPeriodicKafkaChecking()
}

//launch the periodical Kafka checking and create Producer if ready
func (self *Kafka) startPeriodicKafkaChecking() {
  fmt.Println("start Kafka checking")
  go func() {
    for {
      if !self.kafkaReady {
        fmt.Println("Kafka not ready")
        config := samara.NewConfig()
        fmt.Println(conf.kafka)
        client, err := samara.NewClient(strings.Split(conf.kafka,","), config)
        if (err == nil) {
          prod, err := samara.NewAsyncProducerFromClient(client)
          if err == nil {
              fmt.Println("Kafka producer ready on topic: "+conf.kafkaLogsTopic)
              self.producer = prod
              self.kafkaReady = true
              self.kafkaInit = false
              self.sendLogFromBuffer()
          }
        }
      }
      if (self.kafkaInit) {
          time.Sleep(time.Duration(3) * time.Second)
      } else {
          time.Sleep(time.Duration(30) * time.Second)
      }
    }
  }()
}

//create the message struct save it in buffer if Kafka not reday, send to Kafka if ready,
func (self *Kafka) sendLog(mes logMessage) {
    if !self.kafkaReady {
        self.saveLogToBuffer(mes)
    } else {
        self.produceLog(mes)
    }
}

//Marshal the message and send it to Kafka
func (self *Kafka) produceLog(mes logMessage) {
    var data string
    if (len(mes.Message)>0 && mes.Message[0:1] == "{") {
        mesMap := make(map[string]string)
        var objmap map[string]*json.RawMessage
        err := json.Unmarshal([]byte(mes.Message), &objmap)
        if (err == nil) {
            for key, value := range objmap {
                if (key != "timestamp") {
                    data, err2 := value.MarshalJSON()
                    if (err2 == nil) {
                        mesMap[key] = strings.Trim(string(data),"\"")
                    }
                }
            }
            mesMap["message"] = mesMap["msg"]
            mesMap["service_id"] = mes.Service_id
            mesMap["service_name"] = mes.Service_name
            mesMap["node_id"] = mes.Node_id
            mesMap["container_id"] = mes.Container_id
            //mesMap["timestamp"] = mes.Timestamp
            mesMap["time_id"] = mes.Time_id
            dat, _ := json.Marshal(mesMap)
            data = string(dat)
            data = fmt.Sprintf("{\"timestamp\": %v, %s",mes.Timestamp.Unix()*1000, data[1:])
        } else {
            dat, _ := json.Marshal(mes)
            data = string(dat)    
        }
    } else {
        dat, _ := json.Marshal(mes)
        data = string(dat)
    }
    select {
        case self.producer.Input() <- &samara.ProducerMessage{Topic: conf.kafkaLogsTopic, Key: nil, Value: samara.StringEncoder(data)}:
            //fmt.Println("sent")
            break
        case err := <-self.producer.Errors():
            fmt.Println("Kafka not ready anymore, error sending message: ", err)
            self.kafkaReady = false
            self.saveLogToBuffer(mes)
            break
    }
}

//Save the message struct in Buffer
func (self *Kafka) saveLogToBuffer(msg logMessage) {
    if self.logBufferIndex < maxBuffer {
        self.logBuffer[self.logBufferIndex] = msg
        self.logBufferIndex++
    }
}

//Send all message in buffer to Kafka
func (self *Kafka) sendLogFromBuffer() {
    fmt.Printf("Write logs from buffer to Kafka (%v)\n", self.logBufferIndex)
    if self.logBufferIndex >0 {
        for  ii := 0; ii < self.logBufferIndex; ii++ {
            self.produceLog(self.logBuffer[ii])
        }
    }
    self.logBufferIndex = 0
    fmt.Println("Write logs from buffer done")
}


//create the message struct save it in buffer if Kafka not reday, send to Kafka if ready,
func (self *Kafka) sendEvent(mes events.Message) {
    if !self.kafkaReady {
        self.saveEventToBuffer(mes)
    } else {
        self.produceEvent(mes)
    }
}

//Marshal the message and send it to Kafka
func (self *Kafka) produceEvent(mes events.Message) {
    data, _ := json.Marshal(mes)
    select {
        case self.producer.Input() <- &samara.ProducerMessage{Topic: conf.kafkaDockerEventsTopic, Key: nil, Value: samara.StringEncoder(string(data))}:
            //fmt.Println("sent")
            break
        case err := <-self.producer.Errors():
            fmt.Println("Kafka not ready anymore, error sending message: ", err)
            self.kafkaReady = false
            self.saveEventToBuffer(mes)
            break
    }
}

//Save the message struct in Buffer
func (self *Kafka) saveEventToBuffer(msg events.Message) {
    if self.eventBufferIndex < maxBuffer {
        self.eventBuffer[self.eventBufferIndex] = msg
        self.eventBufferIndex++
    }
}

//Send all message in buffer to Kafka
func (self *Kafka) sendEventFromBuffer() {
    fmt.Printf("Write event from buffer to Kafka (%v)\n", self.eventBufferIndex)
    if self.eventBufferIndex >0 {
        for  ii := 0; ii < self.eventBufferIndex; ii++ {
            self.produceEvent(self.eventBuffer[ii])
        }
    }
    self.eventBufferIndex = 0
    fmt.Println("Write events from buffer done")
}

//Close Kafka producer
func (self *Kafka) close() error {
    if self.producer == nil {
        return nil
    }
    if err := self.producer.Close(); err != nil {
        return err
    }
    return nil
}
