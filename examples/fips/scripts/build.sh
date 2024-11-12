#!/usr/bin/env bash
set -e
image_name=go_builder_image:latest
temp_dir=build
container_name=go_builder_container

docker stop $container_name || true
docker rm $container_name || true
docker rmi $image_name || true

# check if directory exists and delete the contemt of dir if exists
rm -rf $temp_dir
if [ -d $temp_dir ]; then
  rm -rf $temp_dir/*
else
    mkdir -p $temp_dir
    chmod -R 777 $temp_dir
fi

cp go.mod $temp_dir
cp fips_web_server.go $temp_dir

# if image doesnet exist, create it
if [[ "$(docker images -q $image_name 2> /dev/null)" == "" ]]; then
  docker build -t $image_name --no-cache --progress=plain . &> build.log
fi

docker run -v $(pwd)/$temp_dir:/mnt --name go_builder_container go_builder_image

