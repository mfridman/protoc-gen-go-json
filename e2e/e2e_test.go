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
		require.IsType(t, []byte{}, val)

		got := new(Basic)
		require.NoError(t, got.Scan(val))
		require.Equal(t, original.GetA(), got.GetA())
		require.Equal(t, original.GetInt(), got.GetInt())
	})

	t.Run("scan string input", func(t *testing.T) {
		original := &Basic{A: "from-string", B: &Basic_Int{Int: 7}}
		val, err := original.Value()
		require.NoError(t, err)

		got := new(Basic)
		require.NoError(t, got.Scan(string(val.([]byte))))
		require.Equal(t, original.GetA(), got.GetA())
		require.Equal(t, original.GetInt(), got.GetInt())
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
}
