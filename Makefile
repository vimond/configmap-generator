VERSION = $(shell git describe --tags --always)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
NAME = $DOCKER_PRIVATE_REPO/k8s-configmap-generator
.PHONY: build test tag_latest release ssh

build:
	go get -t -d ./...
	go build  ./...

test:
	go test

release:
	go get github.com/mitchellh/gox

	CGO_ENABLED=0 gox -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(DATE)" -output "dist/configmapgen_{{.OS}}_{{.Arch}}" -arch "amd64" -os "linux windows darwin" ./...

docker:
	docker login -u $(ARTIFACTORY_USER) -p $(ARTIFACTORY_PASSWORD) -e developers@vimond.com $(DOCKER_PRIVATE_REPO)
	docker build  -t $(NAME):$(VERSION)  .
	docker tag  $(NAME):$(VERSION) $(NAME):latest
	@if ! docker images $(NAME) | awk '{ print $$2 }' | grep -qs -F $(VERSION); then echo "$(NAME) version $(VERSION) is not yet built. Please run 'make build'"; false; fi
	docker push $(NAME)
