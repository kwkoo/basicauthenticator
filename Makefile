PREFIX=github.com/kwkoo
PACKAGE=basicauthenticator

GOPATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GOBIN=$(GOPATH)/bin
IMAGENAME="kwkoo/$(PACKAGE)"
VERSION="0.1"
DB=userdb.tsv

.PHONY: run build clean test coverage image runcontainer
run:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run \
	  $(GOPATH)/src/$(PREFIX)/$(PACKAGE)/cmd/$(PACKAGE)/main.go \
	  -userdb config/$(DB)

build:
	@echo "Building..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/$(PACKAGE) $(PREFIX)/$(PACKAGE)/cmd/$(PACKAGE)

clean:
	rm -f \
	  $(GOPATH)/bin/$(PACKAGE) \
	  $(GOPATH)/pkg/*/$(PACKAGE).a \
	  $(GOPATH)/$(COVERAGEOUTPUT) \
	  $(GOPATH)/$(COVERAGEHTML)

image: 
	docker build --rm -t $(IMAGENAME):$(VERSION) $(GOPATH)

runcontainer:
	docker run \
	  --rm \
	  -it \
	  --name $(PACKAGE) \
	  -p 8080:8080 \
	  -e USERDB=/config/$(DB) \
	  -v $(GOPATH)/config:/config \
	  $(IMAGENAME):$(VERSION)
