package script

// AST 抽象语法树
type AST struct {
	root *astNode
}

func buildAST(sentence string) (ast *AST) {
	words := splitWords(sentence)

	ast = new(AST)
	ast.root = buildASTNode(words)

	return
}

// Run 执行语法树
func (ast *AST) Run(vars map[string]string, exports Setter) {
	if ast.root == nil {
		return
	}

	ast.root.Run(vars, exports)
	return
}
