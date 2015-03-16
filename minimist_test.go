package minimist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var emptyStrings = []string{}
var emptyMap = map[string]interface{}{}

func TestDoubleDash(t *testing.T) {
	assert := assert.New(t)

	// TODO
	// must = required
	// should = recommended
	// may = optional

	// MustBool("help", "h", "?")
	// MayBool(false, "help", "h", "?")

	// cmd --arg
	res := ParseArgv([]string{"--arg"})
	assert.Equal(true, res.MayBool(false, "arg"))
	assert.Equal(0, len(res.Rest()))

	// cmd --arg   # arg in strings
	res = ParseArgv([]string{"--arg"})
	assert.Equal("true", res.MayString("", "arg"))
	assert.Equal(0, len(res.Rest()))

	// cmd --arg 1
	res = ParseArgv([]string{"--arg", "1"})
	assert.Equal(1, res.MayInt(100, "arg"))
	assert.Equal(0, len(res.Rest()))

	// cmd --arg true
	res = ParseArgv([]string{"--arg=true"})
	assert.Equal("true", res.MustString("arg"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"--arg", "true"})
	assert.Equal("true", res.MustString("arg"))
	assert.Equal(0, len(res.Rest()))

	// cmd --arg 1
	res = ParseArgv([]string{"--arg", "0"})
	assert.Equal(false, res.MustBool("arg"))
	assert.Equal(0, len(res.Rest()))

	// cmd --arg false
	res = ParseArgv([]string{"--arg", "false"})
	assert.Equal(false, res.MustBool("arg"))
	assert.Equal(0, len(res.Rest()))

	// cmd --arg -b
	res = ParseArgv([]string{"--arg", "-b"})
	assert.Equal(true, res.MustBool("arg"))
	assert.Equal(true, res.MustBool("b"))
	assert.Equal(0, len(res.Rest()))

	// cmd --arg=1
	res = ParseArgv([]string{"--arg=1"})
	assert.Equal(1, res.MustInt("arg"))
	assert.Equal(0, len(res.Rest()))

	// cmd --arg1 --arg2
	res = ParseArgv([]string{"--arg1", "--arg2"})
	assert.Equal(true, res.MustBool("arg1"))
	assert.Equal(true, res.MustBool("arg2"))
}

func TestSingleDash(t *testing.T) {
	assert := assert.New(t)

	// cmd -a
	res := ParseArgv([]string{"-a"})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(0, len(res.Rest()))

	// cmd -a foo
	res = ParseArgv([]string{"-a", "foo"})
	assert.Equal("foo", res.MustString("a"))
	assert.Equal(0, len(res.Rest()))

	// cmd -a1.24
	res = ParseArgv([]string{"-a1.24"})
	assert.Equal(1.24, res.MustFloat("a"))
	assert.Equal(0, len(res.Rest()))

	// cmd -ab1
	res = ParseArgv([]string{"-ab1"})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(1, res.MustInt("b"))
	assert.Equal(0, len(res.Rest()))

	// # cmd -a码农
	res = ParseArgv([]string{"-a码农"})
	assert.Equal("码农", res.MustString("a"))
	assert.Equal(0, len(res.Rest()))

	// cmd -ab
	res = ParseArgv([]string{"-ab"})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(true, res.MustBool("b"))
	assert.Equal(0, len(res.Rest()))

	// cmd -af test.py
	res = ParseArgv([]string{"-af", "test.go"})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal("test.go", res.MustString("f"))
	assert.Equal(0, len(res.Rest()))

	// # cmd -af false  # f in bools
	res = ParseArgv([]string{"-af", "false"})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(false, res.MustBool("f"))
	assert.Equal(0, len(res.Rest()))

	// cmd -af -b  # f in bools
	res = ParseArgv([]string{"-af", "-b"})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(true, res.MustBool("f"))
	assert.Equal(true, res.MustBool("b"))
	assert.Equal(0, len(res.Rest()))
}

