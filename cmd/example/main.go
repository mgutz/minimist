package main

import (
	"fmt"

	"github.com/mgutz/minimist"
)

var usage = `
example - an example os minimist

Usage:
  -h, --help     This screen
      --verbose  Log verbosely
  -v, --version  Print version
  -w, --watch    Watch tasks
`

func main() {
	argm := minimist.Parse()
	// cmd --help || cmd --h || cmd -?
	if argm.MayBool(false, "help", "h", "?") {
		fmt.Println(usage)
	}

	// cmd -v || cmd --version
	if argm.ZeroBool("v", "version") {
		fmt.Println("1.0")
	}

	// cmd foo -- ...
	// argm.SubCommand("foo", func(a *ArgMap) {
	// })

	// argm.SubExec("talk", "echo")
}
