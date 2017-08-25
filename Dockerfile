FROM golang:1.8

ADD . /go/src/github.com/eric-kansas/cross-pollinators-server

RUN go install github.com/eric-kansas/cross-pollinators-server

ENTRYPOINT /go/bin/cross-pollinators-server