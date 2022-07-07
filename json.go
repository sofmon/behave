// Package behave as part of project https://github.com/sofmon/behave
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package behave

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Then_result_is_json sample object
func Then_result_is_json() *JSONMatch {
	return &JSONMatch{}
}

// JSONMatch action
type JSONMatch struct {
	v  any
	vf func() any
	e  any
}

// Having_match_with specific object
func (x *JSONMatch) Having_match_with(v any) *JSONMatch {
	x.v = v
	return x
}

// Having_match_with specific object
func (x *JSONMatch) Having_match_with_func(f func() any) *JSONMatch {
	x.vf = f
	return x
}

// Also_extracted_to specific object
func (x *JSONMatch) Also_extracted_to(v any) *JSONMatch {
	x.e = v
	return x
}

/* Action implementation */

func (x *JSONMatch) String(res any) string {
	sb := bytes.NewBufferString("Then result is JSON object")

	jsonObj, ok := res.(JSONResult)

	if !ok || jsonObj == nil {
		panic(errors.New("privies operation did not produce object that can provide json"))
	}

	if x.e != nil {
		err := json.Unmarshal(jsonObj.JSON(), x.e)
		if err != nil {
			panic(err)
		}
	}

	if x.vf != nil {
		x.v = x.vf()
	}

	if x.v != nil {
		bytes, err := json.Marshal(x.v)
		if err != nil {
			panic(err)
		}
		sb.WriteString(", having match with object like:\n")
		sb.WriteString("  ")
		sb.WriteString(strings.Replace(string(bytes), "\n", "\n  ", -1))
	}

	return sb.String()
}

// Do the action
func (x *JSONMatch) Do(res any) any {

	jsonObj, ok := res.(JSONResult)

	if !ok || jsonObj == nil {
		panic(errors.New("privies operation did not produce object that can provide json"))
	}

	if x.e != nil {
		err := json.Unmarshal(jsonObj.JSON(), x.e)
		if err != nil {
			panic(err)
		}
	}

	if x.vf != nil {
		x.v = x.vf()
	}

	if x.v != nil {
		expData, err := json.Marshal(x.v)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(jsonObj.JSON(), x.v)
		if err != nil {
			panic(err)
		}

		resData, err := json.Marshal(x.v)
		if err != nil {
			panic(err)
		}

		if string(expData) != string(resData) {
			panic(fmt.Errorf("expected json object does not match expectations; received object: `%s`", string(resData)))
		}
	}

	return res
}
