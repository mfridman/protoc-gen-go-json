package e2e

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// The generated streaming methods take precedence over MarshalJSON and
// UnmarshalJSON in encoding/json v1 and v2.
var (
	_ jsonv2.MarshalerTo     = (*Basic)(nil)
	_ jsonv2.UnmarshalerFrom = (*Basic)(nil)
	_ jsonv2.MarshalerTo     = (*Nested_Message)(nil)
	_ jsonv2.UnmarshalerFrom = (*Nested_Message)(nil)
)

func TestJSONV2(t *testing.T) {
	t.Parallel()

	present := "present"
	msg := &Basic{
		A:   "hello",
		B:   &Basic_Int{Int: 42},
		Map: map[string]string{"b": "2", "a": "1"},
		O:   &present,
	}
	type container struct {
		ID      string `json:"id"`
		Message *Basic `json:"message"`
	}
	in := container{ID: "x", Message: msg}

	// The output must be the same as protojson with the options from
	// e2e/buf.gen.yaml.
	want, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	require.NoError(t, err)
	wantMsg, err := jsonv2.Marshal(jsonv1.RawMessage(want))
	require.NoError(t, err)

	for name, api := range map[string]struct {
		marshal   func(any) ([]byte, error)
		unmarshal func([]byte, any) error
	}{
		"v1": {jsonv1.Marshal, jsonv1.Unmarshal},
		"v2": {
			func(v any) ([]byte, error) { return jsonv2.Marshal(v) },
			func(b []byte, v any) error { return jsonv2.Unmarshal(b, v) },
		},
	} {
		got, err := api.marshal(in)
		require.NoError(t, err, name)
		require.JSONEq(t, `{"id":"x","message":`+string(wantMsg)+`}`, string(got), name)

		var out container
		require.NoError(t, api.unmarshal(got, &out), name)
		require.True(t, proto.Equal(msg, out.Message), "%s: got %v", name, out.Message)
	}
}
