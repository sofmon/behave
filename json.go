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
	v interface{}
	e interface{}
}

// Having_match_with specific object
func (x *JSONMatch) Having_match_with(v interface{}) *JSONMatch {
	x.v = v
	return x
}

// Also_extracted_to specific object
func (x *JSONMatch) Also_extracted_to(v interface{}) *JSONMatch {
	x.e = v
	return x
}

/* Action implementation */

func (x *JSONMatch) String() string {
	sb := bytes.NewBufferString("Then JSON object is having match with object like:\n")

	bytes, err := json.Marshal(x.v)
	if err != nil {
		panic(err)
	}

	sb.WriteString("  ")
	sb.WriteString(strings.Replace(string(bytes), "\n", "\n  ", -1))

	return sb.String()
}

// Do the action
func (x *JSONMatch) Do(res Result) Result {

	jsonObj, ok := res.(JSONResult)

	if !ok || jsonObj == nil {
		panic(errors.New("privies operation did not produce object that can provide json"))
	}

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

	if x.e != nil {
		err = json.Unmarshal(jsonObj.JSON(), x.e)
		if err != nil {
			panic(err)
		}
	}

	return res
}
