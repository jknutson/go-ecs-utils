# Makefile template borrowed from https://sohlich.github.io/post/go_makefile/
PROJECT=ecs-utils
PROJECT_VERSION=`cat VERSION.txt`
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVERSION=1.15
GOFLAGS="-X main.BuildVersion=$(PROJECT_VERSION)"
BINARY_NAME=$(PROJECT)

build-all: test build-linux build-arm build-darwin build-raspi
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f "bin/$(BINARY_NAME)_*"

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o "./bin/$(BINARY_NAME)_linux" -ldflags $(GOFLAGS) -v "main.go"
build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o "./bin/$(BINARY_NAME)_arm" -ldflags $(GOFLAGS) -v "main.go"
build-raspi:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o "./bin/$(BINARY_NAME)_raspi" -ldflags $(GOFLAGS) -v "main.go"
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o "./bin/$(BINARY_NAME)_darwin" -ldflags $(GOFLAGS) -v "main.go"
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w "/go/src/github.com/novu/$(PROJECT)" golang:$(GOVERSION) $(GOBUILD) -o "./bin/$(BINARY_NAME)" -ldflags $(GOFLAGS) -v "main.go"
