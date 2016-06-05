package internal

import (
	"fmt"
)

func Debug(data []byte, err error) {
	if !DebugLogging {
		return
	}
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", data)
	}
}
