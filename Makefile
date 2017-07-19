JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -X github.com/storageos/discovery/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/storageos/discovery/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/storageos/discovery/version.BuildDate=$(JOBDATE)

test:
	go test -v `go list ./... | egrep -v /vendor/`

release:
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo  -ldflags "$(LDFLAGS)" -o discovery .	
