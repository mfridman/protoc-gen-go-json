package e2e

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	type basicWrapper struct{ Basic }
	var optionalPresent = "present"
	var optionalEmpty = ""
	var cases = []struct {
		Name string

		// Value and Expected MUST be pointers to structs. If Expected is
		// nil, then it is expected to be identical to Value.
		Value    interface{}
		Expected interface{}
	}{
		{
			"basic",
			&Basic{
				A: "hello",
				B: &Basic_Int{
					Int: 42,
				},
			},
			nil,
		},

		{
			"basic wrapped in Go struct",
			&basicWrapper{
				Basic: Basic{
					A: "hello",
					B: &Basic_Int{
						Int: 42,
					},
				},
			},
			nil,
		},

		{
			"nested",
			&Nested_Message{
				Basic: &Basic{
					A: "hello",
					B: &Basic_Int{
						Int: 42,
					},
				},
			},
			nil,
		},

		{
			"optional present",
			&Basic{
				A: "hello",
				B: &Basic_Int{
					Int: 42,
				},
				O: &optionalPresent,
			},
			nil,
		},

		{
			"optional empty",
			&Basic{
				A: "hello",
				B: &Basic_Int{
					Int: 42,
				},
				O: &optionalEmpty,
			},
			nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			// Verify marshaling doesn't error
			bs, err := json.Marshal(tt.Value)
			require.NoError(err)
			require.NotEmpty(bs)

			// Determine what we expect the result to be
			expected := tt.Expected
			if expected == nil {
				expected = tt.Value
			}

			// Unmarshal. We want to do this into a concrete type so we
			// use reflection here (you can't just decode into interface{})
			// and have that work.
			val := reflect.New(reflect.ValueOf(expected).Elem().Type())
			require.NoError(json.Unmarshal(bs, val.Interface()))
			require.Equal(val.Interface(), expected)
		})
	}
}

