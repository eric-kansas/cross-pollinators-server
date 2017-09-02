FROM golang:1.8

# Install OS-level language locales
ENV DEBIAN_FRONTEND=noninteractive LANG=en_US.UTF-8 LANGUAGE=en_US:en LC_ALL=en_US.UTF-8
RUN apt-get update -q &&\
    apt-get install -y --no-install-recommends locales locales-all &&\
    locale-gen "$LANG" &&\
    update-locale LANG="$LANG" &&\
    rm -rf /var/lib/apt/lists/* /tmp/*

# GRPC stuff based on Dockerfile for gRPC Go
## https://github.com/grpc/grpc-docker-library/blob/master/1.0/golang/Dockerfile
# GRPC: install protobuf
ENV PROTOBUF_VERSION=3.3.0
RUN apt-get update -q &&\
    apt-get -y --no-install-recommends install unzip &&\
    mkdir -p /tmp/protoc &&\
    curl -sSL https://github.com/google/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip -o /tmp/protoc/protoc.zip &&\
    cd /tmp/protoc &&\
    unzip protoc.zip &&\
    cp /tmp/protoc/bin/protoc /usr/local/bin &&\
    chmod go+rx /usr/local/bin/protoc &&\
    cd /tmp &&\
    rm -rf /var/lib/apt/lists/* /tmp/*

# GRPC: Get the source from GitHub and install protoc-gen-go
RUN go get google.golang.org/grpc &&\
    go get -u github.com/golang/protobuf/protoc-gen-go

# Install vendoring tool and CompileDeamon
ENV GO15VENDOREXPERIMENT=1
RUN go get -u github.com/kardianos/govendor &&\
    go get -u github.com/githubnemo/CompileDaemon

WORKDIR /go/src/github.com/tapjoy/beeswax-log-ingestion
ENTRYPOINT ["./deploy/docker-entrypoint.sh"]
