# Use the official Golang image as the base image
FROM golang:1.23.2

RUN mkdir -p /mnt /src && chmod 777 /mnt /src

WORKDIR /src
COPY scripts/compile.sh /src/compile.sh

RUN pwd && ls -la 

CMD /src/compile.sh