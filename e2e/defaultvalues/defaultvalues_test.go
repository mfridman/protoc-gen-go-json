package defaultvalues

import (
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	_ jsonv2.MarshalerTo     = (*Item)(nil)
	_ jsonv2.UnmarshalerFrom = (*Item)(nil)
)

// TestJSONV2Options checks the options in buf.gen.yaml at run time: the
// generated methods must give the same result as protojson with the same
// options.
func TestJSONV2Options(t *testing.T) {
	t.Parallel()

	marshalOptions := protojson.MarshalOptions{EmitDefaultValues: true}
	unmarshalOptions := protojson.UnmarshalOptions{}
	const unknownFieldsAllowed = false

	apis := map[string]struct {
		marshal   func(any) ([]byte, error)
		unmarshal func([]byte, any) error
	}{
		"v1": {jsonv1.Marshal, jsonv1.Unmarshal},
		"v2": {
			func(v any) ([]byte, error) { return jsonv2.Marshal(v) },
			func(b []byte, v any) error { return jsonv2.Unmarshal(b, v) },
		},
	}
	messages := []*Item{
		{},
		{DisplayName: "name", Kind: Item_KIND_A, Tags: []string{"a", "b"}},
		{Count: 3, Note: new("")},
	}
	for name, api := range apis {
		for _, msg := range messages {
			want, err := marshalOptions.Marshal(msg)
			require.NoError(t, err)
			got, err := api.marshal(msg)
			require.NoError(t, err, name)
			require.Equal(t, compact(t, want), compact(t, got), "%s: marshal %v", name, msg)

			out := new(Item)
			require.NoError(t, api.unmarshal(got, out), name)
			require.True(t, proto.Equal(msg, out), "%s: got %v, want %v", name, out, msg)
		}

		input := []byte(`{"display_name":"x","unknown":1,"kind":1}`)
		wantMsg := new(Item)
		wantErr := unmarshalOptions.Unmarshal(input, wantMsg)
		gotMsg := new(Item)
		gotErr := api.unmarshal(input, gotMsg)
		require.Equal(t, unknownFieldsAllowed, wantErr == nil, "protojson: %v", wantErr)
		require.Equal(t, wantErr == nil, gotErr == nil, "%s: protojson %v, generated %v", name, wantErr, gotErr)
		if gotErr == nil {
			require.True(t, proto.Equal(wantMsg, gotMsg), "%s: got %v, want %v", name, gotMsg, wantMsg)
		}
	}
}

func compact(t *testing.T, b []byte) string {
	t.Helper()
	v := jsontext.Value(append([]byte(nil), b...))
	require.NoError(t, v.Compact())
	return string(v)
}
