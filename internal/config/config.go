package config

import "fmt"

type KafkaClusterConfig struct {
	Name                              string   `yaml:"name"`
	KafkaServers                      []string `yaml:"kafka-servers"`
	KafkaServerCA                     string   `yaml:"kafka-server-ca"`
	KafkaAuthCertificateStore         string   `yaml:"kafka-auth-certificate-store"`
	KafkaAuthCertificateStorePassword string   `yaml:"kafka-auth-certificate-store-password"`
}

func ReadConfig() ([]KafkaClusterConfig, error) {
	// Read dispatcher config file
	routes, err := readDispatcherConfigs()
	if err != nil {
		return nil, fmt.Errorf("error reading dispatcher configs: %w", err)
	}

	for _, route := range routes.MessageRoute {
		fmt.Printf("Route for message type: %s\n", route.MessageType)
	}

	// Read config file containing Kafka cluster definitions
	//	KafkaConfigs, err := getKafkaClients()
	_, err = getKafkaClients()
	if err != nil {
		return nil, fmt.Errorf("error reading Kafka configs: %w", err)
	}
	return nil, nil
}
