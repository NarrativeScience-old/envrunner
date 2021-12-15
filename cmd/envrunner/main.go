package main

import (
	"os"

	"github.com/NarrativeScience/envrunner/internal/envrunner"
)

func main() {
	envrunner.Run(os.Args[1:])
}
