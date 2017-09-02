FROM golang:1.9

WORKDIR /go/src/github.com/eric-kansas/cross-pollinators-server
COPY . .

RUN go-wrapper download 
RUN go-wrapper install  

#CMD ["go-wrapper", "run"]