package script

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrTooManyWords 单词太多
var ErrTooManyWords = errors.New("too many words")

// ErrInvalidLeftValue 非法的左值
var ErrInvalidLeftValue = errors.New("invalid left value")

// ErrInvalidRightValue 非法的右值
var ErrInvalidRightValue = errors.New("invalid right value")

// ErrInvalidOperator 非法的操作符
var ErrInvalidOperator = errors.New("invalid operator")

// astNode AST 节点
type astNode struct {
	value *Word
	left  *astNode
	right *astNode
}

func buildASTNode(wordsOfSentence []*Word) (n *astNode) {
	if len(wordsOfSentence) == 0 {
		return
	}

	value, index := findOperator(wordsOfSentence)
	if value == nil {
		// 没匹配到关键词
		n = buildValueNode(wordsOfSentence)
		return
	}

	var wordsLeft []*Word
	var wordsRight []*Word

	wordsLeft = wordsOfSentence[:index]
	wordsRight = wordsOfSentence[index+1:]

	n = new(astNode)
	n.value = value
	n.left = buildASTNode(wordsLeft)
	n.right = buildASTNode(wordsRight)

	return
}

// buildValueNode 构造值节点
func buildValueNode(wordsOfSentence []*Word) (n *astNode) {
	if len(wordsOfSentence) > 1 {
		err := errors.Wrapf(ErrTooManyWords, "%v", wordsOfSentence)
		panic(err)
	}

	n = new(astNode)
	n.value = wordsOfSentence[0]

	return
}

func findOperator(wordsOfSentence []*Word) (value *Word, index int) {
	for _, dw := range sortedWords {
		found, i := findOperatorWord(wordsOfSentence, dw)
		if found {
			value = wordsOfSentence[i]
			index = i
			break
		}
	}
	return
}

func findOperatorWord(words []*Word, dw string) (found bool, index int) {
	for i, theW := range words {
		if theW.wordType != Operator {
			continue
		}
		if theW.data == dw {
			found = true
			index = i
			break
		}
	}
	return
}

func (n *astNode) Run(vars map[string]string, exports Setter) (value *Word) {
	var valueLeft *Word
	var valueRight *Word

	if n.left != nil {
		valueLeft = n.left.Run(vars, exports)
	}
	if n.right != nil {
		valueRight = n.right.Run(vars, exports)
	}

	switch n.value.wordType {
	case Operator:
		value = excute(valueLeft, n.value.data, valueRight, vars, exports)
	default:
		value = n.value
	}

	return
}

func excute(left *Word, operator string, right *Word, vars map[string]string, exports Setter) (value *Word) {
	switch operator {
	case Import:
		_import(right, vars, exports)
		value = &(*right)
	case Export:
		export(right, vars, exports)
		value = &(*right)
	case Assign:
		assign(left, right, vars)
		// 赋值号的返回值为左值
		value = &(*left)
	default:
		err := errors.WithMessage(ErrInvalidOperator, operator)
		panic(err)
	}

	return
}

func _import(key *Word, vars map[string]string, exports Setter) (out string) {
	if key.wordType != VarName {
		err := errors.WithMessage(ErrInvalidLeftValue,
			fmt.Sprint(key.data, "-", key.wordType))
		panic(err)
	}

	realKey := key.data
	value := exports.Get(realKey)

	var _value string
	if value != nil {
		_value = fmt.Sprint(value)
	}

	vars[realKey] = _value
	out = _value
	return
}

func export(key *Word, vars map[string]string, exports Setter) (out string) {
	if key.wordType != VarName {
		err := errors.WithMessage(ErrInvalidLeftValue,
			fmt.Sprint(key.data, "-", key.wordType))
		panic(err)
	}

	value := vars[key.data]

	exports.Set(key.data, value)

	out = value
	return
}

func assign(key, value *Word, vars map[string]string) (out string) {
	if key.wordType != VarName {
		err := errors.WithMessage(ErrInvalidLeftValue,
			fmt.Sprint(key.data, "-", key.wordType))
		panic(err)
	}

	switch value.wordType {
	case VarName:
		vars[key.data] = vars[value.data]
	case Const:
		vars[key.data] = value.data
	default:
		err := errors.WithMessage(ErrInvalidRightValue,
			fmt.Sprint(value.data, "-", value.wordType))
		panic(err)
	}

	out = vars[key.data]
	return
}
