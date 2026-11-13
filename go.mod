module github.com/mfridman/protoc-gen-go-json

go 1.27

require (
	github.com/bufbuild/protoplugin v0.0.0-20260414125817-25d1d281b46b
	github.com/stretchr/testify v1.12.1
	google.golang.org/protobuf v1.36.12
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

retract (
	v1.1.0 // Invalid module reference (old repo and tag)
	v1.0.0 // Invalid module reference (old repo and tag)
)
