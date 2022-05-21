// Package behave as part of project https://github.com/sofmon/behave
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package behave

import (
	"fmt"
	"strings"
)

// Action to be performed, usually starts with
// Given... - we set specific state upfront
// When... - we perform an action
// Then.. - we check for specific result
type Action interface {
	String() string
	Do(interface{}) interface{}
}

// JSONResult is a result that can be read as JSON object
type JSONResult interface {
	JSON() []byte
}

const (
	prefixSeparator = "  "
)

var (
	prefix = ""
)

// Do set of actions and indicate success
func Do(acts ...Action) (ok bool) {

	ok = true

	var res interface{}
	for i, act := range acts {

		if act == nil {
			continue
		}

		func() {

			defer func() {
				recErr := recover()
				if recErr != nil {
					prefixLogf(fmt.Sprintf("FAILED: %v", recErr))
					ok = false
				}
			}()

			increasePrefix()
			prefixLogf("(%d) %s", i+1, act.String())
			res = act.Do(res)
			decreasePrefix()

		}()

		if !ok {
			return false
		}

	}

	return ok
}

func increasePrefix() {
	prefix += prefixSeparator
}

func decreasePrefix() {
	prefix = prefix[:len(prefix)-len(prefixSeparator)]
}

func prefixLogf(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	msg = "\n" + prefix + strings.Replace(msg, "\n", "\n"+prefix, -1)
	fmt.Println(msg)
}
