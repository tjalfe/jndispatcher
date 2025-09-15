package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Dispatcher struct {
	DefaultOutputKafka string  `yaml:"default_output_kafka"`
	DefaultInputKafka  string  `yaml:"default_input_kafka"`
	MessageRoute       []Route `yaml:"message_types"`
}

type Route struct {
	MessageType  string        `yaml:"message_type"`
	InputKafka   string        `yaml:"input_kafka,omitempty"`
	OutputTopics []RouteOutput `yaml:"output_topics"`
}

type RouteOutput struct {
	OutputKafka string `yaml:"output_kafka,omitempty"`
	OutputTopic string `yaml:"output_topic"`
}

// ReadKafkaConfigs reads a YAML file and unmarshals it into an array of KafkaClusterConfig structs.
func readDispatcherConfigs() (Dispatcher, error) {
	yamlFile, err := os.ReadFile(DispatcherConf)
	if err != nil {
		return Dispatcher{}, fmt.Errorf("error reading YAML file: %w", err) //Using %w for error wrapping
	}

	var configs Dispatcher

	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		return Dispatcher{}, fmt.Errorf("error unmarshaling YAML: %w", err) //Using %w for error wrapping
	}

	return configs, nil
}
