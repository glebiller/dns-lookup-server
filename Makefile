PACKAGES="./..."
BUILD_FOLDER = dist

all: build

build: clean dist/artifacts.json dist/manifests.yaml
.PHONY: build

dist/artifacts.json:
	goreleaser build --single-target --config .github/.goreleaser.yaml --snapshot --rm-dist --output dist/dns-lookup-server

dist/manifests.yaml:
	mkdir -p dist
	kustomize build deploy > dist/manifests.yaml

install: go.sum
	go install cmd/dns-lookup-server/main.go

clean:
	@echo clean build folder $(BUILD_FOLDER)
	rm -rf $(BUILD_FOLDER)
	@echo done
.PHONY: clean

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test: go-test kubernetes-test
.PHONY: test

go-test:
	@go test -mod=readonly $(PACKAGES) -cover -race
.PHONY: go-test

kubernetes-test: dist/manifests.yaml
	kube-score score dist/manifests.yaml
.PHONY: kubernetes-test

lint: go-lint kubernetes-lint
.PHONY: lint

go-lint:
	@golangci-lint run --config .github/.golangci.yaml
	@go mod verify
.PHONY: go-lint

kubernetes-lint:
	docker run --rm -v $(shell pwd):/data cytopia/yamllint -s .
.PHONY: lint
