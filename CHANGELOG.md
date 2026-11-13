# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Template errors no longer fail silently. Before, the plugin skipped the file and reported
  success.
- An interrupt signal no longer makes the plugin exit with code 0 while generation is in progress.

### Changed

- Go 1.27 or later is now required to build and install the plugin. The generated code does not
  change and still works with older Go versions.
- Update `google.golang.org/protobuf` to v1.36.12 and `github.com/bufbuild/protoplugin` to the
  latest version.

## [v1.6.0] - 2026-03-05

### Added

- `emit_scanner_valuer` option to generate `database/sql.Scanner` and
  `database/sql/driver.Valuer` implementations for protobuf messages. Useful for storing JSON
  representations of protobuf messages in database columns.

## [v1.5.0] - 2025-01-09

- Correctly handle Protobuf Editions.

## [v1.4.1] - 2024-10-16

- General housekeeping and dependency updates.

## [v1.4.0] - 2024-04-14

- Add a changelog to the project, based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
- Add new `emit_defaults_without_null` option (#5)
  ([EmitDefaultValues](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson#MarshalOptions))
- Only generate fields that are not set to their zero value. This will cause the generated code to
  be smaller, but should not affect the behavior of the generated code (since all fields were
  booleans and the zero value is `false`).

```diff
// MarshalJSON implements json.Marshaler
func (msg *Basic) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
-		UseEnumNumbers:  false,
-		EmitUnpopulated: false,
+		UseProtoNames:   true,
	}.Marshal(msg)
}
```

[Unreleased]: https://github.com/mfridman/protoc-gen-go-json/compare/v1.6.0...HEAD
[v1.6.0]: https://github.com/mfridman/protoc-gen-go-json/compare/v1.5.0...v1.6.0
[v1.5.0]: https://github.com/mfridman/protoc-gen-go-json/releases/tag/v1.5.0
[v1.4.1]: https://github.com/mfridman/protoc-gen-go-json/releases/tag/v1.4.1
[v1.4.0]: https://github.com/mfridman/protoc-gen-go-json/releases/tag/v1.4.0
