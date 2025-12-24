package internal

import (
	"runtime"
	"strings"
)

type StackTracer interface {
	StackTrace() string
}

// stack captures the current stack trace
func stack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// formatStack formats the stack trace for better readability
func formatStack(stack string) string {
	var sb strings.Builder
	lines := strings.Split(stack, "\n")

	// Skip the first line (it's just "goroutine X [running]:")
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		sb.WriteString(lines[i])
		if i < len(lines)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
