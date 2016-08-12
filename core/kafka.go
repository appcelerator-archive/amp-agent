package core

import (
	"encoding/json"
	"fmt"
	samara "github.com/Shopify/sarama"
	"github.com/docker/engine-api/types/events"
	"strings"
	"time"
)

//Kafka status
type Kafka struct {
	producer   samara.AsyncProducer
	kafkaReady bool
	kafkatest  bool
	kafkaInit  bool
}

//Log message mimimum parameter list
type logMessage struct {
	Timestamp   time.Time `json:"timestamp"`
	TimeID      string    `json:"time_id"`
	ServiceID   string    `json:"service_id"`
	ServiceName string    `json:"service_name"`
	ServiceUUID string    `json:"service_uuid"` //obsolet to be removed
	Message     string    `json:"message"`
	ContainerID string    `json:"container_id"`
	NodeID      string    `json:"node_id"`
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
func (kfk *Kafka) startPeriodicKafkaChecking() {
	go func() {
		time.Sleep(30 * time.Second)
		fmt.Println("start Kafka checking")
		for {
			if kfk.kafkaReady {
				kfk.kafkaInit = false
			} else {
				fmt.Println("Kafka is not ready")
				config := samara.NewConfig()
				client, err := samara.NewClient(strings.Split(conf.kafka, ","), config)
				if err == nil {
					prod, err := samara.NewAsyncProducerFromClient(client)
					if err == nil {
						kfk.producer = prod
						mes := logMessage{
							ServiceName: "test",
							ServiceUUID: "test",
							ServiceID:   "test",
							NodeID:      "test",
							ContainerID: "test",
							Message:     "test",
							Timestamp:   time.Now(),
							TimeID:      fmt.Sprintf("%v", time.Now().UnixNano()),
						}
						//kafka is ready, but verifying that topic is created
						kfk.kafkatest = true
						kafka.sendLog(mes)
						time.Sleep(12 * time.Second) //time needed to be sure to get kafka error
						if kfk.kafkatest {
							fmt.Println("Kafka is ready")
							kfk.kafkaReady = true
						}
					}
				}
			}
			if kfk.kafkaInit {
				time.Sleep(10 * time.Second)
			} else {
				time.Sleep(30 * time.Second)
			}

		}
	}()
}

//Marshal the logMessage and send it to Kafka, unmarshal message field if json to distribute it on separate fields
func (kfk *Kafka) sendLog(mes logMessage) {
	var data string
	if len(mes.Message) == 0 || mes.Message[0:1] != "{" {
		dat, _ := json.Marshal(mes)
		data = string(dat)
	} else {
		var message = mes.Message
		mesMap := make(map[string]string)
		var objmap map[string]*json.RawMessage
		err := json.Unmarshal([]byte(message), &objmap)
		if err == nil {
			for key, value := range objmap {
				if key != "timestamp" {
					data, err2 := value.MarshalJSON()
					if err2 == nil {
						mesMap[key] = strings.Trim(string(data), "\"")
					}
				}
			}
			mesValue, ok := mesMap["msg"]
			if ok {
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
	case kfk.producer.Input() <- &samara.ProducerMessage{Topic: conf.kafkaLogsTopic, Key: nil, Value: samara.StringEncoder(data)}:
		break
	case err := <-kfk.producer.Errors():
		if mes.ContainerID != "test" {
			fmt.Println("Kafka not ready anymore, error sending logs: ", err)
		} else {
			kfk.kafkatest = false
		}
		kfk.kafkaReady = false
		break
	}
}

//Marshal the message and send it to Kafka
func (kfk *Kafka) sendEvent(mes events.Message) {
	data, _ := json.Marshal(mes)
	select {
	case kfk.producer.Input() <- &samara.ProducerMessage{Topic: conf.kafkaDockerEventsTopic, Key: nil, Value: samara.StringEncoder(string(data))}:
		break
	case err := <-kfk.producer.Errors():
		fmt.Println("Kafka not ready anymore, error sending event: ", err)
		break
	}
}

//Close Kafka producer
func (kfk *Kafka) close() error {
	if kfk.producer == nil {
		return nil
	}
	if err := kfk.producer.Close(); err != nil {
		return err
	}
	return nil
}
