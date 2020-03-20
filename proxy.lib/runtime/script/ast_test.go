package script

import (
	"log"
	"testing"
)

func Test_buildAST(t *testing.T) {
	sentence := `import a`
	ast := buildAST(sentence)
	log.Printf("%+v", ast.root)
}
