package minimist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var emptyStrings = []string{}
var emptyMap = map[string]interface{}{}

func TestDoubleDash(t *testing.T) {
	assert := assert.New(t)

	// cmd --arg
	res := parseArgs([]string{"--arg"}, nil, nil, nil)
	assert.Equal(true, res["arg"].(bool))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg   # arg in strings
	res = parseArgs([]string{"--arg"}, nil, []string{"arg"}, nil)
	assert.Equal("", res["arg"].(string))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg 1
	res = parseArgs([]string{"--arg", "1"}, nil, nil, nil)
	assert.Equal(1, res["arg"].(int))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg 1 # arg in bools
	res = parseArgs([]string{"--arg", "1"}, []string{"arg"}, nil, nil)
	assert.Equal(true, res["arg"].(bool))
	assert.Equal(1, res.Leftover()[0].(int))

	// cmd --arg true # arg in bools
	res = parseArgs([]string{"--arg", "true"}, []string{"arg"}, nil, nil)
	assert.Equal(true, res["arg"].(bool))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg -b
	res = parseArgs([]string{"--arg", "-b"}, nil, nil, nil)
	assert.Equal(true, res["arg"].(bool))
	assert.Equal(true, res["b"].(bool))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg=1
	res = parseArgs([]string{"--arg=1"}, nil, nil, nil)
	assert.Equal(1, res["arg"].(int))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg1 --arg2
	res = parseArgs([]string{"--arg1", "--arg2"}, nil, nil, nil)
	assert.Equal(true, res["arg1"].(bool))
	assert.Equal(true, res["arg2"].(bool))
}

func TestSingleDash(t *testing.T) {
	assert := assert.New(t)

	// cmd -a
	res := parseArgs([]string{"-a"}, nil, nil, nil)
	assert.Equal(true, res["a"].(bool))
	assert.Equal(0, len(res.Leftover()))

	// cmd -a1.24
	res = parseArgs([]string{"-a1.24"}, nil, nil, nil)
	assert.Equal(1.24, res["a"].(float64))
	assert.Equal(0, len(res.Leftover()))

	// cmd -ab1
	res = parseArgs([]string{"-ab1"}, nil, nil, nil)
	assert.Equal(true, res["a"].(bool))
	assert.Equal(1, res["b"].(int))
	assert.Equal(0, len(res.Leftover()))

	// # cmd -a码农
	res = parseArgs([]string{"-a码农"}, nil, nil, nil)
	assert.Equal("码农", res["a"].(string))
	assert.Equal(0, len(res.Leftover()))

	// cmd -ab
	res = parseArgs([]string{"-ab"}, nil, nil, nil)
	assert.Equal(true, res["a"].(bool))
	assert.Equal(true, res["b"].(bool))
	assert.Equal(0, len(res.Leftover()))

	// cmd -af test.py
	res = parseArgs([]string{"-af", "test.go"}, nil, nil, nil)
	assert.Equal(true, res["a"].(bool))
	assert.Equal("test.go", res["f"].(string))
	assert.Equal(0, len(res.Leftover()))

	// # cmd -af false  # f in bools
	res = parseArgs([]string{"-af", "false"}, []string{"f"}, nil, nil)
	assert.Equal(true, res["a"].(bool))
	assert.Equal(false, res["f"].(bool))
	assert.Equal(0, len(res.Leftover()))

	// cmd -af -b  # f in bools
	res = parseArgs([]string{"-af", "-b"}, []string{"f"}, nil, nil)
	assert.Equal(true, res["a"].(bool))
	assert.Equal(true, res["f"].(bool))
	assert.Equal(true, res["b"].(bool))
	assert.Equal(0, len(res.Leftover()))
}

func TestLeftover(t *testing.T) {
	assert := assert.New(t)

	// cmd a b
	res := parseArgs([]string{"a", "b"}, nil, nil, nil)
	assert.Contains(res.Leftover(), "a")
	assert.Contains(res.Leftover(), "b")
	assert.Equal(2, len(res.Leftover()))

	// cmd -a b c d
	res = parseArgs([]string{"-a", "b", "c", "d"}, nil, nil, nil)
	assert.Equal("b", res["a"].(string))
	assert.Contains(res.Leftover(), "c")
	assert.Contains(res.Leftover(), "d")
	assert.Equal(2, len(res.Leftover()))

}

func TestUnparsed(t *testing.T) {
	assert := assert.New(t)

	// cmd -a b c d -- -g --x
	res := parseArgs([]string{"-a", "b", "c", "d", "--", "-g", "--x"}, nil, nil, nil)
	assert.Equal("b", res["a"].(string))
	assert.Contains(res.Leftover(), "c")
	assert.Contains(res.Leftover(), "d")
	assert.Contains(res.Unparsed(), "-g")
	assert.Contains(res.Unparsed(), "--x")
	assert.Equal(2, len(res.Leftover()))
	assert.Equal(2, len(res.Unparsed()))
}

func TestNegate(t *testing.T) {
	assert := assert.New(t)
	// cmd --no-input
	res := parseArgs([]string{"--no-input"}, nil, nil, nil)
	assert.Equal(false, res["input"].(bool))
	assert.Equal(0, len(res.Leftover()))
}

func TestDefaultsOption(t *testing.T) {
	assert := assert.New(t)

	// cmd -a2  # with a = 100 as default
	res := parseArgs([]string{"-a2"}, nil, nil, map[string]interface{}{"a": 100})
	assert.Equal(2, res["a"].(int))
	assert.Equal(0, len(res.Leftover()))

	// cmd -a  # with b = 2 as default
	res = parseArgs([]string{"-a"}, nil, nil, map[string]interface{}{"b": 2})
	assert.Equal(true, res["a"].(bool))
	assert.Equal(2, res["b"].(int))
	assert.Equal(0, len(res.Leftover()))
}
