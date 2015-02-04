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
	res := parseArgs([]string{"--arg"}, nil)
	assert.Equal(true, res.MustBool("arg"))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg   # arg in strings
	res = parseArgs([]string{"--arg"}, &Options{Strings: []string{"arg"}})
	assert.Equal("", res.MustString("arg"))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg 1
	res = parseArgs([]string{"--arg", "1"}, nil)
	assert.Equal(1, res.MustInt("arg"))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg 1 # arg in bools
	res = parseArgs([]string{"--arg", "1"}, &Options{Bools: []string{"arg"}})
	assert.Equal(true, res.MustBool("arg"))
	assert.Equal(1, res.Leftover()[0].(int))

	// cmd --arg true # arg in bools
	res = parseArgs([]string{"--arg", "true"}, &Options{Bools: []string{"arg"}})
	assert.Equal(true, res.MustBool("arg"))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg -b
	res = parseArgs([]string{"--arg", "-b"}, nil)
	assert.Equal(true, res.MustBool("arg"))
	assert.Equal(true, res.MustBool("b"))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg=1
	res = parseArgs([]string{"--arg=1"}, nil)
	assert.Equal(1, res.MustInt("arg"))
	assert.Equal(0, len(res.Leftover()))

	// cmd --arg1 --arg2
	res = parseArgs([]string{"--arg1", "--arg2"}, nil)
	assert.Equal(true, res.MustBool("arg1"))
	assert.Equal(true, res.MustBool("arg2"))
}

func TestSingleDash(t *testing.T) {
	assert := assert.New(t)

	// cmd -a
	res := parseArgs([]string{"-a"}, nil)
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(0, len(res.Leftover()))

	// cmd -a1.24
	res = parseArgs([]string{"-a1.24"}, nil)
	assert.Equal(1.24, res.MustFloat("a"))
	assert.Equal(0, len(res.Leftover()))

	// cmd -ab1
	res = parseArgs([]string{"-ab1"}, nil)
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(1, res.MustInt("b"))
	assert.Equal(0, len(res.Leftover()))

	// # cmd -a码农
	res = parseArgs([]string{"-a码农"}, nil)
	assert.Equal("码农", res.MustString("a"))
	assert.Equal(0, len(res.Leftover()))

	// cmd -ab
	res = parseArgs([]string{"-ab"}, nil)
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(true, res.MustBool("b"))
	assert.Equal(0, len(res.Leftover()))

	// cmd -af test.py
	res = parseArgs([]string{"-af", "test.go"}, nil)
	assert.Equal(true, res.MustBool("a"))
	assert.Equal("test.go", res.MustString("f"))
	assert.Equal(0, len(res.Leftover()))

	// # cmd -af false  # f in bools
	res = parseArgs([]string{"-af", "false"}, &Options{Bools: []string{"f"}})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(false, res.MustBool("f"))
	assert.Equal(0, len(res.Leftover()))

	// cmd -af -b  # f in bools
	res = parseArgs([]string{"-af", "-b"}, &Options{Bools: []string{"f"}})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(true, res.MustBool("f"))
	assert.Equal(true, res.MustBool("b"))
	assert.Equal(0, len(res.Leftover()))
}

func TestLeftover(t *testing.T) {
	assert := assert.New(t)

	// cmd a b
	res := parseArgs([]string{"a", "b"}, nil)
	assert.Contains(res.Leftover(), "a")
	assert.Contains(res.Leftover(), "b")
	assert.Equal(2, len(res.Leftover()))

	// cmd -a b c d
	res = parseArgs([]string{"-a", "b", "c", "d"}, nil)
	assert.Equal("b", res.MustString("a"))
	assert.Contains(res.Leftover(), "c")
	assert.Contains(res.Leftover(), "d")
	assert.Equal(2, len(res.Leftover()))

}

func TestUnparsed(t *testing.T) {
	assert := assert.New(t)

	// cmd -a b c d -- -g --x
	res := parseArgs([]string{"-a", "b", "c", "d", "--", "-g", "--x"}, nil)
	assert.Equal("b", res.MustString("a"))
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
	res := parseArgs([]string{"--no-input"}, nil)
	assert.Equal(false, res.MustBool("input"))
	assert.Equal(0, len(res.Leftover()))
}

func TestDefaultsOption(t *testing.T) {
	assert := assert.New(t)

	// cmd -a2  # with a = 100 as default
	res := parseArgs([]string{"-a2"}, &Options{Defaults: map[string]interface{}{"a": 100}})
	assert.Equal(2, res.MustInt("a"))
	assert.Equal(0, len(res.Leftover()))

	// cmd -a  # with b = 2 as default
	res = parseArgs([]string{"-a"}, &Options{Defaults: map[string]interface{}{"b": 2}})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(2, res.MustInt("b"))
	assert.Equal(0, len(res.Leftover()))
}

func TestAliases(t *testing.T) {
	assert := assert.New(t)

	// cmd -a2  # with a = 100 as default
	res := parseArgs([]string{"-z2"}, NewOptions().Alias("z", "zm", "zoom"))
	assert.Equal(2, res.MustInt("z"))
	assert.Equal(2, res.MustInt("zm"))
	assert.Equal(2, res.MustInt("zoom"))
	assert.Equal(0, len(res.Leftover()))
}
