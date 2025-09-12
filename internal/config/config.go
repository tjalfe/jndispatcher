package config

import "crypto/x509"

type KafkaTarget struct {
	KafkaServers     []string
	Certificate      *x509.Certificate
	CertificateChain []*x509.Certificate
	CACertPool       *x509.CertPool
}
