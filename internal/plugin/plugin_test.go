package plugin

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/bufbuild/protoplugin"
	"github.com/mfridman/protoc-gen-go-json/e2e"
	"github.com/mfridman/protoc-gen-go-json/internal/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestHandleMatchesCommittedE2E(t *testing.T) {
	t.Parallel()

	// Use the same parameters as the proto-protoc target in the Makefile.
	response := runPlugin(t, e2e.File_e2e_proto.Path(), "orig_name=true,emit_scanner_valuer=true")
	require.Empty(t, response.GetError())
	require.Len(t, response.GetFile(), 1)

	file := response.GetFile()[0]
	assert.Equal(t, "e2e.pb.json.go", filepath.Base(file.GetName()))
	want, err := os.ReadFile(filepath.Join("..", "..", "e2e", "e2e.pb.json.go"))
	require.NoError(t, err)
	assert.Equal(t, string(want), file.GetContent(), "e2e/e2e.pb.json.go is stale, run make proto")
}

func TestHandleSkipsFilesWithoutMessages(t *testing.T) {
	t.Parallel()

	response := runPlugin(t, e2e.File_e2e_service_proto.Path(), "")
	require.Empty(t, response.GetError())
	assert.Empty(t, response.GetFile())
}

func TestHandleInvalidParameter(t *testing.T) {
	t.Parallel()

	request := newRequest(e2e.File_e2e_proto.Path(), "orig_name")
	err := protoplugin.Run(
		context.Background(),
		protoplugin.Env{
			Stdin:  bytes.NewReader(mustMarshal(t, request)),
			Stdout: new(bytes.Buffer),
			Stderr: new(bytes.Buffer),
		},
		protoplugin.HandlerFunc(Handle),
	)
	require.ErrorContains(t, err, "must be in the form of name=value")
}

func TestParseOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		want    *gen.Options
		wantErr string
	}{
		{
			name: "empty",
			raw:  "",
			want: &gen.Options{},
		},
		{
			name: "all options",
			raw: "enums_as_ints=true,emit_defaults=true,emit_defaults_without_null=true," +
				"orig_name=true,allow_unknown=true,emit_scanner_valuer=true",
			want: &gen.Options{
				EnumsAsInts:        true,
				EmitDefaults:       true,
				EmitDefaultValues:  true,
				OrigName:           true,
				AllowUnknownFields: true,
				EmitScannerValuer:  true,
			},
		},
		{
			name: "unknown options are ignored",
			raw:  "paths=source_relative,orig_name=1",
			want: &gen.Options{OrigName: true},
		},
		{
			name:    "missing value",
			raw:     "orig_name",
			wantErr: "must be in the form of name=value",
		},
		{
			name:    "invalid bool",
			raw:     "orig_name=yes",
			wantErr: "invalid value for orig_name",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseOptions(test.raw)
			if test.wantErr != "" {
				require.ErrorContains(t, err, test.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func runPlugin(t *testing.T, fileToGenerate string, parameter string) *pluginpb.CodeGeneratorResponse {
	t.Helper()
	request := newRequest(fileToGenerate, parameter)
	stdout := new(bytes.Buffer)
	err := protoplugin.Run(
		context.Background(),
		protoplugin.Env{
			Stdin:  bytes.NewReader(mustMarshal(t, request)),
			Stdout: stdout,
			Stderr: new(bytes.Buffer),
		},
		protoplugin.HandlerFunc(Handle),
	)
	require.NoError(t, err)
	response := new(pluginpb.CodeGeneratorResponse)
	require.NoError(t, proto.Unmarshal(stdout.Bytes(), response))
	return response
}

func newRequest(fileToGenerate string, parameter string) *pluginpb.CodeGeneratorRequest {
	// The service file imports the message file, so the message file must come
	// first. This is the topological order that protoc and buf use.
	protoFiles := []*descriptorpb.FileDescriptorProto{
		protodesc.ToFileDescriptorProto(e2e.File_e2e_proto),
		protodesc.ToFileDescriptorProto(e2e.File_e2e_service_proto),
	}
	return &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{fileToGenerate},
		Parameter:      proto.String(parameter),
		ProtoFile:      protoFiles,
	}
}

func mustMarshal(t *testing.T, message proto.Message) []byte {
	t.Helper()
	data, err := proto.Marshal(message)
	require.NoError(t, err)
	return data
}
