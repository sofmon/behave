package behave

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"
)

// HTTPAction towards HTTP server
type HTTPAction struct {
	method     string
	url        string
	headers    map[string]string
	bodyString string
	bodyJSON   interface{}

	req *http.Request
}

// WithURL changes the HTTP method
func (x *HTTPAction) WithURL(url string) *HTTPAction {
	x.url = url
	return x
}

// WithMethod changes the HTTP method
func (x *HTTPAction) WithMethod(method string) *HTTPAction {
	x.method = method
	return x
}

// WithHeader added to the HTTP call
func (x *HTTPAction) WithHeader(key, val string) *HTTPAction {
	if x.headers == nil {
		x.headers = make(map[string]string)
	}
	x.headers[key] = val
	return x
}

// WithJSONBody changes the HTTP body to JSON object
func (x *HTTPAction) WithJSONBody(body interface{}) *HTTPAction {
	x.bodyJSON = body
	return x
}

// WithStringBody changes the HTTP body to JSON object
func (x *HTTPAction) WithStringBody(body string) *HTTPAction {
	x.bodyString = body
	return x
}

/* Action implementation */

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
