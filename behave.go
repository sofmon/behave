// Package behave as part of project https://github.com/sofmon/behave
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package behave

import (
	"fmt"
	"os"
	"strings"
)

// Action to be performed, usually starts with
// Given... - we set specific state upfront
// When... - we perform an action
// Then.. - we check for specific result
type Action interface {
	String() string
	Do(Result) Result
}

// Result of an action
type Result interface {
	String() string
}

// JSONResult is a result that can be read as JSON object
type JSONResult interface {
	JSON() []byte
}

const (
	prefixSeperator = "    "
)

var (
	prefix = ""
)

// Do set of actions
func Do(acts ...Action) (res Result) {

	for i, act := range acts {

		if act == nil {
			continue
		}

		func() {

			defer func() {
				recErr := recover()
				if recErr != nil {
					msg := fmt.Sprintf("FAILED: %v", recErr)
					prefixLogf(msg)
					os.Exit(1)
				}
			}()

			increasePrefix()
			prefixLogf("(%d) %s", i+1, act.String())
			res = act.Do(res)
			decreasePrefix()

		}()
	}

	return
}

func increasePrefix() {
	prefix += prefixSeperator
}

func decreasePrefix() {
	prefix = prefix[:len(prefix)-len(prefixSeperator)]
}

func prefixLogf(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	msg = "\n" + prefix + strings.Replace(msg, "\n", "\n"+prefix, -1)
	fmt.Println(msg)
}
