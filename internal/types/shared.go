package types

const (
	PcryptInit = "1foqrXBga8wUmRkTx@406A^UhLEb5xn%"
)

type Config struct {
	KafkaServers                      []string      `yaml:"kafka-servers"`
	KafkaServerCA                     string        `yaml:"kafka-server-ca"`
	KafkaAuthCertificateStore         string        `yaml:"kafka-auth-certificate-store"`
	KafkaAuthCertificateStorePassword string        `yaml:"kafka-auth-certificate-store-password"`
	InputTopic                        string        `yaml:"input-topic"`
	KafkaConsumerGroup                string        `yaml:"kafka-consumer-group"`
	TrustServerCa                     bool          `yaml:"trust-server-ca"`
	TrustExtraCa                      []string      `yaml:"trust-extra-ca"`
	MessageTypes                      []MessageType `yaml:"message_types"`
}

type MessageType struct {
	MessageType string   `yaml:"message_type"`
	OutputTopic []string `yaml:"output_topic"`
}
