package fips

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"

	_ "crypto/tls/fipsonly"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

const data = "Hello, World!"

func main() {

	// TLS examples

	safe_tls()

	unsafe_sha1_example()

}

func unsafe_sha1_example() {
	// SHA1 is no longer considered safe
	fmt.Println(sha1.BlockSize)
	h := sha1.New()
	h.Write([]byte(data))
	fmt.Printf("SHA1: %x\n", h.Sum(nil))
}

func safe_tls() {
	// keyFileName is the name of the file where the private key will be saved
	keyFileName := "private.pem"
	// certFileName is the name of the file where the self-signed certificate will be saved
	certFileName := "cert.pem"
	err := CreateSelfSignedKeyAndCertFiles(keyFileName, certFileName)
	if err != nil {
		log.Fatalf("Failed to create key and cert files: %v", err)
	}
	// Create a new TLS configuration
	cfg := &tls.Config{
		MinVersion:       tls.VersionTLS11,
		MaxVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{},
	}

	// Create a new HTTP server mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	// Set an HTTP server instance with configuration with the TLS configuration
	srv := &http.Server{
		Addr:      ":443",
		Handler:   mux,
		TLSConfig: cfg,
	}

	// Start the server
	fmt.Printf("Starting server on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServeTLS(certFileName, keyFileName))

}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is an example server.\n"))
}

// CreateSelfSignedKeyAndCertFiles generates a private key and a self-signed certificate
// and saves them to the specified files
func CreateSelfSignedKeyAndCertFiles(keyFileName, certFileName string) error {

	// Generate a private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Encode the private key to the PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}

	// Save the private key to a file
	privateKeyFile, err := os.Create(keyFileName)
	if err != nil {
		return fmt.Errorf("error creating private key file: %v", err)
	}
	defer privateKeyFile.Close()

	// Encode the private key to the PEM format
	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		return fmt.Errorf("error encoding private key to PEM: %v", err)
	}

	// Create a template for the certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"logo"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	// Create a self-signed certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	// convert certificate DER to PEM
	cert := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}

	// Save the certificate to a file
	certFile, err := os.Create(certFileName)
	if err != nil {
		return fmt.Errorf("error creating cert file:%v", err)
	}

	defer func() error {
		err = certFile.Close()
		if err != nil {
			return fmt.Errorf("error closing cert file:%v ", err)
		}
		return nil
	}()

	// Encode the certificate to the PEM format
	err = pem.Encode(certFile, cert)
	if err != nil {
		return fmt.Errorf("error encoding cert to PEM:%v", err)
	}

	return nil
}
