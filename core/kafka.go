package core

import (
    "time"
    "fmt"
    "strings"
    "encoding/json"
    samara "github.com/Shopify/sarama"
)

const maxBuffer = 10000

type Kafka struct {
    producer samara.AsyncProducer
    kafkaReady bool
    kafkaInit bool
    messageBuffer []logMessage
    messageBufferIndex int
}

type logMessage struct {
    Timestamp string            `json:"timestamp"`
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
    kafka.messageBuffer = make([]logMessage, maxBuffer)
    kafka.messageBufferIndex = 0
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
        fmt.Println("1", err)
        if (err == nil) {
          prod, err := samara.NewAsyncProducerFromClient(client)
          fmt.Println("2", err)
          if err == nil {
              fmt.Println("Kafka producer ready on topic: "+conf.kafkaLogsTopic)
              self.producer = prod
              self.kafkaReady = true
              self.kafkaInit = false
              self.sendMessageBuffer()
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
func (self *Kafka) sendMessage(mes logMessage) {
    if !self.kafkaReady {
        self.saveMessageOnBuffer(mes)
    } else {
        self.sendToKafka(mes)
    }
}

//Marshal the message and send it to Kafka
func (self *Kafka) sendToKafka(mes logMessage) {
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
            mesMap["timestamp"] = mes.Timestamp
            mesMap["time_id"] = mes.Timestamp
            dat, _ := json.Marshal(mesMap)
            data = string(dat)
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
            self.saveMessageOnBuffer(mes)
            break
    }
}

//Save the message struct in Buffer
func (self *Kafka) saveMessageOnBuffer(msg logMessage) {
    if self.messageBufferIndex < maxBuffer {
        self.messageBuffer[self.messageBufferIndex] = msg
        self.messageBufferIndex++
    }
}

//Send all message in buffer to Kafka
func (self *Kafka) sendMessageBuffer() {
    fmt.Printf("Write message buffer to Kafka (%v)\n", self.messageBufferIndex)
    if self.messageBufferIndex >0 {
        for  ii := 0; ii < self.messageBufferIndex; ii++ {
            self.sendToKafka(self.messageBuffer[ii])
        }
    }
    self.messageBufferIndex = 0
    fmt.Println("Write message buffer done")
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
