module github.com/mfridman/protoc-gen-go-json

go 1.23.2

require (
	github.com/bufbuild/protoplugin v0.0.0-20250106231243-3a819552c9d9
	github.com/mfridman/buildversion v0.3.0
	github.com/stretchr/testify v1.10.0
	google.golang.org/protobuf v1.36.2
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
