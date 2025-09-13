package structs

import "crypto/x509"

type KafkaCertificate struct {
	Certificate      *x509.Certificate
	PrivateKey       any
	CertificateChain []*x509.Certificate
}


