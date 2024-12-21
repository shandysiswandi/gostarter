package debugger

import (
	"os"
	"runtime/debug"
	"strings"
)

func Stack(prefix string) {
	lines := strings.Split(string(debug.Stack()), "\n")
	var formattedLines []string
	for _, line := range lines {
		if strings.Contains(line, ".go") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, prefix) {
				formattedLines = append(formattedLines, strings.TrimPrefix(line, prefix))
			}
		}
	}

	os.Stderr.WriteString(strings.Join(formattedLines, "\n"))
}
