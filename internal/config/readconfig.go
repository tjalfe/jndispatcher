package config

import (
	"fmt"
	"os"

	"github.com/tjalfe/jndispatcher/internal/types"

	"gopkg.in/yaml.v3"
)

func configValidate(config types.Config) error {
	// Validate that required fields are present
	if len(config.KafkaServers) == 0 {
		return fmt.Errorf("no Kafka servers defined in config")
	}
	if config.KafkaAuthCertificateStore == "" {
		return fmt.Errorf("no PKCS12 file for Kafka auth certificate defined in config")
	}
	if config.KafkaAuthCertificateStorePassword == "" {
		return fmt.Errorf("no password for %s for PKCS12 Kafka auth certificate defined in config", config.KafkaAuthCertificateStore)
	}
	if config.InputTopic == "" {
		return fmt.Errorf("no input topic defined in config")
	}
	if len(config.MessageTypes) == 0 {
		return fmt.Errorf("no message types defined in config")
	}
	// Check message types
	messageTypeMap := make(map[string]bool)
	for _, mt := range config.MessageTypes {
		// Check that message types are unique
		if _, exists := messageTypeMap[mt.MessageType]; exists {
			return fmt.Errorf("duplicate message type found in config: %s", mt.MessageType)
		}
		messageTypeMap[mt.MessageType] = true
		// Check that each message type has at least one output topic
		if len(mt.OutputTopic) == 0 {
			return fmt.Errorf("no output topics defined for message type: %s", mt.MessageType)
		}
		// Check that output topic names are not equal to input topic name
		for _, ot := range mt.OutputTopic {
			if ot == config.InputTopic {
				return fmt.Errorf("output topic %s for message type %s is the same as input topic", ot, mt.MessageType)
			}
		}
	}
	return nil
}

func readConfig() (types.Config, error) {
	var config types.Config
	yamlFile, err := os.ReadFile(ConfigYamlFile)
	if err != nil {
		return config, fmt.Errorf("error reading config YAML file: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshaling config YAML: %w", err)
	}
	err = configValidate(config)
	if err != nil {
		return config, err
	}
	return config, nil
}
