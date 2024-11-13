# Building a FIPS-140 Compliant Binary

This guide provides instructions on how to build a Go binary in a FIPS-compliant environment using Docker. The process involves using a build script to create a Docker image and then running a compile script to generate the FIPS-140 compliant binary.

For more detailed instructions and background information, refer to this article: [Navigating FIPS Compliance for Go Applications: Libraries, Integration, and Security](https://medium.com/cyberark-engineering/navigating-fips-compliance-for-go-applications-libraries-integration-and-security-42ac87eec40b).

## Prerequisites

- Docker installed on your machine.
- Go source code that you want to compile.

## Steps

1. **Prepare the Build Environment**

   Ensure that your Go source code and the necessary scripts are in place. The following files should be present in the `examples/fips` directory:
   - `Dockerfile`
   - `scripts/compile.sh`
   - `scripts/build.sh`
   - `go.mod`
   - `fips_web_server.go`

2. **Build the Docker Image**

   Run the `build.sh` script to create the Docker image and prepare the build environment.

   ```bash
   cd /path/to/crypto-scanner/examples/fips
   ./scripts/build.sh
   ```

   This script will:
   - Stop and remove any existing Docker containers with the same name.
   - Remove any existing Docker images with the same name.
   - Create a temporary build directory and copy the necessary files into it.
   - Build the Docker image if it does not already exist.
   - Run the Docker container to compile the Go binary.

3. **Compile the FIPS-140 Compliant Binary**

   The `compile.sh` script inside the Docker container will compile the Go binary with FIPS-140 compliance.

   ```bash
   # This step is handled by the build.sh script, which runs the Docker container and executes the compile.sh script inside it.
   ```

   The `compile.sh` script will:
   - Clean the Go module cache.
   - Download the necessary Go modules.
   - Compile the Go binary with FIPS-140 compliance using the `boringcrypto` tag.

4. **Verify the Output**

   After the build process is complete, the FIPS-140 compliant binary will be available in the `build` directory.

   ```bash
   ls -la build
   ```

   You should see the `fips_web_server_linux` binary in the directory.

## Example Commands

To build the FIPS-140 compliant binary, use the following commands:
