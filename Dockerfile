FROM golang:1.6
MAINTAINER "CoreOS, Inc"
EXPOSE 8087

COPY . /go/src/github.com/storageos/discovery.etcd.io
RUN go install -v github.com/storageos/discovery.etcd.io

CMD ["discovery.etcd.io"]
