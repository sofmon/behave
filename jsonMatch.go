package behave

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// JSONMatch action
type JSONMatch struct {
	v interface{}
}

// HavingMatchWith specific object
func (x *JSONMatch) HavingMatchWith(v interface{}) *JSONMatch {
	x.v = v
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

	err = json.Unmarshal(jsonObj.JSON(), &x.v)
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

	return res
}
