#!/usr/bin/env bash

# Check if git is installed
if ! command -v git &> /dev/null
then
    echo "git could not be found. Please install git to proceed."
    exit 1
fi

REPO_NAME=https://github.com/GilAddaCyberark/crypto-scanner/
BRANCH=draft
# Git clone it
git clone --branch $BRANCH $REPO_NAME


# Check if golang compiler  is installed
if ! command -v go &> /dev/null
then
    echo "go language compiler could not be found. Please install go to proceed."
    exit 1
fi

# Change to the repository directory
cd crypto-scanner

# Build the binary-checker tool
cd cmd/binary-checker
go build -o binarychecker

echo "binarychecker tool has been built successfully."