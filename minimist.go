package minimist

import (
	"regexp"
	//"strings"
	"strconv"

	"github.com/mgutz/str"
)

// Args is the result of parsing command-line arguments.
type Args map[string]interface{}

// Leftover are arguments which were not parsed as flags before "--"
func (a Args) Leftover() []interface{} {
	return a["_"].([]interface{})
}

// Unparsed are args that came after "--"
func (a Args) Unparsed() []string {
	return a["--"].([]string)
}

func nextString(list []string, i int) *string {
	if i+1 < len(list) {
		return &list[i+1]
	}
	return nil
}

func sliceContains(strings []string, needle string) bool {
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

func parseArgs(args, bools, strings []string, defaults map[string]interface{}) Args {
	if bools == nil {
		bools = []string{}
	}

	if strings == nil {
		strings = []string{}
	}
	if defaults == nil {
		defaults = map[string]interface{}{}
	}

	inBools := func(key string) bool {
		return sliceContains(bools, key)
	}

	leftover := []interface{}{}
	result := map[string]interface{}{
		"_": leftover,
	}

	setArg := func(key string, val interface{}) {
		if !sliceContains(strings, key) && isNumber(val) {
			result[key] = parseNumber(val)
			return
		}
		result[key] = val
	}

	iifInStrings := func(key string, s string, v interface{}) interface{} {
		if sliceContains(strings, key) {
			return s
		}
		return v
	}

	l := len(args)
	argsAt := func(i int) string {
		if i > -1 && i < l {
			return args[i]
		}
		return ""
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == "--" {
			result["--"] = args[i+1:]
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

				if i+1 < len(args) {
					nextArg := args[i+1]
					if !dashesRe.MatchString(nextArg) && !inBools(key) {
						setArg(key, nextArg)
						i++
					} else if trueFalseRe.MatchString(args[i+1]) {
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

	for key := range defaults {
		if result[key] == nil {
			setArg(key, defaults[key])
		}
	}

	return result
}
