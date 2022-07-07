package behave

import (
	"bytes"
)

// Then_do simple do function
func Then_do(desc string, do func()) *DoSomething {
	return (&DoSomething{}).Also_do(desc, do)
}

// DoSomething action
type DoSomething struct {
	dos map[string]func()
}

// Also_do another
func (x *DoSomething) Also_do(desc string, do func()) *DoSomething {
	if x.dos == nil {
		x.dos = make(map[string]func())
	}
	x.dos[desc] = do
	return x
}

/* Action implementation */

func (x *DoSomething) String(res any) string {
	sb := bytes.NewBufferString("Then we do: \n")

	for k := range x.dos {
		sb.WriteString("  ")
		sb.WriteString(k)
		sb.WriteString("\n")
	}

	return sb.String()
}

// Do the action
func (x *DoSomething) Do(res any) any {

	for _, f := range x.dos {
		f()
	}

	return res
}
