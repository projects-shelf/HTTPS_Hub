package cert

import (
	"os"
	"time"
)

type StoredCert struct {
	Cert       []byte
	PrivateKey []byte
	IssuerCert []byte
	NotAfter   time.Time
}

func SaveCertFiles(cert *CertResult, basePath string, notAfter time.Time) error {
	if err := os.WriteFile(basePath+".crt", cert.CertPEM, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(basePath+".key", cert.PrivateKeyPEM, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(basePath+".issuer.crt", cert.IssuerCertPEM, 0600); err != nil {
		return err
	}
	return os.WriteFile(basePath+".expiry", []byte(notAfter.Format(time.RFC3339)), 0600)
}

func LoadCertInfo(basePath string) (*StoredCert, error) {
	cert, err := os.ReadFile(basePath + ".crt")
	if err != nil {
		return nil, err
	}
	key, err := os.ReadFile(basePath + ".key")
	if err != nil {
		return nil, err
	}
	issuer, err := os.ReadFile(basePath + ".issuer.crt")
	if err != nil {
		return nil, err
	}
	expiryData, err := os.ReadFile(basePath + ".expiry")
	if err != nil {
		return nil, err
	}
	notAfter, err := time.Parse(time.RFC3339, string(expiryData))
	if err != nil {
		return nil, err
	}
	return &StoredCert{
		Cert:       cert,
		PrivateKey: key,
		IssuerCert: issuer,
		NotAfter:   notAfter,
	}, nil
}
