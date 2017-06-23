FROM golang:1.8.3
COPY . /go/src/github.com/storageos/discovery
WORKDIR /go/src/github.com/storageos/discovery
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo  -ldflags  -'w' -o discovery .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/storageos/discovery/discovery /bin/discovery
ENTRYPOINT ["/bin/discovery"]

EXPOSE 8081