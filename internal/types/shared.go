package types

import _ "embed"

// Create a file named internal/types/pcryptinit.secret
// The file should contain a 32 bytes random bytes. This is used for encryption initialization.
// Sample command to generate such a string: `openssl rand 32 > internal/types/pcryptinit.secret`
// This file should be ignored by git (add it to .gitignore).
// The content of the file will be embedded into the binary at compile time.
// Make sure to keep this file secure and do not share it publicly.
// If the file is lost, and project is recompiled any data encrypted with it cannot be decrypted.
//
//go:embed pcryptinit.secret
var PcryptInit []byte

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
