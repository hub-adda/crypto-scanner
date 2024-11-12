#!/usr/bin/env bash

# Enable debugging and verbose output
set -x
set -v

# Check if git is installed
if ! command -v git &> /dev/null
then
    echo "git could not be found. Please install git to proceed."
    exit 1
fi

# Define variables
REPO_NAME=https://github.com/GilAddaCyberark/crypto-scanner/
REPO_DIR=crypto-scanner
BRANCH=draft
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TMP_DIR_NAME=tmp
TMP_DIR_PATH=$SCRIPT_DIR/$TMP_DIR_NAME

# Create a temporary directory and navigate into it
mkdir -p $TMP_DIR_NAME
cd $TMP_DIR_PATH

# Clone the specified branch of the repository
git clone --branch $BRANCH $REPO_NAME

# Check if Go compiler is installed
if ! command -v go &> /dev/null
then
    echo "Go language compiler could not be found. Please install Go to proceed."
    exit 1
fi

# Navigate to the cloned repository and download dependencies
cd $TMP_DIR_PATH/$REPO_DIR
go mod download

# Navigate to the build directory and build the binary
cd $TMP_DIR_PATH/$REPO_DIR/cmd/crypto-checker
go build -o ./binary-checker $TMP_DIR_PATH/$REPO_DIR/cmd/binary-checker

# List the built binary to confirm successful build
ls -l $TMP_DIR_PATH/$REPO_DIR/cmd/binary-checker
echo "binary-checker tool has been built successfully"

# Copy the built binary and configuration files to the parent directory
cp $TMP_DIR_PATH/$REPO_DIR/binary-checker ../../
cp $TMP_DIR_PATH/$REPO_DIR/profiles/default.yaml ../../
cp $TMP_DIR_PATH/$REPO_DIR/profiles/fips.yaml ../../

# Clean up by removing the temporary directory
rm -rf $TMP_DIR_PATH
