package utils

import "strings"

func FormatPythonCode(raw string) string {
	// Strip the triple backtick markers
	code := strings.TrimSpace(raw)
	code = strings.TrimPrefix(code, "```python")
	code = strings.TrimPrefix(code, "```") // in case it's just ```
	code = strings.TrimSuffix(code, "```")
	code = strings.TrimSpace(code)

	lines := strings.Split(code, "\n")
	var formatted []string
	indent := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Adjust indentation level
		if strings.HasPrefix(trimmed, "elif") || strings.HasPrefix(trimmed, "else") || strings.HasPrefix(trimmed, "except") || strings.HasPrefix(trimmed, "finally") {
			indent--
		}
		if strings.HasPrefix(trimmed, "#") {
			// Preserve comment position
			formatted = append(formatted, strings.Repeat("    ", indent)+trimmed)
			continue
		}
		formatted = append(formatted, strings.Repeat("    ", indent)+trimmed)

		// Increase indent after colons that start a block
		if strings.HasSuffix(trimmed, ":") &&
			!strings.HasPrefix(trimmed, "#") &&
			!strings.Contains(trimmed, "lambda") {
			indent++
		}

		// Reduce indentation for dedent keywords (crude but practical)
		if trimmed == "" {
			indent = 0
		}
	}

	return strings.Join(formatted, "\n")
}
