package minimist

// ArgMap is the result of parsing command-line arguments.
type ArgMap map[string]interface{}

// Leftover are arguments which were not parsed as flags before "--"
func (am ArgMap) Leftover() []string {
	return am["_"].([]string)
}

// Unparsed are args that came after "--"
func (am ArgMap) Unparsed() []string {
	return am["--"].([]string)
}