func TestScannerValuer(t *testing.T) {
	t.Run("round trip", func(t *testing.T) {
		original := &Basic{
			A: "hello",
			B: &Basic_Int{Int: 42},
		}
		val, err := original.Value()
		require.NoError(t, err)
		require.IsType(t, "", val)

		got := new(Basic)
		require.NoError(t, got.Scan(val))
		require.Equal(t, original.GetA(), got.GetA())
		require.Equal(t, original.GetInt(), got.GetInt())
	})

	t.Run("scan string input", func(t *testing.T) {
		original := &Basic{A: "from-string", B: &Basic_Int{Int: 7}}
		val, err := original.Value()
		require.NoError(t, err)
		require.IsType(t, "", val)

		got := new(Basic)
		require.NoError(t, got.Scan(val))
		require.Equal(t, original.GetA(), got.GetA())
		require.Equal(t, original.GetInt(), got.GetInt())
	})

	t.Run("scan bytes input", func(t *testing.T) {
		got := new(Basic)
		require.NoError(t, got.Scan([]byte(`{"a":"from-bytes","int":9}`)))
		require.Equal(t, "from-bytes", got.GetA())
		require.Equal(t, int32(9), got.GetInt())
	})

	t.Run("scan nil resets message", func(t *testing.T) {
		optional := "set"
		msg := &Basic{
			A: "hello",
			B: &Basic_Int{Int: 42},
			O: &optional,
		}
		require.NoError(t, msg.Scan(nil))
		require.Equal(t, &Basic{}, msg)
	})

	t.Run("scan unsupported type returns error", func(t *testing.T) {
		msg := new(Basic)
		err := msg.Scan(12345)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected []byte or string")
	})

	t.Run("nil receiver scan returns error", func(t *testing.T) {
		var msg *Basic
		err := msg.Scan([]byte(`{"a":"x"}`))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nil receiver")
	})

	t.Run("nil receiver value returns nil", func(t *testing.T) {
		var msg *Basic
		val, err := msg.Value()
		require.NoError(t, err)
		require.Nil(t, val)
	})

	t.Run("valuer interface", func(t *testing.T) {
		var _ driver.Valuer = (*Basic)(nil)
		var _ driver.Valuer = (*Nested)(nil)
		var _ driver.Valuer = (*Nested_Message)(nil)
	})

	t.Run("nested message round trip", func(t *testing.T) {
		original := &Nested_Message{
			Basic: &Basic{
				A: "nested",
				B: &Basic_Str{Str: "value"},
			},
		}
		val, err := original.Value()
		require.NoError(t, err)

		got := new(Nested_Message)
		require.NoError(t, got.Scan(val))
		require.Equal(t, original.GetBasic().GetA(), got.GetBasic().GetA())
		require.Equal(t, original.GetBasic().GetStr(), got.GetBasic().GetStr())
	})

	t.Run("map fields round trip", func(t *testing.T) {
		original := &Basic{
			A:   "with-map",
			Map: map[string]string{"key1": "val1", "key2": "val2"},
		}
		val, err := original.Value()
		require.NoError(t, err)

		got := new(Basic)
		require.NoError(t, got.Scan(val))
		require.Equal(t, original.GetMap(), got.GetMap())
	})

	t.Run("optional fields round trip", func(t *testing.T) {
		present := "present"
		original := &Basic{
			A: "with-optional",
			B: &Basic_Str{Str: "test"},
			O: &present,
		}
		val, err := original.Value()
		require.NoError(t, err)

		got := new(Basic)
		require.NoError(t, got.Scan(val))
		require.NotNil(t, got.O)
		require.Equal(t, present, *got.O)
	})

	t.Run("optional field absent round trip", func(t *testing.T) {
		original := &Basic{
			A: "no-optional",
			B: &Basic_Int{Int: 1},
		}
		val, err := original.Value()
		require.NoError(t, err)

		got := new(Basic)
		require.NoError(t, got.Scan(val))
		require.Nil(t, got.O)
	})

	t.Run("oneof string variant round trip", func(t *testing.T) {
		original := &Basic{
			A: "oneof-str",
			B: &Basic_Str{Str: "hello"},
		}
		val, err := original.Value()
		require.NoError(t, err)

		got := new(Basic)
		require.NoError(t, got.Scan(val))
		require.Equal(t, "hello", got.GetStr())
		require.Equal(t, int32(0), got.GetInt())
	})

	t.Run("all fields populated round trip", func(t *testing.T) {
		optional := "opt"
		original := &Basic{
			A:   "full",
			B:   &Basic_Int{Int: 99},
			Map: map[string]string{"a": "1", "b": "2"},
			O:   &optional,
		}
		val, err := original.Value()
		require.NoError(t, err)

		got := new(Basic)
		require.NoError(t, got.Scan(val))
		require.Equal(t, original.GetA(), got.GetA())
		require.Equal(t, original.GetInt(), got.GetInt())
		require.Equal(t, original.GetMap(), got.GetMap())
		require.NotNil(t, got.O)
		require.Equal(t, optional, *got.O)
	})

	t.Run("reused receiver with partial json clears omitted fields", func(t *testing.T) {
		optional := "set"
		msg := &Basic{
			A:   "seed",
			B:   &Basic_Int{Int: 42},
			Map: map[string]string{"k": "v"},
			O:   &optional,
		}

		require.NoError(t, msg.Scan([]byte(`{"a":"only-a"}`)))
		require.Equal(t, "only-a", msg.GetA())
		require.Nil(t, msg.GetB())
		require.Equal(t, int32(0), msg.GetInt())
		require.Nil(t, msg.GetMap())
		require.Nil(t, msg.O)
	})

	t.Run("reused receiver oneof switches variants", func(t *testing.T) {
		msg := new(Basic)

		require.NoError(t, msg.Scan([]byte(`{"int":7}`)))
		require.Equal(t, int32(7), msg.GetInt())
		require.Equal(t, "", msg.GetStr())

		require.NoError(t, msg.Scan([]byte(`{"str":"seven"}`)))
		require.Equal(t, "seven", msg.GetStr())
		require.Equal(t, int32(0), msg.GetInt())

		require.NoError(t, msg.Scan([]byte(`{"int":11}`)))
		require.Equal(t, int32(11), msg.GetInt())
		require.Equal(t, "", msg.GetStr())
	})

	t.Run("json null literal differs from sql null", func(t *testing.T) {
		optional := "set"
		msg := &Basic{
			A:   "seed",
			B:   &Basic_Int{Int: 42},
			Map: map[string]string{"k": "v"},
			O:   &optional,
		}

		require.NoError(t, msg.Scan(nil))
		require.Equal(t, &Basic{}, msg)

		msg = &Basic{
			A:   "seed",
			B:   &Basic_Int{Int: 42},
			Map: map[string]string{"k": "v"},
			O:   &optional,
		}
		err := msg.Scan([]byte("null"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected token null")
		require.Equal(t, "", msg.GetA())
		require.Nil(t, msg.GetB())
		require.Nil(t, msg.GetMap())
		require.Nil(t, msg.O)

		msg = &Basic{
			A:   "seed",
			B:   &Basic_Int{Int: 42},
			Map: map[string]string{"k": "v"},
			O:   &optional,
		}
		err = msg.Scan("null")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected token null")
		require.Equal(t, "", msg.GetA())
		require.Nil(t, msg.GetB())
		require.Nil(t, msg.GetMap())
		require.Nil(t, msg.O)
	})

	t.Run("scan unknown field returns error", func(t *testing.T) {
		msg := new(Basic)
		err := msg.Scan([]byte(`{"a":"x","unknown":1}`))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown field")
	})

	t.Run("scan invalid payloads return error", func(t *testing.T) {
		msg := new(Basic)
		require.Error(t, msg.Scan([]byte(`{"map":"not-object"}`)))
		require.Error(t, msg.Scan([]byte(`{"int":"abc"}`)))
		require.Error(t, msg.Scan([]byte(`{"a":"x"`)))
	})

	t.Run("scan error allows subsequent successful scan on reused receiver", func(t *testing.T) {
		msg := &Basic{
			A: "seed",
			B: &Basic_Int{Int: 42},
		}
		require.Error(t, msg.Scan([]byte(`{"int":"abc"}`)))
		require.NoError(t, msg.Scan([]byte(`{"a":"ok","int":3}`)))
		require.Equal(t, "ok", msg.GetA())
		require.Equal(t, int32(3), msg.GetInt())
	})

	t.Run("nil receiver scan nil source returns error", func(t *testing.T) {
		var msg *Basic
		err := msg.Scan(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nil receiver")
	})
}
