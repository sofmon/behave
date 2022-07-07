package behave

import (
	"bytes"
	"fmt"
)

// Then_check_that simple check function
func Then_check_that(desc string, check func() bool) *Check {
	return (&Check{}).Also_that(desc, check)
}

// Check action
type Check struct {
	checks map[string]func() bool
}

// Also_that do another check
func (x *Check) Also_that(desc string, check func() bool) *Check {
	if x.checks == nil {
		x.checks = make(map[string]func() bool)
	}
	x.checks[desc] = check
	return x
}

/* Action implementation */

func (x *Check) String(res any) string {
	sb := bytes.NewBufferString("Then we check that: \n")

	for k := range x.checks {
		sb.WriteString("  ")
		sb.WriteString(k)
		sb.WriteString("\n")
	}

	return sb.String()
}

// Do the action
func (x *Check) Do(res any) any {

	for k, f := range x.checks {
		if !f() {
			panic(fmt.Errorf("check for '%s' failed", k))
		}
	}

	return res
}
