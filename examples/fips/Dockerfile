# Use the official Golang image as the base image
FROM golang:1.23.2
RUN mkdir -p /mnt
RUN chmod 777 /mnt

RUN mkdir -p /src
RUN chmod 777 /src


WORKDIR /src
COPY scripts/compile.sh /src/compile.sh

RUN pwd
RUN ls -la 

CMD /src/compile.sh