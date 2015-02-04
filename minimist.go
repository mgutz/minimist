package minimist

import (
	"os"
	"regexp"
	//"strings"
	"strconv"

	"github.com/mgutz/str"
)

func nextString(list []string, i int) *string {
	if i+1 < len(list) {
		return &list[i+1]
	}
	return nil
}

func sliceContains(strings []string, needle string) bool {
	if strings == nil {
		return false
	}
	for _, s := range strings {
		if s == needle {
			return true
		}
	}
	return false
}

var integerRe = regexp.MustCompile(`^-?\d+$`)
var numberRe = regexp.MustCompile(`^-?\d+(\.\d+)?(e-?\d+)?$`)

// --port=8000
var longFormEqualRe = regexp.MustCompile(`^--.+=`)
var longFormEqualValsRe = regexp.MustCompile(`^--([^=]+)=(.*)$`)

// --port 8000
var longFormRe = regexp.MustCompile(`^--.+`)
var longFormKeyRe = regexp.MustCompile(`^--(.+)`)

//longFormSpaceValsRe := regexp.MustCompile(`^--([^=])=([\s\S]*)$`)

// --no-debug
var negateRe = regexp.MustCompile(`^--no-.+`)
var negateValsRe = regexp.MustCompile(`^--no-(.+)`)

// -abc
var shortFormRe = regexp.MustCompile(`^-[^-]+`)

var lettersRe = regexp.MustCompile(`^[A-Za-z]`)

var notWordRe = regexp.MustCompile(`\W`)

var dashesRe = regexp.MustCompile(`^(-|--)`)

var trueFalseRe = regexp.MustCompile(`^(true|false)`)

func isNumber(v interface{}) bool {
	return isInteger(v) || isFloat(v)
}

func isFloat(v interface{}) bool {
	switch t := v.(type) {
	case float64, float32:
		return true
	case string:
		return numberRe.MatchString(t)
	default:
		return false
	}
}

func toFloat64(v interface{}) float64 {
	switch t := v.(type) {
	case float64, float32:
		return t.(float64)
	case string:
		val, _ := strconv.ParseFloat(t, 64)
		return val
	default:
		return 0
	}
}

func isInteger(v interface{}) bool {
	switch t := v.(type) {
	case int32, int64, uint, uint32, uint64:
		return true
	case string:
		return integerRe.MatchString(t)
	default:
		return false
	}
}
func toInt(v interface{}) int {
	switch t := v.(type) {
	case int32, int64, uint, uint32, uint64:
		return t.(int)
	case string:
		val, _ := strconv.Atoi(t)
		return val
	default:
		return 0
	}
}

func parseNumber(v interface{}) interface{} {
	if isInteger(v) {
		return toInt(v)
	}
	return toFloat64(v)
}

// Options are parse options.
type Options struct {
	// Bools are flags which should always be treated as bool
	Bools []string
	// Strings are flags which should always be treated as string
	Strings []string
	// Defaults are default values for flags.
	Defaults map[string]interface{}
	// Aliases are aliases for a flag
	Aliases map[string][]string
}

// NewOptions creates an instance of Options.
func NewOptions() *Options {
	return &Options{}
}

// Alias adds an alias.
func (o *Options) Alias(from string, to ...string) *Options {
	if len(to) == 0 {
		return o
	}

	if o.Aliases == nil {
		o.Aliases = map[string][]string{}
	}
	o.Aliases[from] = to
	return o
}

// Parse parses os.Args excluding os.Args[0].
func Parse() *argv {
	return parseArgs(os.Args[1:], nil)
}

