FROM golang:1.8

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/golang/example/outyet

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/eric-kansas/cross-pollinators-server

WORKDIR /go/src/github.com/eric-kansas/cross-pollinators-server

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/cross-pollinators-server