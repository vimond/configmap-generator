package configmap_generator

import (
	"os"
)

var Debug = false

func debug(msg string) {
	if Debug {
		os.Stderr.WriteString(msg)
	}
}
