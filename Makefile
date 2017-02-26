VERSION = $(shell git describe --tags --always)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")


.PHONY: build test tag_latest release ssh

build:
	go get -t -d ./...
	go build  ./...

test:
	go test

release:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr
	CGO_ENABLED=0 gox -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(DATE)" -output "dist/configmapgen_{{.OS}}_{{.Arch}}" -arch "amd64" -os "linux windows darwin" ./...

