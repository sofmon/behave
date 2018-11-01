package behave

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

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
