.PHONY: proto
proto:
	mkdir -p build
	go build -o build/protoc-gen-go-json .
	export PATH=$(CURDIR)/build/:$$PATH && \
	    protoc --go_out=. -I./e2e --go-json_out=orig_name=true:. e2e/*.proto

.PHONY: test
test:
	go test -count=1 -v -race -cover ./...
