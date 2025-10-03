package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tjalfe/jndispatcher/internal/arguments"
	"github.com/tjalfe/jndispatcher/internal/config"
	"github.com/tjalfe/jndispatcher/internal/types"
	"github.com/tjalfe/jndispatcher/internal/verification"
	"github.com/tjalfe/pcrypt"
	"github.com/tjalfe/psign"
	"github.com/twmb/franz-go/pkg/kgo"
)

// PromptForPassword prompts the user for a password and encrypts it using pcrypt.
// The encrypted password is used in config yaml file so that the real password is not stored in plaintext.
func PromptForPassword() error {
	crypt, err := pcrypt.Init(types.PcryptInit)
	defer crypt.Zero()
	if err != nil {
		return err
	}
	err = crypt.Prompt()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	args, err := arguments.ParseArguments()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	if args.Encrypt {
		err := PromptForPassword()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encrypting password: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Normal execution of the application starts here
	conf, err := config.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	kafkaClient, err := config.InitKafkaClient(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Kafka client: %v\n", err)
		os.Exit(1)
	}
	defer kafkaClient.Close()

	// Load pool of trusted CA certificates
	trustedCaPool, err := verification.LoadTrustedCaPool(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading pool of CA for verification of signed messages: %v\n", err)
	}

	// Main loop to consume messages
	for {
		fetches := kafkaClient.PollFetches(context.Background())
		fetches.EachRecord(func(record *kgo.Record) {
			log.Printf("Received:")
			// Load the record.Value into JSON structure
			var message types.Message
			err := json.Unmarshal(record.Value, &message)
			if err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				//	return
			}
			signingCert, err := x509.ParseCertificate(message.Certificate)
			if err != nil {
				log.Printf("Error parsing signing certificate: %v", err)
				//	return
			}
			// Verify the signing certificate
			CertOK := verification.VerifyTrustSigningCertificate(signingCert, trustedCaPool)
			if CertOK != nil {
				log.Printf("Error verifying signing certificate: %v", CertOK)
			} else {
				log.Printf("Signing certificate is valid and trusted")
			}
			// TO BE MADE: HANDELING OF MESSAGE
			log.Printf("Message UUID: %v\n", message.MessageUUID)
			log.Printf("Signing Certificate: %v\n", string(message.CertificateCommonName))

			// Verify signature
			signerCertificate, _ := x509.ParseCertificate(message.Certificate)
			signer, _ := psign.NewPverifier(signerCertificate.PublicKey)
			signatureOK := signer.Verify(message.Payload, message.Signature)

			if signatureOK != nil {
				log.Printf("Signature NOT ok: %v\n", signatureOK)
			} else {
				log.Printf("Signature ok\n")
			}

			log.Printf("Payload: %v\n\n", string(message.Payload))
		})
		if err := fetches.Err(); err != nil {
			log.Printf("Fetch error: %v", err)
		}
	}

}
