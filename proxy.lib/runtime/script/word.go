package script

import (
	"strings"

	"github.com/lancelotXie/proxy/proxy.lib/container/set"
)

// WordType 单词的类型
type WordType int

const (
	// Operator 操作符
	Operator WordType = iota
	// VarName 变量名
	VarName
	// Const 常量
	Const
)

// definedWords 已定义关键词
var definedWords *set.Set

// sortedWords 按照权重排序的关键词
var sortedWords []string

// Word 单词
type Word struct {
	data     string
	wordType WordType
}

func newWord(data string, t WordType) (w *Word) {
	w = new(Word)
	w.data = data
	w.wordType = t
	return
}

const (
	// Assign 赋值
	Assign = "="
	// Import 导入
	Import = "import"
	// Export 导出
	Export = "export"
)

func init() {
	definedWords = set.NewSet()

	registerWords([]string{
		Export,
		Assign,
		Import,
	})
}

func registerWords(words []string) {
	for _, w := range words {
		registerWord(w)
	}
}

func registerWord(word string) {
	definedWords.Set(word)
	sortedWords = append(sortedWords, word)
}

func splitWords(sentence string) (words []*Word) {
	_words := strings.Split(sentence, " ")

	for _, w := range _words {
		t, w := getWordType(w)
		_w := newWord(w, t)
		words = append(words, _w)
	}

	words = rmEmptyWords(words)

	return
}

// getWordType 获取单词的类型，并获取真实值
func getWordType(w string) (t WordType, realValue string) {
	realValue = w

	if isDefined(w) {
		t = Operator
		return
	}

	if isConst(w) {
		t = Const
		// 去除首尾的引号
		realValue = realValue[1:]
		realValue = realValue[:len(realValue)-1]
		return
	}

	t = VarName
	return
}

func isDefined(w string) bool {
	return definedWords.Exist(w)
}

func isConst(w string) (is bool) {
	if len(w) < 2 {
		return
	}

	if w[0] == '"' && w[len(w)-1] == '"' {
		is = true
	}
	return
}

func rmEmptyWords(src []*Word) (dst []*Word) {
	for _, w := range src {
		if w.data == "" {
			continue
		}
		dst = append(dst, w)
	}
	return
}
