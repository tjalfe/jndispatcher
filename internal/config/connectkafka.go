package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/tjalfe/jndispatcher/internal/types"

	"github.com/tjalfe/pcrypt"
	"github.com/twmb/franz-go/pkg/kgo"
	"software.sslmate.com/src/go-pkcs12"
)

func readPkcs12(p12Path string, password string) (interface{}, *x509.Certificate, []*x509.Certificate, error) {
	// Decrypt the password for the PKCS#12 file
	decrypter, err := pcrypt.Init(types.PcryptInit)
	if err != nil {
		return nil, &x509.Certificate{}, nil, err
	}
	decrypted_password, err := pcrypt.Decrypt(decrypter, password)
	if err != nil {
		return nil, &x509.Certificate{}, nil, fmt.Errorf("error decrypting password for %s: %w", p12Path, err)
	}

	// Read the PKCS#12 file
	p12Data, err := os.ReadFile(p12Path)
	if err != nil {
		return nil, &x509.Certificate{}, nil, fmt.Errorf("error reading p12 file %s: %w", p12Path, err)
	}
	privateKey, certificate, certificateChain, err := pkcs12.DecodeChain(p12Data, decrypted_password)
	if err != nil {
		return nil, &x509.Certificate{}, nil, fmt.Errorf("error decoding p12 file %s: %w", p12Path, err)
	}
	return privateKey, certificate, certificateChain, nil
}

func generateRootCAPool(config types.Config) (*x509.CertPool, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("error loading system CA certificates: %w", err)
	}
	if config.KafkaServerCA != "" {
		caCert, err := os.ReadFile(config.KafkaServerCA)
		if err != nil {
			return nil, fmt.Errorf("error reading Kafka server CA certificate %s: %w", config.KafkaServerCA, err)
		}
		if ok := certPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("error appending Kafka server CA certificate %s to pool", config.KafkaServerCA)
		}
	}
	return certPool, nil
}

func InitKafkaClient(config types.Config) (*kgo.Client, error) {
	privateKey, clientCertificate, clientCertificateChain, err := readPkcs12(config.KafkaAuthCertificateStore, config.KafkaAuthCertificateStorePassword)
	if err != nil {
		return nil, fmt.Errorf("error reading PKCS#12 file: %w", err)
	}

	var kafkaClientAuth tls.Certificate
	kafkaClientAuth.PrivateKey = privateKey
	kafkaClientAuth.Certificate = append(kafkaClientAuth.Certificate, clientCertificate.Raw)
	for _, cert := range clientCertificateChain {
		kafkaClientAuth.Certificate = append(kafkaClientAuth.Certificate, cert.Raw)
	}

	rootCA, err := generateRootCAPool(config)
	if err != nil {
		return nil, fmt.Errorf("error generating root CA pool: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{kafkaClientAuth},
		RootCAs:      rootCA,
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(config.KafkaServers...),
		kgo.ConsumeTopics(config.InputTopic),
		kgo.ConsumerGroup(config.KafkaConsumerGroup),
		kgo.DialTLSConfig(tlsConfig),
		kgo.DialTimeout(3 * time.Second),
		kgo.ProducerBatchCompression(kgo.ZstdCompression(), kgo.SnappyCompression(), kgo.NoCompression()),
		kgo.RequireStableFetchOffsets(),
		kgo.TransactionTimeout(10 * time.Second),
		kgo.RebalanceTimeout(30 * time.Second),
	}

	// if debug {
	// 	opts = append(opts, kgo.WithLogger(&stdoutLogger{}))
	// }

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %v", err)
	}

	return client, nil
}
