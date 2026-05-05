package cli

import (
	"flag"
	"fmt"
)

func ShowUsage() {
	fmt.Printf("\x1b[1;34m********* Fast Grep *********\x1b[0m\n")
	fmt.Printf("Usage: fastgrep [OPTIONS] <search-pattern> <file-or-directory>...\n\n")
	fmt.Printf("Options:\n")
	fmt.Printf("  -i    Ignore case in search\n")
	fmt.Printf("  -r    Search directories recursively\n")
	fmt.Printf("  -h    Show this help message\n\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  fastgrep \"hello\" main.go\n")
	fmt.Printf("  fastgrep -i \"error\" ./logs\n")
	fmt.Printf("  fastgrep -r \"TODO\" *\n")
}

func ParseArgs() (string, []string, bool, bool, error) {
	flag.Usage = ShowUsage
	ignoreCase := flag.Bool("i", false, "Ignore case in search")
	recursive := flag.Bool("r", false, "Search directories recursively")

	flag.Parse()
	args := flag.Args()

	var files = make([]string, 0)

	if len(args) > 0 && args[0] == "help" {
		ShowUsage()
		return "", files, false, false, nil
	}

	if len(args) < 2 {
		if len(args) == 0 {
			ShowUsage()
			return "", files, false, false, nil
		}
		return "", files, false, false, fmt.Errorf("invalid arguments passed: expected search pattern and at least one file/directory")
	}

	searchString := args[0]
	files = args[1:]
	return searchString, files, *ignoreCase, *recursive, nil
}
