package main

import (
	"os"

	"enflame/enflame-container-toolkit/internal/runtime"
)

func main() {
	r := runtime.New()
	err := r.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