func TestRest(t *testing.T) {
	assert := assert.New(t)

	// cmd a b
	res := ParseArgv([]string{"a", "b"})
	assert.Contains(res.Rest(), "a")
	assert.Contains(res.Rest(), "b")
	assert.Equal(2, len(res.Rest()))

	// cmd -a b c d
	res = ParseArgv([]string{"-a", "b", "c", "d"})
	assert.Equal("b", res.MustString("a"))
	assert.Contains(res.Rest(), "c")
	assert.Contains(res.Rest(), "d")
	assert.Equal(2, len(res.Rest()))

}

func TestUnparsed(t *testing.T) {
	assert := assert.New(t)

	// cmd -a b c d -- -g --x
	res := ParseArgv([]string{"-a", "b", "c", "d", "--", "-g", "--x"})
	assert.Equal("b", res.MustString("a"))
	assert.Contains(res.Rest(), "c")
	assert.Contains(res.Rest(), "d")
	assert.Contains(res.Unparsed(), "-g")
	assert.Contains(res.Unparsed(), "--x")
	assert.Equal(2, len(res.Rest()))
	assert.Equal(2, len(res.Unparsed()))

	// cmd -z2 -- foo
	res = ParseArgv([]string{"-z2", "--", "foo"})
	assert.Equal(2, res.MustInt("z"))
	assert.Equal(0, len(res.Rest()))
	assert.Contains(res.Unparsed(), "foo")
	assert.Equal(1, len(res.Unparsed()))

	// cmd -z2 # no unparsed args
	res = ParseArgv([]string{"-z2"})
	assert.Equal(0, len(res.Unparsed()))

}

func TestNegate(t *testing.T) {
	assert := assert.New(t)
	// cmd --no-input
	res := ParseArgv([]string{"--no-input"})
	assert.Equal(false, res.MustBool("input"))
	assert.Equal(0, len(res.Rest()))
}

func TestDefaultsOption(t *testing.T) {
	assert := assert.New(t)

	// cmd -a2  # with a = 100 as default
	res := ParseArgv([]string{"-a2"})
	assert.Equal(2, res.MayInt(100, "a"))
	assert.Equal(0, len(res.Rest()))

	// cmd -a
	res = ParseArgv([]string{"-a"})
	assert.Equal(true, res.MustBool("a"))
	assert.Equal(2, res.MayInt(2, "b"))
	assert.Equal(0, len(res.Rest()))
}

func TestAliases(t *testing.T) {
	assert := assert.New(t)

	// cmd -z2
	res := ParseArgv([]string{"-z2"})
	assert.Equal(2, res.MustInt("zoom", "zm", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"--zm", "3"})
	assert.Equal(3, res.MustInt("zoom", "zm", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"--zoom", "4"})
	assert.Equal(4, res.MustInt("zoom", "zm", "z"))
	assert.Equal(0, len(res.Rest()))
}

func TestMayFuncs(t *testing.T) {
	assert := assert.New(t)

	// cmd -z2
	res := ParseArgv([]string{"-z2"})
	assert.Equal(2, res.MayInt(100, "zoom", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"--zoom=2"})
	assert.Equal(2, res.MayInt(100, "zoom", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"-z", "bird"})
	assert.Equal("bird", res.MustString("zoom", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"--zoom=bird"})
	assert.Equal("bird", res.MustString("zoom", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"-z"})
	assert.Equal(true, res.MayBool(false, "zoom", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"--zoom"})
	assert.Equal(true, res.MayBool(false, "zoom", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"--zom"})
	assert.Equal(false, res.MayBool(false, "zoom", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"-z1.0"})
	assert.Equal(1.0, res.MayFloat(2.0, "zoom", "z"))
	assert.Equal(0, len(res.Rest()))

	res = ParseArgv([]string{"--zoom=1.0"})
	assert.Equal(1.0, res.MayFloat(2.0, "zoom", "z"))
	assert.Equal(0, len(res.Rest()))
}

func TestRestString(t *testing.T) {
	argm := ParseArgv([]string{"--zoom=1.0", "--", "one two"})
	argu := ParseArgv(argm.Unparsed())
	assert.Equal(t, "one two", argu.Rest()[0])
}
