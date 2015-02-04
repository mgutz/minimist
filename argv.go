package minimist

import (
	"fmt"

	"github.com/mgutz/go-nestedjson"
)

// argv is the result of parsing command-line arguments.
type argv struct {
	*nestedjson.Map
}

// NewFromMap creates a new ARGV instance from an existing map.
func newFromMap(m map[string]interface{}) *argv {
	return &argv{nestedjson.NewFromMap(m)}
}

// Leftover are arguments which were not parsed as flags before "--"
func (a *argv) Leftover() []interface{} {
	return a.MustArray("_")
}

// Unparsed are args that came after "--"
func (a *argv) Unparsed() []string {
	v, err := a.Get("--")
	if err != nil {
		fmt.Println(err.Error())
		panic(`"--" key is not a string slice`)
	}
	if slice, ok := v.([]string); ok {
		return slice
	}
	return nil
}
