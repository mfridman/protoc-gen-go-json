.PHONY: build
build:
	mkdir -p build
	go build -o build/protoc-gen-go-json .
	GOBIN=$(CURDIR)/build go install google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: proto
proto: proto-buf proto-protoc

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
	go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...

.PHONY: proto-buf
proto-buf: build
	export PATH=$(CURDIR)/build/:$$PATH && \
		cd examples && buf generate proto

.PHONY: proto-protoc
proto-protoc: build
	export PATH=$(CURDIR)/build/:$$PATH && \
		protoc --go_out=. -I./e2e --go-json_out=orig_name=true,emit_scanner_valuer=true:. e2e/*.proto

.PHONY: pinact
pinact:
	GITHUB_TOKEN=$$(gh auth token) pinact run --update --verify
