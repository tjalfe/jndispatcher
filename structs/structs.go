package structs

import "crypto/x509"

type KafkaCertificate struct {
	Certificate      *x509.Certificate
	PrivateKey       any
	CertificateChain []*x509.Certificate
}

type KafkaTarget struct {
	KafkaServers     []string
	Certificate      *x509.Certificate
	PrivateKey       any
	CertificateChain []*x509.Certificate
	CACertPool       *x509.CertPool
}

type KafkaTopic struct {
	TopicName        string
	KafkaCertificate *KafkaCertificate
	KafkaTarget      *KafkaTarget
}
