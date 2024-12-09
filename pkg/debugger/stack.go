package debugger

import (
	"os"
	"runtime/debug"
	"strings"
)

func Stack() {
	lines := strings.Split(string(debug.Stack()), "\n")
	var formattedLines []string
	for _, line := range lines {
		if strings.Contains(line, ".go") {
			formattedLines = append(formattedLines, line)
		}
	}

	os.Stderr.WriteString(strings.Join(formattedLines, "\n"))
}
