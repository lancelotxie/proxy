package script

import "strings"

func splitSentences(script string) (sentences []string) {
	sentences = strings.Split(script, "\n")
	return
}
