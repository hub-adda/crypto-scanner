#!/usr/bin/env bash
set -x
set -v
# Check if git is installed
if ! command -v git &> /dev/null
then
    echo "git could not be found. Please install git to proceed."
    exit 1
fi

REPO_NAME=https://github.com/GilAddaCyberark/crypto-scanner/
BRANCH=draft
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
echo - script dir $SCRIPT_DIR
TMP_DIR_NAME=tmp
echo script dir $SCRIPT_DIR
TMP_DIR_PATH=$SCRIPT_DIR/$TMP_DIR_NAME
echo temp dir: $TMP_DIR_PATH

mkdir -p $TMP_DIR_NAME
cd $TMP_DIR_PATH
echo temp - $(pwd)

# Git clone it
git clone --branch $BRANCH $REPO_NAME

# Check if golang compiler  is installed
if ! command -v go &> /dev/null
then
    echo "go language compiler could not be found. Please install go to proceed."
    exit 1
fi

cd $TMP_DIR_PATH/crypto-scanner
echo $(pwd)
go mod download

cd $TMP_DIR_PATH/crypto-scanner/cmd/crypto-checker
echo build dir: $(pwd)

go build -o ./binary-checker $TMP_DIR_PATH/crypto-scanner/cmd/binary-checker

ls -l $TMP_DIR_PATH/crypto-scanner/cmd/binary-checker

echo "binary-checker tool has been built successfully"

cp $TMP_DIR_PATH/crypto-scanner/binary-checker ../../

rm -rf $TMP_DIR_PATH
