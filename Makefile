VERSION ?= $(shell git describe --always --tags)
BIN = prometheus-eventhub-adapter
BUILD_CMD = go build -o build/$(BIN)-$(VERSION)-$${GOOS}-$${GOARCH} &
IMAGE_REPO = azurets.azurecr.io
FMT_CMD = $(gofmt -s -l -w $(find . -type f -name '*.go' -not -path './vendor/*') | tee /dev/stderr)

default:
	$(MAKE) release
	$(MAKE) image-push

test: bootstrap
	test -z '$(FMT_CMD)'
	go vet $(go list ./... | grep -v /vendor/)
	golint -set_exit_status $(shell go list ./... | grep -v vendor)
	gas ./...
	ginkgo -r -v
bootstrap:
	glide install
build:
	go build -o $(BIN)
clean:
	rm -rf build vendor
	rm -f release image bootstrap $(BIN)
release: bootstrap
	@echo "Running build command..."
	bash -c '\
		export GOOS=linux; export GOARCH=amd64; export CGO_ENABLED=0; $(BUILD_CMD) \
		wait \
	'
image:
	@echo "Building the Docker image..."
	docker build -t $(IMAGE_REPO)/$(BIN):$(VERSION) .
	docker tag $(IMAGE_REPO)/$(BIN):$(VERSION) $(IMAGE_REPO)/$(BIN):latest

image-push: image
	docker push $(IMAGE_REPO)/$(BIN):$(VERSION)
	docker push $(IMAGE_REPO)/$(BIN):latest

.PHONY: test build clean image-push

