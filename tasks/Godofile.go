package main

import . "gopkg.in/godo.v1"

// Tasks is the entry point for Godo.
func Tasks(p *Project) {
	p.Task("test", W{"*.go"}, func() {
		Run("go test")
	})
}

func main() {
	Godo(Tasks)
}
