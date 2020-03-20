package script

import "testing"

func Test_Format(t *testing.T) {
	script := `
	import a
	a = 1
	export a
	`

	script = format(script)
	t.Log(script)
}
