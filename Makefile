# Makefile template borrowed from https://sohlich.github.io/post/go_makefile/
# lowercase helper
lc = $(subst A,a,$(subst B,b,$(subst C,c,$(subst D,d,$(subst E,e,$(subst F,f,$(subst G,g,$(subst H,h,$(subst I,i,$(subst J,j,$(subst K,k,$(subst L,l,$(subst M,m,$(subst N,n,$(subst O,o,$(subst P,p,$(subst Q,q,$(subst R,r,$(subst S,s,$(subst T,t,$(subst U,u,$(subst V,v,$(subst W,w,$(subst X,x,$(subst Y,y,$(subst Z,z,$1))))))))))))))))))))))))))
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
LIST_SERVICES=listServices

ifeq ($(OS),Windows_NT)
	@echo "Windows not supported by this Makefile, sorry!"
else
	UNAME_S := $(shell uname -s)
	UNAME_LOWER = $(call lc,$(UNAME_S))

	UNAME_P := $(shell uname -p)
	ifneq ($(filter arm%,$(UNAME_P)),)
		# TODO: support arm/raspberrypi
		LOWER_UNAME = "raspi"
	endif
endif

build: build-$(UNAME_LOWER)
build-all: test build-linux build-arm build-darwin build-raspi
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f "bin/$(BINARY_NAME)_*"

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o "./bin/$(LIST_SERVICES)_linux" -ldflags $(GOFLAGS) -v "cmd/$(LIST_SERVICES)/main.go"
build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o "./bin/$(LIST_SERVICES)_arm" -ldflags $(GOFLAGS) -v "cmd/$(LIST_SERVICES)/main.go"
build-raspi:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o "./bin/$(LIST_SERVICES)_raspi" -ldflags $(GOFLAGS) -v "cmd/$(LIST_SERVICES)/main.go"
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o "./bin/$(LIST_SERVICES)_darwin" -ldflags $(GOFLAGS) -v "cmd/$(LIST_SERVICES)/main.go"
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w "/go/src/github.com/novu/$(PROJECT)" golang:$(GOVERSION) $(GOBUILD) -o "./bin/$(LIST_SERVICES)" -ldflags $(GOFLAGS) -v "cmd/$(LIST_SERVICES)/main.go"
