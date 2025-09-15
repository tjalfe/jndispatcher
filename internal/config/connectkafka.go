package config

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/tjalfe/jndispatcher/internal/types"

	"github.com/tjalfe/pcrypt"
	"github.com/twmb/franz-go/pkg/kgo"
	"gopkg.in/yaml.v3"
	"software.sslmate.com/src/go-pkcs12"
)

// ReadKafkaConfigs reads a YAML file and unmarshals it into an array of KafkaClusterConfig structs.
func readKafkaConfigs() ([]KafkaClusterConfig, error) {
	yamlFile, err := os.ReadFile(KafkaServerConf)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err) //Using %w for error wrapping
	}

	var configs []KafkaClusterConfig

	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err) //Using %w for error wrapping
	}

	return configs, nil
}

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

func initKafkaClient(config KafkaClusterConfig) (*kgo.Client, error) {
	// privateKey, certificate, certificateChain, err := readPkcs12(config.KafkaAuthCertificateStore, config.KafkaAuthCertificateStorePassword)
	fmt.Printf("Cluster name: %s\n", config.Name)
	_, certificate, certificateChain, err := readPkcs12(config.KafkaAuthCertificateStore, config.KafkaAuthCertificateStorePassword)
	if err != nil {
		return nil, fmt.Errorf("error reading PKCS#12 file: %w", err)
	}
	fmt.Printf("Certificate subject: %s\n", certificate.Subject)
	for _, cert := range certificateChain {
		fmt.Printf("\tCertificate in chain subject: %s\n", cert.Subject)
	}
	fmt.Println()
	return nil, nil
}

func getKafkaClients() (map[string]*kgo.Client, error) {
	// Read Kafka configurations from the YAML file
	kafkaConfigs, err := readKafkaConfigs()
	if err != nil {
		return nil, fmt.Errorf("error reading Kafka configs: %w", err)
	}

	//kafkaClient := make(map[string]*kgo.Client)
	for _, config := range kafkaConfigs {
		_, err = initKafkaClient(config)
		if err != nil {
			return nil, fmt.Errorf("error initializing Kafka client for cluster %s: %w", config.Name, err)
		}

	}

	// Print names of clusters
	for _, cluster := range kafkaConfigs {
		fmt.Println(cluster.Name)
	}

	return nil, nil
}
