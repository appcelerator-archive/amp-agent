package core

import (
    "time"
    "fmt"
    "strings"
    "encoding/json"
    samara "github.com/Shopify/sarama"
    "github.com/docker/engine-api/types/events"
)

type Kafka struct {
    producer samara.AsyncProducer
    kafkaReady bool
    kafkatest bool
    kafkaInit bool
}

//Log message mimimum parameter list
type logMessage struct {
    Timestamp time.Time         `json:"timestamp"`
    Time_id string              `json:"time_id"`
    Service_id string           `json:"service_id"`
    Service_name string         `json:"service_name"`
    Service_uuid string         `json:"service_uuid"` //obsolet to be removed
    Message string              `json:"message"`
    Container_id string         `json:"container_id"`
    Node_id string              `json:"node_id"`
}

var kafka Kafka

//init Kafka instance
func initKafka() {
    fmt.Println("Init Kafka")
    kafka.kafkaReady = false
    kafka.kafkaInit = true
    kafka.startPeriodicKafkaChecking()
}

//launch the periodical Kafka checking trying to create Producer
func (self *Kafka) startPeriodicKafkaChecking() {
    go func() {
        time.Sleep(30 * time.Second)
        fmt.Println("start Kafka checking")
        for {
            if self.kafkaReady {
                self.kafkaInit = false
            } else {
                fmt.Println("Kafka is not ready")
                config := samara.NewConfig()
                client, err := samara.NewClient(strings.Split(conf.kafka,","), config)
                if (err == nil) {
                    prod, err := samara.NewAsyncProducerFromClient(client)
                    if err == nil {
                        self.producer = prod
                        mes := logMessage {
                            Service_name : "test",
                            Service_uuid: "test",
                            Service_id : "test",
                            Node_id : "test",
                            Container_id : "test",
                            Message : "test",
                            Timestamp : time.Now(),
                            Time_id : fmt.Sprintf("%v", time.Now().UnixNano()),
                        }
                        //kafka is ready, but verifying that topic is created
                        self.kafkatest = true
                        kafka.sendLog(mes)
                        time.Sleep(12 * time.Second) //time needed to be sure to get kafka error
                        if (self.kafkatest) {
                            fmt.Println("Kafka is ready")
                            self.kafkaReady = true
                        }
                    }
                }
            }
            if (self.kafkaInit) {
                time.Sleep(10 * time.Second)
            } else {
                time.Sleep(30 * time.Second)
            }

        }
    }()
}

//Marshal the logMessage and send it to Kafka, unmarshal message field if json to distribute it on separate fields 
func (self *Kafka) sendLog(mes logMessage) {
    var data string
    if (len(mes.Message)==0 || mes.Message[0:1] != "{") {
        dat, _ := json.Marshal(mes)
        data = string(dat)
    } else {
        var message string = mes.Message
        mesMap := make(map[string]string)
        var objmap map[string]*json.RawMessage
        err := json.Unmarshal([]byte(message), &objmap)
        if (err == nil) {
            for key, value := range objmap {
                if (key != "timestamp") {
                    data, err2 := value.MarshalJSON()
                    if (err2 == nil) {
                        mesMap[key] = strings.Trim(string(data),"\"")
                    }
                }
            }
            mesValue, ok := mesMap["msg"]
            if (ok) {
                mes.Message = mesValue
            }
            dat, _ := json.Marshal(mes)
            data = string(dat)
            dataExt, _ := json.Marshal(mesMap)
            data = fmt.Sprintf("%s,%s", dataExt[0:len(dataExt)-1], data[1:])
        } else {
            dat, _ := json.Marshal(mes)
            data = string(dat)
        }
    }
    select {
        case self.producer.Input() <- &samara.ProducerMessage{Topic: conf.kafkaLogsTopic, Key: nil, Value: samara.StringEncoder(data)}:
            break
        case err := <-self.producer.Errors():
            if mes.Container_id != "test" {
                fmt.Println("Kafka not ready anymore, error sending logs: ", err)
            } else {
                self.kafkatest = false
            }
            self.kafkaReady = false
            break
    }
}

//Marshal the message and send it to Kafka
func (self *Kafka) sendEvent(mes events.Message) {
    data, _ := json.Marshal(mes)
    select {
        case self.producer.Input() <- &samara.ProducerMessage{Topic: conf.kafkaDockerEventsTopic, Key: nil, Value: samara.StringEncoder(string(data))}:
            break
        case err := <-self.producer.Errors():
            fmt.Println("Kafka not ready anymore, error sending event: ", err)
            break
    }
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
