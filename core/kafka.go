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
    fmt.Println("start Kafka checking")
    go func() {
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
                time.Sleep(time.Duration(3) * time.Second)
            } else {
                time.Sleep(time.Duration(30) * time.Second)
            }
        }
    }()
}

//Marshal the logMessage and send it to Kafka, unmarshal message field if json to distribute it on separate fields 
func (self *Kafka) sendLog(mes logMessage) {
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
            mesValue, ok := mesMap["msg"]
            if (ok) {
                mesMap["message"] = mesValue
            } else {
                mesMap["message"] = mes.Message
            }
            mesMap["service_id"] = mes.Service_id
            mesMap["service_uuid"] = mes.Service_uuid
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
            break
        case err := <-self.producer.Errors():
            if mes.Container_id != "test" {
                fmt.Println("Kafka not ready anymore, error sending message: ", err)
            }
            self.kafkatest = false
            self.kafkaReady = false
            break
    }
}

//Marshal the message and send it to Kafka
func (self *Kafka) sendEvent(mes events.Message) {
    data, _ := json.Marshal(mes)
    select {
        case self.producer.Input() <- &samara.ProducerMessage{Topic: conf.kafkaDockerEventsTopic, Key: nil, Value: samara.StringEncoder(string(data))}:
            //fmt.Println("sent")
            break
        case err := <-self.producer.Errors():
            fmt.Println("Kafka not ready anymore, error sending message: ", err)
            self.kafkaReady = false
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
