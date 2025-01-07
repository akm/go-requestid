package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		basic()
		return
	}
	switch os.Args[1] {
	case "basic":
		basic()
	case "namespace":
		namespace()
	}
}