// ParseArgs parses an argv for options.
func parseArgs(argv []string, options *Options) *argv {
	if options == nil {
		options = &Options{}
	}
	bools := options.Bools
	strings := options.Strings
	defaults := options.Defaults
	aliases := options.Aliases

	inBools := func(key string) bool {
		return sliceContains(bools, key)
	}

	leftover := []interface{}{}
	result := map[string]interface{}{
		"_": leftover,
	}

	setArg := func(key string, val interface{}) {
		var keys []string
		if aliases != nil {
			if aka := aliases[key]; aka != nil {
				keys = append(aka, key)
			}
		}

		if keys == nil {
			keys = []string{key}
		}

		for _, keyName := range keys {
			if !sliceContains(strings, keyName) && isNumber(val) {
				result[keyName] = parseNumber(val)
				continue
			}
			result[keyName] = val
		}

	}

	iifInStrings := func(key string, s string, v interface{}) interface{} {
		if sliceContains(strings, key) {
			return s
		}
		return v
	}

	l := len(argv)
	argsAt := func(i int) string {
		if i > -1 && i < l {
			return argv[i]
		}
		return ""
	}

	i := 0
	for i < len(argv) {
		arg := argv[i]

		if arg == "--" {
			result["--"] = argv[i+1:]
			break
		}

		argAt := func(i int) string {
			if i >= 0 && i < len(arg) {
				return arg[i : i+1]
			}
			return ""
		}
		if longFormEqualRe.MatchString(arg) {
			// --long-form=value

			m := longFormEqualValsRe.FindStringSubmatch(arg)
			//fmt.Printf("--long-form= %s\n", arg)
			setArg(m[1], m[2])

		} else if negateRe.MatchString(arg) {
			//fmt.Printf("--no-flag %s\n", arg)

			m := negateValsRe.FindStringSubmatch(arg)
			setArg(m[1], false)

		} else if longFormRe.MatchString(arg) {
			// --long-form
			//fmt.Printf("--long-form %s\n", arg)

			key := longFormKeyRe.FindStringSubmatch(arg)[1]
			next := argsAt(i + 1)

			if next == "" {
				setArg(key, iifInStrings(key, "", true))
			} else if next[0:1] == "-" {
				setArg(key, iifInStrings(key, "", true))
			} else if !sliceContains(bools, key) {
				setArg(key, next)
				i++
			} else if next == "true" || next == "false" {
				setArg(key, next == "true")
				i++
			} else {
				setArg(key, true)
			}
		} else if shortFormRe.MatchString(arg) {
			// -abc a, b are boolean c is undetermined
			//fmt.Printf("-short-form %s\n", arg)

			letters := arg[1:]

			L := len(letters)
			lettersAt := func(i int) string {
				if i < L {
					return letters[i : i+1]
				}
				return ""
			}

			broken := false
			k := 0
			for k < len(letters) {
				next := arg[k+2:]
				if next == "-" {
					setArg(lettersAt(k), next)
					k++
					continue
				}
				if lettersRe.MatchString(lettersAt(k)) && numberRe.MatchString(next) {
					setArg(lettersAt(k), next)
					broken = true
					break
				}
				if k+1 < len(letters) && notWordRe.MatchString(lettersAt(k+1)) {
					setArg(lettersAt(k), next)
					broken = true
					break
				}

				setArg(lettersAt(k), iifInStrings(lettersAt(k), "", true))
				k++
			}

			key := argAt(len(arg) - 1)
			if !broken && key != "-" {

				if i+1 < len(argv) {
					nextArg := argv[i+1]
					if !dashesRe.MatchString(nextArg) && !inBools(key) {
						setArg(key, nextArg)
						i++
					} else if trueFalseRe.MatchString(argv[i+1]) {
						setArg(key, nextArg == "true")
						i++
					}
				} else {
					setArg(key, iifInStrings(key, "", true))
				}
			}
		} else {
			if str.IsNumeric(arg) {
				leftover = append(leftover, parseNumber(arg))
			} else {
				leftover = append(leftover, arg)
			}
			result["_"] = leftover
		}

		i++
	}

	if defaults != nil {
		for key := range defaults {
			if result[key] == nil {
				setArg(key, defaults[key])
			}
		}
	}

	return newFromMap(result)
}
