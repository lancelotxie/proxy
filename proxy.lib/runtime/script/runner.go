package script

// Runner 脚本执行器
type Runner interface {
	Run(script string)
}

// NewRunner 构造新的 Runner
func NewRunner(setter Setter) (r Runner) {
	r = newRunner(setter)
	return
}

type runner struct {
	setter Setter
	vars   map[string]string
}

func newRunner(setter Setter) (r *runner) {
	r = new(runner)
	r.setter = setter
	r.vars = make(map[string]string)
	return
}

func (r *runner) Run(script string) {
	r.runScript(script)
}

func (r *runner) runScript(script string) {
	script = format(script)

	sentences := splitSentences(script)

	for _, sentence := range sentences {
		r.runSentence(sentence)
	}
}

func (r *runner) runSentence(sentence string) {
	ast := buildAST(sentence)
	ast.Run(r.vars, r.setter)
}
