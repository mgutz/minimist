package minimist

import "github.com/mgutz/go-nestedjson"

// ArgMap is the result of parsing command-line arguments.
type ArgMap struct {
	*nestedjson.Map
}

// NewFromMap creates a new ARGV instance from an existing map.
func newFromMap(m map[string]interface{}) *ArgMap {
	return &ArgMap{nestedjson.NewFromMap(m)}
}

// Leftover are arguments which were not parsed as flags before "--"
func (am *ArgMap) Leftover() []interface{} {
	return am.MustArray("_")
}

// Unparsed are args that came after "--"
func (am *ArgMap) Unparsed() []string {
	v, err := am.Get("--")
	if err != nil {
		panic(`key "--" should not return an error`)
	}

	if slice, ok := v.([]string); ok {
		return slice
	}
	return nil
}

// SafeString should get value from path or return val.
func (am *ArgMap) SafeString(path, alias string, val string) string {
	v, err := am.Map.String(path)
	if err != nil {
		if alias == "" {
			return val
		}
		v, err = am.Map.String(alias)
		if err != nil {
			return val
		}
	}
	return v
}

// SafeInt should get value from path or return val.
func (am *ArgMap) SafeInt(path, alias string, val int) int {
	v, err := am.Map.Int(path)
	if err != nil {
		if alias == "" {
			return val
		}
		v, err = am.Map.Int(alias)
		if err != nil {
			return val
		}
	}
	return v
}

// SafeFloat should get value from path or return val.
func (am *ArgMap) SafeFloat(path, alias string, val float64) float64 {
	v, err := am.Map.Float(path)
	if err != nil {
		if alias == "" {
			return val
		}
		v, err = am.Map.Float(alias)
		if err != nil {
			return val
		}
	}
	return v
}

// SafeBool should get value from path or return val.
func (am *ArgMap) SafeBool(path, alias string, val bool) bool {
	v, err := am.Map.Bool(path)
	if err != nil {
		if alias == "" {
			return val
		}
		v, err = am.Map.Bool(alias)
		if err != nil {
			return val
		}
	}
	return v
}

// SafeArray should get value from path or return val.
func (am *ArgMap) SafeArray(path, alias string, val []interface{}) []interface{} {
	v, err := am.Map.Array(path)
	if err != nil {
		if alias == "" {
			return val
		}
		v, err = am.Map.Array(alias)
		if err != nil {
			return val
		}
	}

	return v
}

// SafeMap should get value from path or return val.
func (am *ArgMap) SafeMap(path, alias string, val map[string]interface{}) map[string]interface{} {
	v, err := am.Map.Map(path)
	if err != nil {
		if alias == "" {
			return val
		}
		v, err = am.Map.Map(alias)
		if err != nil {
			return val
		}
	}

	return v
}
