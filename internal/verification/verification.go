package verification

import (
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/tjalfe/jndispatcher/internal/types"
)

func LoadTrustedCaPool(config types.Config) (*x509.CertPool, error) {
	// Create empty CA pool
	certPool := x509.NewCertPool()
	if config.TrustServerCa {
		// Add system CAs to the pool if configured to do so
		var err error
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("error loading system CA certificates: %w", err)
		}
	}
	// Add extra CAs from config
	for _, caPath := range config.TrustExtraCa {
		caCert, err := os.ReadFile(caPath)
		if err != nil {
			return nil, fmt.Errorf("error reading extra CA certificate %s: %w", caPath, err)
		}
		if ok := certPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("error appending extra CA certificate %s to pool", caPath)
		}
	}
	return certPool, nil
}

func VerifyTrustSigningCertificate(signingCertificate *x509.Certificate, signingCertChain *x509.CertPool, trustedCaPool *x509.CertPool) error {
	// Start checking the notBefore and notAfter validity period
	now := time.Now()
	if now.Before(signingCertificate.NotBefore) {
		return fmt.Errorf("signing certificate is not valid yet: notBefore=%s", signingCertificate.NotBefore)
	}
	if now.After(signingCertificate.NotAfter) {
		return fmt.Errorf("signing certificate has expired: notAfter=%s", signingCertificate.NotAfter)
	}
	// Verify the signing certificate against the trusted CA pool
	opts := x509.VerifyOptions{
		Roots:         trustedCaPool,
		Intermediates: signingCertChain,
	}
	if _, err := signingCertificate.Verify(opts); err != nil {
		return fmt.Errorf("signing certificate is not trusted: %w", err)
	}
	// All checks passed
	return nil
}
