package cert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
)

type User struct {
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string { return "" }
func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

type CertResult struct {
	CertPEM       []byte
	PrivateKeyPEM []byte
	IssuerCertPEM []byte
	NotAfter      time.Time
}

func obtain(domains []string) (*CertResult, error) {
	privateKey, err := certcrypto.GeneratePrivateKey(certcrypto.RSA2048)
	if err != nil {
		return nil, err
	}
	user := &User{key: privateKey}
	config := lego.NewConfig(user)
	config.CADirURL = lego.LEDirectoryProduction
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	provider, err := cloudflare.NewDNSProvider()
	if err != nil {
		return nil, err
	}
	if err := client.Challenge.SetDNS01Provider(provider, dns01.AddRecursiveNameservers([]string{"1.1.1.1:53", "8.8.8.8:53"})); err != nil {
		return nil, err
	}
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.Registration = reg

	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}
	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}
	notAfter, err := extractNotAfter(cert.Certificate)
	if err != nil {
		return nil, err
	}
	return &CertResult{
		CertPEM:       cert.Certificate,
		PrivateKeyPEM: cert.PrivateKey,
		IssuerCertPEM: cert.IssuerCertificate,
		NotAfter:      notAfter,
	}, nil
}

func extractNotAfter(certPEM []byte) (time.Time, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return time.Time{}, fmt.Errorf("invalid PEM data")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, err
	}
	return cert.NotAfter, nil
}

func AutoRenew() (*CertResult, error) {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("fail to read DOMAIN env")
	}
	domains := []string{domain, "*." + domain}
	basePath := "/certs/certs"
	renewBefore := 30 * 24 * time.Hour

	stored, err := LoadCertInfo(basePath)
	needRenew := err != nil || time.Until(stored.NotAfter) < renewBefore
	if !needRenew {
		return &CertResult{
			CertPEM:       stored.Cert,
			PrivateKeyPEM: stored.PrivateKey,
			IssuerCertPEM: stored.IssuerCert,
			NotAfter:      stored.NotAfter,
		}, nil
	}

	cert, err := obtain(domains)
	if err != nil {
		return nil, err
	}
	if err := SaveCertFiles(cert, basePath, cert.NotAfter); err != nil {
		return nil, err
	}
	return cert, nil
}
