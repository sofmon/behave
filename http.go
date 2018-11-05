// Package behave as part of project https://github.com/sofmon/behave
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package behave

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

// When_we_make_http_call an HTTP service
func When_we_make_http_call(url string) *HTTPAction {
	return (&HTTPAction{}).With_url(url)
}

// Then_http_response_is matching
func Then_http_response_is(status int) *HTTPActionCheck {
	return (&HTTPActionCheck{}).Having_status(status)
}

// HTTPAction towards HTTP server
type HTTPAction struct {
	method     string
	url        string
	headers    map[string]string
	bodyString string
	bodyJSON   interface{}

	req *http.Request
}

// With_url changes the HTTP method
func (x *HTTPAction) With_url(url string) *HTTPAction {
	x.url = url
	return x
}

// With_method changes the HTTP method
func (x *HTTPAction) With_method(method string) *HTTPAction {
	x.method = method
	return x
}

// With_header added to the HTTP call
func (x *HTTPAction) With_header(key, val string) *HTTPAction {
	if x.headers == nil {
		x.headers = make(map[string]string)
	}
	x.headers[key] = val
	return x
}

// With_json_body changes the HTTP body to JSON object
func (x *HTTPAction) With_json_body(body interface{}) *HTTPAction {
	x.bodyJSON = body
	return x
}

// With_string_body changes the HTTP body to JSON object
func (x *HTTPAction) With_string_body(body string) *HTTPAction {
	x.bodyString = body
	return x
}

func (x *HTTPAction) ensureRequest() {
	if x.req != nil {
		return
	}
	body := []byte{}

	switch {
	case x.bodyJSON != nil:
		var err error
		body, err = json.Marshal(x.bodyJSON)
		if err != nil {
			panic(err)
		}

	case x.bodyString != "":
		body = []byte(x.bodyString)
	}

	var err error
	x.req, err = http.NewRequest(x.method, x.url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	for k, v := range x.headers {
		x.req.Header.Set(k, v)
	}
}

func (x *HTTPAction) String() string {

	x.ensureRequest()

	data, err := httputil.DumpRequestOut(x.req, true)
	if err != nil {
		panic(err)
	}

	sb := bytes.NewBufferString("When we do HTTP call like:\n")
	sb.WriteString("  ")
	sb.WriteString(strings.Replace(string(data), "\n", "\n  ", -1))

	return sb.String()
}

// Do the action
func (x *HTTPAction) Do(res Result) Result {

	x.ensureRequest()

	resp, err := http.DefaultClient.Do(x.req)
	if err != nil {
		panic(err)
	}

	return newHTTPResult(resp)
}

// HTTPActionCheck is matched to
type HTTPActionCheck struct {
	statusCode int
	headers    map[string]string
}

// Having_status for the http call
func (x *HTTPActionCheck) Having_status(status int) *HTTPActionCheck {
	x.statusCode = status
	return x
}

// HavingHeader for the http call
func (x *HTTPActionCheck) Having_header(key, value string) *HTTPActionCheck {
	if x.headers == nil {
		x.headers = make(map[string]string)
	}
	x.headers[key] = value
	return x
}

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
		panic(fmt.Errorf("expecting status code of %d but received %d", x.statusCode, resp.statusCode))
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

func newHTTPResult(resp *http.Response) (res *httpResult) {
	if resp == nil {
		return
	}

	res = &httpResult{}
	res.statusCode = resp.StatusCode
	res.header = resp.Header

	var err error
	res.dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	res.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return
}

type httpResult struct {
	statusCode int
	header     http.Header
	body       []byte
	dump       []byte
}

func (x *httpResult) String() string {
	if x == nil {
		return "nil"
	}
	return string(x.dump)
}

func (x *httpResult) JSON() []byte {
	if x == nil {
		return []byte("null")
	}

	return x.body
}
