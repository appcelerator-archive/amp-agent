package core

import (
  "os"
  "fmt"
  "strconv"
)

//Json format of conffile
type AgentConfig struct {
  dockerEngine string
  kafka string
  kafkaLogsTopic string
  kafkaDockerEventsTopic string
  port string
  period int
}

var conf AgentConfig

//Load Json conffile and instanciate new Config
func (self *AgentConfig) init(version string) {
  self.setDefault()
  self.loadConfigUsingEnvVariable()
  self.displayConfig(version)
}

//Set default value of configuration
func (self *AgentConfig) setDefault() {
  self.dockerEngine="unix:///var/run/docker.sock"
  self.kafka= "kafka:9092"
  self.kafkaLogsTopic="amp-logs"
  self.kafkaDockerEventsTopic="amp-docker-events"
  self.port="3000"
  self.period=10
}

//Update config with env variables
func (self *AgentConfig) loadConfigUsingEnvVariable() {
  self.dockerEngine = getStringParameter("AMPAGENT_DOCKER", self.dockerEngine)
  self.kafka = getStringParameter("AMPAGENT_KAFKA", self.kafka)
  self.kafkaLogsTopic = getStringParameter("AMPAGENT_LOGS_TOPIC", self.kafkaLogsTopic)
  self.kafkaDockerEventsTopic = getStringParameter("AMPAGENT_DOCKEREVENTS_TOPIC", self.kafkaDockerEventsTopic)
  self.port = getStringParameter("AMPAGENT_PORT", self.port)
  self.period = getIntParameter("AMPAGENT_PERIOD", self.period)
}

//display amp-pilot configuration
func (self * AgentConfig) displayConfig(version string) {
  fmt.Printf("amp-agent version: %v\n", version)
  fmt.Println("----------------------------------------------------------------------------")
  fmt.Println("Configuration:")
  fmt.Printf("Docker-engine: %s\n", conf.dockerEngine)
  fmt.Printf("Kafka addr: %s\n", conf.kafka)
  fmt.Printf("Kafka log topic: %s\n", conf.kafkaLogsTopic)
  fmt.Printf("Kafka docker events topic: %s\n", conf.kafkaDockerEventsTopic)
  fmt.Println("----------------------------------------------------------------------------")
}

//return env variable value, if empty return default value
func getStringParameter(envVariableName string, def string) string {
  value := os.Getenv(envVariableName)
  if value == "" {
    return def
  }
  return value;
}

//return env variable value convert to int, if empty return default value
func getIntParameter(envVariableName string, def int) int {
    value := os.Getenv(envVariableName)
    if value != "" {
        ivalue, err := strconv.Atoi(value)
        if err != nil {
            return def
        }
        return ivalue
    } else {
        return def
    }
}
