FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y curl && \
    apt-get install tini && \
    curl -O https://dl.google.com/go/go1.20.1.linux-amd64.tar.gz && \
    tar -xvf go1.20.1.linux-amd64.tar.gz && \
    mv go /usr/local

ENV PATH="/usr/local/go/bin:${PATH}"

RUN bash

ENV MNT_DIR /mnt/nfs/filestore

WORKDIR ukeleleweb

COPY . .

RUN sed "s/localhost/0.0.0.0/g" cmd/ukuleleweb/main.go -i

RUN chmod +x /ukeleleweb/run.sh

RUN cd cmd/ukuleleweb/ && \
    go mod tidy && \
    GOOS=linux GOARCH=amd64 go build && \
    mkdir ukuleleweb-data
    
EXPOSE 8080

ENTRYPOINT ["/usr/bin/tini", "--"]

CMD ["/ukeleleweb/run.sh"]
