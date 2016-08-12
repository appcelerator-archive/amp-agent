package core

import (
	"fmt"
	"os"
	"strconv"
)

//AgentConfig configuration parameters
type AgentConfig struct {
	dockerEngine           string
	kafka                  string
	kafkaLogsTopic         string
	kafkaDockerEventsTopic string
	elasticsearchURL       string
	apiPort                string
	period                 int
}

var conf AgentConfig

//update conf instance with default value and environment variables
func (cfg *AgentConfig) init(version string) {
	cfg.setDefault()
	cfg.loadConfigUsingEnvVariable()
	cfg.displayConfig(version)
}

//Set default value of configuration
func (cfg *AgentConfig) setDefault() {
	cfg.dockerEngine = "unix:///var/run/docker.sock"
	cfg.kafka = "kafka:9092"
	cfg.kafkaLogsTopic = "amp-logs"
	cfg.kafkaDockerEventsTopic = "amp-docker-events"
	cfg.elasticsearchURL = "elasticsearch:9200/amp-logs/_search"
	cfg.apiPort = "3000"
	cfg.period = 10
}

//Update config with env variables
func (cfg *AgentConfig) loadConfigUsingEnvVariable() {
	cfg.dockerEngine = getStringParameter("AMPAGENT_DOCKER", cfg.dockerEngine)
	cfg.kafka = getStringParameter("AMPAGENT_KAFKA", cfg.kafka)
	cfg.kafkaLogsTopic = getStringParameter("AMPAGENT_LOGS_TOPIC", cfg.kafkaLogsTopic)
	cfg.kafkaDockerEventsTopic = getStringParameter("AMPAGENT_DOCKEREVENTS_TOPIC", cfg.kafkaDockerEventsTopic)
	cfg.apiPort = getStringParameter("AMPAGENT_PORT", cfg.apiPort)
	cfg.elasticsearchURL = getStringParameter("AMPAGENT_ELASTICSEARCH", cfg.elasticsearchURL)
	cfg.period = getIntParameter("AMPAGENT_PERIOD", cfg.period)
}

//display amp-pilot configuration
func (cfg *AgentConfig) displayConfig(version string) {
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
	return value
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
	}
	return def
}
