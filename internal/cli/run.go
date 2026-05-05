package cli

import (
	"fastgrep/internal/orchestrator"
	"fmt"
)

func Run() {
	searchString, files, ignoreCase, recursive, err := ParseArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(files) == 0 {
		return
	}

	err = orchestrator.Execute(searchString, files, ignoreCase, recursive)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
