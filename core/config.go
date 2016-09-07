package core

import (
	"fmt"
	"os"
	"strconv"
)

//AgentConfig configuration parameters
type AgentConfig struct {
	dockerEngine     string
	kafkaHost        string
	elasticsearchURL string
	apiPort          string
	period           int
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
	cfg.kafkaHost = "kafka:9092"
	cfg.elasticsearchURL = "elasticsearch:9200/amp-logs/_search"
	cfg.apiPort = "3000"
	cfg.period = 10
}

//Update config with env variables
func (cfg *AgentConfig) loadConfigUsingEnvVariable() {
	cfg.dockerEngine = getStringParameter("DOCKER", cfg.dockerEngine)
	cfg.kafkaHost = getStringParameter("KAFKA_HOST", cfg.kafkaHost)
	cfg.apiPort = getStringParameter("API_PORT", cfg.apiPort)
	cfg.elasticsearchURL = getStringParameter("ELASTICSEARCH", cfg.elasticsearchURL)
	cfg.period = getIntParameter("PERIOD", cfg.period)
}

//display amp-pilot configuration
func (cfg *AgentConfig) displayConfig(version string) {
	fmt.Printf("amp-agent version: %v\n", version)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("Docker-engine: %s\n", conf.dockerEngine)
	fmt.Printf("Kafka host: %s\n", conf.kafkaHost)
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
