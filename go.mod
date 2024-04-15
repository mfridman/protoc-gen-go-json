module github.com/mfridman/protoc-gen-go-json

go 1.22.1

require (
	github.com/bufbuild/protoplugin v0.0.0-20240323223605-e2735f6c31ee
	github.com/mfridman/buildversion v0.3.0
	github.com/stretchr/testify v1.9.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v1.1.0 // Invalid module reference (old repo and tag)
	v1.0.0 // Invalid module reference (old repo and tag)
)
