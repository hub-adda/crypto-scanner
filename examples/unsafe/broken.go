package main

import (
	"crypto/des"
	"crypto/dsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/htruong/go-md2"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
)

const data = "Hello, World!"

func main() {

	// Unsafe hashes examples

	unsafe_md2_example()

	unsafe_md4_example()

	unsafe_md5_example()

	unsafe_sha1_example()

	unsafe_ripemd160_example()

	// unsafe symmetric crypto examples

	// unsafe_rc2_example() - need to complete usage of RC2

	unsafe_rc4_example()

	unsafe_blowfish_example()

	unsafe_cast5_example()

	unsafe_des_example()

	// Signature Algorithm examples

	unsafe_dsa_example()

	safe_rsa_example()

	// TLS examples

	safe_tls()

}

// func unsafe_rc2_example(){
// 	key := []byte("example key 1234")
// 	cipher, err := rc2.NewCipher(key)
// 	if err != nil {
// 		fmt.Println("Error creating cipher:", err)
// 		return
// 	}
// 	fmt.Printf("Encrypted: %v\n", cipher)
// }

func safe_rsa_example() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key:", err)
		return
	}
	fmt.Printf("Generated key: %v\n", priv)
}

func unsafe_des_example() {
	key := []byte("example key 1234")
	cipher, err := des.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return
	}
	fmt.Printf("Encrypted: %v\n", cipher)
}

func unsafe_cast5_example() {
	// key := []byte("example key 1234")
	// cipher, err := cast5.NewCipher(key)
	// if err != nil {
	// 	fmt.Println("Error creating cipher:", err)
	// 	return
	// }
	// fmt.Printf("Encrypted: %v\n", cipher)
}

func unsafe_blowfish_example() {
	key := []byte("example key 1234")
	cipher, err := blowfish.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return
	}
	fmt.Printf("Encrypted: %v\n", cipher)
}

func unsafe_ripemd160_example() {
	hash := ripemd160.New()
	hash.Write([]byte(data))
}

func unsafe_md5_example() {
	// hash := md5.New()
	// hash.Write([]byte(data))
}

func unsafe_md4_example() {
	hash := md4.New()
	hash.Write([]byte(data))
}

func unsafe_sha1_example() {
	// hash := sha1.New()
	// hash.Write([]byte(data))
}

func unsafe_md2_example() {
	hash := md2.New()
	hash.Write([]byte(data))
}

func unsafe_rc4_example() {
	// key := []byte("examplekey123")

	// Create the RC4 cipher
	// cipher, err := rc4.NewCipher(key)
	// if err != nil {
	// 	fmt.Println("Error creating cipher:", err)
	// 	return
	// }
	// fmt.Printf("Encrypted: %v\n", cipher)
}

func unsafe_dsa_example() {
	// Note that FIPS 186-5 no longer approves DSA for signature generation.
	// The DSA algorithm is only used for verifying signatures.
	// The ECDSA algorithm is the recommended replacement for DSA.

	hashed := []byte("testing")
	var priv dsa.PrivateKey
	params := &priv.Parameters

	err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160)
	if err != nil {
		fmt.Println("Error generating parameters:", err)
		return
	}
	_, _, err = dsa.Sign(rand.Reader, &priv, hashed)
	if err != nil {
		fmt.Println("Error signing:", err)
		return
	}
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
