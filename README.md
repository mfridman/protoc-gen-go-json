# protoc-gen-go-json

[![Build](https://github.com/mfridman/protoc-gen-go-json/actions/workflows/ci.yaml/badge.svg)][badges_ci]
[![Report
Card](https://goreportcard.com/badge/github.com/mfridman/protoc-gen-go-json)][badges_goreportcard]
[![Go Reference](https://pkg.go.dev/badge/github.com/mfridman/protoc-gen-go-json.svg)][badges_godoc]

This is a Protobuf plugin for Go that generates code to implement
[json.Marshaler](https://golang.org/pkg/encoding/json/#Marshaler) and
[json.Unmarshaler](https://golang.org/pkg/encoding/json/#Unmarshaler) using
[protojson](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson).

This enables Go-generated protobuf messages to be embedded directly within other structs and encoded
with the standard JSON library, since the standard `encoding/json` library can't encode certain
protobuf messages such as those that contain `oneof` fields.

> [!TIP]
>
> For list of chanegs between versions, see the [CHANGELOG.md](CHANGELOG.md).

## Install

```
go install github.com/mfridman/protoc-gen-go-json@latest
protoc-gen-go-json --version
```

Also required:

- [protoc](https://github.com/google/protobuf) or [buf](https://github.com/bufbuild/buf)
- [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf)

## Usage

Define your messages like normal:

```proto
syntax = "proto3";

package api.v1;

message Request {
  oneof kind {
    string name = 1;
    int32 code = 2;
  }
}
```

The example message purposely uses a `oneof` since this won't work by default with `encoding/json`.
Next, generate the code:

#### Using `protoc`

```
protoc --go_out=. --go-json_out=. request.proto
```

#### Using `buf`

```yaml
version: v2
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen/go
    opt: paths=source_relative
  - remote: buf.build/community/mfridman-go-json
    out: gen/go
    opt:
      - paths=source_relative
      - orig_name=true
```

And then run:

```sh
buf generate request.proto
```

Your output should contain a file `request.pb.json.go` which contains the implementation of
`json.Marshal/Unmarshal` for all your message types. You can then encode your messages using
standard `encoding/json`:

```go
// Marshal
by, err := json.Marshal(&apiv1.Request{
  Kind: &apiv1.Request_Name{
    Name: "alice",
	},
})
if err != nil {
  log.Fatal(err)
}
fmt.Println(string(by))
// {"name":"alice"}

// Unmarshal
var request apiv1.Request
if err := json.Unmarshal(by, &request); err != nil {
	log.Fatal(err)
}
fmt.Println(request.GetName())
// alice
```

From the root of this repository, you can run the following to see the example in action:

```sh
go run ./examples
```

## Options

The generator supports options to control the behavior of the generated code. The options are passed
as a comma-separated list to the `--go-json_out` flag.

#### Marshal Options

See the
[protojson.MarshalOptions](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson#MarshalOptions)
for documentation on these options.

| Option                       | Description                                        | Default |
| ---------------------------- | -------------------------------------------------- | ------- |
| `enums_as_ints`              | Render enums as integers instead of strings        | `false` |
| `emit_defaults`              | Render fields with zero values                     | `false` |
| `emit_defaults_without_null` | Render fields with zero values without null values | `false` |
| `orig_name`                  | Use original (.proto file) name for fields         | `false` |

#### Unmarshal Options

See the
[protojson.UnmarshalOptions](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson#UnmarshalOptions)
for documentation on these options.

| Option          | Description                               | Default |
| --------------- | ----------------------------------------- | ------- |
| `allow_unknown` | Disallow unknown fields when unmarshaling | `false` |

It also includes the "standard" options available to all
[protogen](https://pkg.go.dev/google.golang.org/protobuf/compiler/protogen?tab=doc)-based plugins:

- `import_path={path}` - Override the import path
- `paths=source_relative` - Derive the output path from the input path
- etc.

These can be set as part of the `--go-json_out` value:

```sh
protoc --go-json_opt=emit_defaults=true:.
```

You can specify multiple using a `,`:

```sh
protoc --go-json_out=enums_as_ints=true,emit_defaults=true:.
```

Alternatively, you may also specify options using the `--go-json_opt` value:

```sh
protoc --go-json_out:. --go-json_opt=emit_defaults=true,enums_as_ints=true
```

## Acknowledgements

This project is a clone of
[mitchellh/protoc-gen-go-json](https://github.com/mitchellh/protoc-gen-go-json). The original
project is no longer maintained and this project is a continuation of it. To learn more see
Mitchell's [Planned Repo
Archive](https://gist.github.com/mitchellh/90029601268e59a29e64e55bab1c5bdc) gist for more
information.

[badges_ci]: https://github.com/mfridman/protoc-gen-go-json/actions/workflows/ci.yaml
[badges_goreportcard]: https://goreportcard.com/report/github.com/mfridman/protoc-gen-go-json
[badges_godoc]: https://pkg.go.dev/github.com/mfridman/protoc-gen-go-json
