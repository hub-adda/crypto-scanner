# crypto-scanner

Golang is a powerful programming language, but like any other language, it is crucial to ensure that cryptographic implementations are safe and secure. Unsafe cryptographic practices can lead to vulnerabilities and potential security breaches.

## About

crypto-scanner is a comprehensive tool designed to enhance the security of Go applications by identifying unsafe cryptographic implementations. It consists of the following components:

- **Binary Checker**: Scans compiled binaries to detect the use of unsafe cryptographic functions. It supports two configurations:
  - **Unsafe Cryptography**: Identifies the use of cryptographic functions that are considered insecure.
  - **FIPS-Compliant**: Ensures that the cryptographic functions used comply with FIPS (Federal Information Processing Standards) specific rules.
- **Code Checker**: Analyzes source code to identify insecure cryptographic practices.
- **Future Plans**: Aims to develop a Go fork that removes unsafe implementations from the standard crypto libraries, ensuring that only secure cryptographic functions are available.

## Usage

To build the example code with different configurations, use the following commands:

### Building a Binary

To build a binary for general use:
``` bash
env GOOS=linux GOARCH=arm go build -o broken ./examples/broken.go
```
To build a FIPS-140 compliant binary:
### build the FIPS binary in a a linux machine 
```
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOEXPERIMENT=boringcrypto go build -tags boringcrypto -o fips_web_server_linux ./fips_web_server.go
```
In some machines a FIPS binary cannot build with as local cc compiler, therefore there is a way to build it in a docer file

### build the FIPS binary in using docker

