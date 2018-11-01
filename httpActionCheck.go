package behave

import (
	"bytes"
	"errors"
	"fmt"
)

// HTTPActionCheck is matched to
type HTTPActionCheck struct {
	statusCode int
	headers    map[string]string
}

// HavingStatus for the http call
func (x *HTTPActionCheck) HavingStatus(status int) *HTTPActionCheck {
	x.statusCode = status
	return x
}

// HavingHeader for the http call
func (x *HTTPActionCheck) HavingHeader(key, value string) *HTTPActionCheck {
	if x.headers == nil {
		x.headers = make(map[string]string)
	}
	x.headers[key] = value
	return x
}

/* Action implementation */

func (x *HTTPActionCheck) String() string {
	sb := bytes.NewBufferString(fmt.Sprintf("Then HTTP response status code must be '%d'", x.statusCode))

	// headers
	for k, v := range x.headers {
		sb.WriteString(fmt.Sprintf(", having header '%s=%s'", k, v))
	}

	return sb.String()
}

// Do the action
func (x *HTTPActionCheck) Do(res Result) Result {

	resp, ok := res.(*httpResult)

	if !ok || resp == nil {
		panic(errors.New("privies operation did not produce HTTP response; please use 'HTTP.When_we_call(...) as prior action"))
	}

	// statusCode
	if resp.statusCode != x.statusCode {
		panic(fmt.Errorf("expecting status code of %d but received %d", resp.statusCode, x.statusCode))
	}

	// headers
	for k, v := range x.headers {
		vv := resp.header.Get(k)
		if vv != v {
			panic(fmt.Errorf("expecting header `%s` to have value of '%s' but it has `%s`", k, v, vv))
		}
	}

	return res
}
