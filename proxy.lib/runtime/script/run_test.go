package script

import "testing"

func Test_Run(t *testing.T) {
	script := `
	import a
	a = "1"
	export a
	export b = import a
	`

	s := NewBaseSetter()
	r := NewRunner(s)

	r.Run(script)
}
