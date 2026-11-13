.PHONY: build
build:
	go build -o build/protoc-gen-go-json .

.PHONY: proto
proto:
	cd e2e && buf generate
	cd examples && buf generate proto

.PHONY: test
test:
	go test -count=1 -v -race -cover ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	golangci-lint fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vulncheck
vulncheck:
	go tool govulncheck ./...

.PHONY: pinact
pinact:
	GITHUB_TOKEN=$$(gh auth token) pinact run --update --verify
