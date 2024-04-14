package plugin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mfridman/protoc-gen-go-json/gen"
)

var supportedOptions = map[string]func(*gen.Options, string) error{
	// Marshal options
	"enums_as_ints":              func(o *gen.Options, value string) error { return parseBool(&o.EnumsAsInts, value) },
	"emit_defaults":              func(o *gen.Options, value string) error { return parseBool(&o.EmitDefaults, value) },
	"emit_defaults_without_null": func(o *gen.Options, value string) error { return parseBool(&o.EmitDefaultValues, value) },
	"orig_name":                  func(o *gen.Options, value string) error { return parseBool(&o.OrigName, value) },
	// Unmarshal options
	"allow_unknown": func(o *gen.Options, value string) error { return parseBool(&o.AllowUnknownFields, value) },
}

func parseOptions(raw string) (*gen.Options, error) {
	opts := new(gen.Options)
	if raw == "" {
		return opts, nil
	}
	all := strings.Split(raw, ",")
	for _, opt := range all {
		name, value, ok := strings.Cut(opt, "=")
		if !ok {
			return nil, fmt.Errorf("invalid option, must be in the form of name=value: %s", opt)
		}
		// Check if the option is supported, but don't error.
		if fn, ok := supportedOptions[name]; ok {
			if err := fn(opts, value); err != nil {
				return nil, fmt.Errorf("invalid value for %s: %w", name, err)
			}
		}
	}
	return opts, nil
}

func parseBool(target *bool, value string) error {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	*target = b
	return nil
}
