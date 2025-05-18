package main

import (
	"back/cert"
	"back/config"
	"back/dns"
	"back/proxy"
	"crypto/tls"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	certMu     sync.RWMutex
	tlsCert    *tls.Certificate
	reloadFreq = 6 * time.Hour
)

func main() {
	// Manage A record
	dns.UpsertARecord()

	// Load subdomain:port table
	portMap, err := config.LoadPortMap()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/*path", proxy.ProxyHandler(portMap))

	// Auto generate SSL cert
	_, err = cert.AutoRenew()
	if err != nil {
		log.Fatalf("Failed to obtain certificate: %v", err)
	}

	if err := loadCertificate(); err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}

	go func() {
		ticker := time.NewTicker(reloadFreq)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Checking certificate renewal...")
			if _, err := cert.AutoRenew(); err != nil {
				log.Printf("Certificate renewal failed: %v", err)
				continue
			}
			if err := loadCertificate(); err != nil {
				log.Printf("Certificate reload failed: %v", err)
			} else {
				log.Println("Certificate reloaded successfully")
			}
		}
	}()

	httpsSrv := &http.Server{
		Addr:    ":443",
		Handler: r,
		TLSConfig: &tls.Config{
			GetCertificate: getCertificateFunc(),
		},
	}

	httpSrv := &http.Server{
		Addr: ":80",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			target := "https://" + req.Host + req.URL.RequestURI()
			http.Redirect(w, req, target, http.StatusMovedPermanently)
		}),
	}

	go func() {
		log.Println("Starting HTTP server (redirecting to HTTPS) on :80...")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Println("Starting HTTPS server (Gin + hot reload) on :443...")
	err = httpsSrv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("HTTPS server error: %v", err)
	}
}

func loadCertificate() error {
	certMu.Lock()
	defer certMu.Unlock()

	cert, err := tls.LoadX509KeyPair("/certs/certs.crt", "/certs/certs.key")
	if err != nil {
		return err
	}
	tlsCert = &cert
	return nil
}

func getCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		certMu.RLock()
		defer certMu.RUnlock()
		return tlsCert, nil
	}
}
