# crypto-scanner

Golang is a powerful programming language, but like any other language, it is crucial to ensure that cryptographic implementations are safe and secure. Unsafe cryptographic practices can lead to vulnerabilities and potential security breaches.

## About

crypto-scanner is a comprehensive tool designed to enhance the security of Go applications by identifying unsafe cryptographic implementations. It consists of the following components:

- **Binary Checker**: Scans compiled binaries to detect the use of unsafe cryptographic functions. It supports two configurations:
  - **Unsafe Cryptography**: Identifies the use of cryptographic functions that are considered insecure.
  - **FIPS-Compliant**: Ensures that the cryptographic functions used comply with FIPS (Federal Information Processing Standards) specific rules.
- **Code Checker** [Future]: Analyzes source code to identify insecure cryptographic practices.
- **Safe Compiler [future]**: A Go compiler that prevents compilation if unsafe implementations are found. It should identify unsafe usage of standard crypto libraries such as MD5, SHA-1. It will ensure that only secure cryptographic functions are available.

## Installation

1. un a single command on a temp directory to download and build the tool
``` bash
curl -s  https://raw.githubusercontent.com/hub-adda/crypto-scanner/refs/heads/main/install.sh | bash
```

2. Or clone the repo and build the tools:
```
git clone https://github.com/GilAddaCyberark/crypto-scanner/
cd crypto-scanner
chmod a+x ./install.sh
./install.sh
```

## Usage

### Checking a binary for safe cryptographic usage
 
Run this command
```
./binary-checker -binary my_binary -profile default.yaml 
``` 
The output should look like 
```
Using profile file: /Users/gil.adda/tmp/crypto-scanner/default.yaml
Scanning binary file: /Users/gil.adda/tmp/crypto-scanner/binary-checker
Check: 'The file is a valid Go binary' found.  
Check: 'MD4 Algorithm Usage' not found.  
Check: 'RIPEMD-160 Algorithm Usage' not found.  
Check: 'RC4 Algorithm Usage' not found.  
Check: 'Blowfish Algorithm Usage' not found.  
Check: 'CAST5 Algorithm Usage' not found.  
```
### What about SHA-1, md5 and RSA with small keys?
TBD...

### Checking a binary for FIPS-140 compliant usage 

Run this command
```
./binary-checker -binary my_binary -profile fips.yaml 
```

To build the example code with different configurations, use the following commands:

### Building a Binary

To build a binary for general use:
``` bash
env GOOS=linux GOARCH=arm go build -o broken ./examples/broken.go
```
To build a FIPS-140 compliant binary:
### build the FIPS binary in a linux machine 
```
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOEXPERIMENT=boringcrypto go build -tags boringcrypto -o fips_web_server_linux ./fips_web_server.go
```
In some machines a FIPS binary cannot build with a local cc compiler, therefore there is a way to build it in a docker file

### build the FIPS binary using docker

## TODOs in next steps
1. binary-checker: add go version conditions as GOOS or minimal go version (output below)
1. binary-checker: download and build the tool using a single command (e.g curl install script| bash
1. complete code usage of unsafe functions as sha1, md5 which are part of the go standard libraries
1. add contribution and license policy

### go version output example
This is the example output of the go version tool. 
``` bash
go version -m binary-checker
```

``` txt
binary-checker: go1.23.1
        path    crypto-scanner/cmd/binary-checker
        mod     crypto-scanner  (devel)
        dep     gopkg.in/yaml.v2        v2.4.0  h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=
        build   -buildmode=exe
        build   -compiler=gc
        build   CGO_ENABLED=1
        build   CGO_CFLAGS=
        build   CGO_CPPFLAGS=
        build   CGO_CXXFLAGS=
        build   CGO_LDFLAGS=
        build   GOARCH=arm64
        build   GOOS=darwin
        build   GOARM64=v8.0
        build   vcs=git
        build   vcs.revision=6bb4da79ff6a79852d58eece9dc88e1ea69c483d
        build   vcs.time=2024-11-12T13:42:30Z
        build   vcs.modified=false
```

``` bash
go version -m binary-checker examples/fips/build/fips_web_server_linux
```

This is the output of a FIPS-140 compliant binary  ( GOEXPERIMENT=boringcrypto, )


``` text
binary-checker: go1.23.1
        path    crypto-scanner/cmd/binary-checker
        mod     crypto-scanner  (devel)
        dep     gopkg.in/yaml.v2        v2.4.0  h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=
        build   -buildmode=exe
        build   -compiler=gc
        build   CGO_ENABLED=1
        build   CGO_CFLAGS=
        build   CGO_CPPFLAGS=
        build   CGO_CXXFLAGS=
        build   CGO_LDFLAGS=
        build   GOARCH=arm64
        build   GOOS=darwin
        build   GOARM64=v8.0
        build   vcs=git
        build   vcs.revision=6bb4da79ff6a79852d58eece9dc88e1ea69c483d
        build   vcs.time=2024-11-12T13:42:30Z
        build   vcs.modified=false
examples/fips/build/fips_web_server_linux: go1.23.2 X:boringcrypto
        path    command-line-arguments
        build   -buildmode=exe
        build   -compiler=gc
        build   -tags=boringcrypto
        build   CGO_ENABLED=1
        build   CGO_CFLAGS=
        build   CGO_CPPFLAGS=
        build   CGO_CXXFLAGS=
        build   CGO_LDFLAGS=
        build   GOARCH=arm64
        build   GOEXPERIMENT=boringcrypto
        build   GOOS=linux
        build   GOARM64=v8.0
```
The future plan is to set conditions and validate the binary under those conditions

