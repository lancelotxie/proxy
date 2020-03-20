package script

import "strings"

// format 格式化脚本
func format(script string) (formatted string) {
	// 统一符号
	script = formatSignal(script)
	// 合并空格
	script = mergeBlank(script)
	// 去除空行
	script = rmEmptyLine(script)

	formatted = script
	return
}

// formatSignal 格式化符号
func formatSignal(script string) (formatted string) {
	// 统一空格符
	script = strings.ReplaceAll(script, "\t", "    ")
	// 统一换行符
	formatted = strings.ReplaceAll(script, "\r\n", "\n")
	return
}

// mergeBlank 合并多余的空格
func mergeBlank(script string) (formatted string) {
	for strings.Contains(script, "  ") {
		script = strings.ReplaceAll(script, "  ", " ")
	}
	formatted = script
	return
}

// rmEmptyLine 去除空行
func rmEmptyLine(script string) (formatted string) {
	lines := strings.Split(script, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		if line == " " {
			continue
		}

		formatted += line + "\n"
	}

	return
}
