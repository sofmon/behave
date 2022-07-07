package main

import (
	"fmt"

	b "github.com/sofmon/behave"
)

// Given_we_can_access specific URL
func Given_we_can_access(url string) *CheckNetworkAccess {
	return (&CheckNetworkAccess{}).With_url(url).With_status(200)
}

// CheckNetworkAccess as an example of custom action
type CheckNetworkAccess struct {
	url    string
	status int
}

// With_url sets the url of the server to check
func (x *CheckNetworkAccess) With_url(url string) *CheckNetworkAccess {
	x.url = url
	return x
}

// With_status sets the expected status code
func (x *CheckNetworkAccess) With_status(status int) *CheckNetworkAccess {
	x.status = status
	return x
}

// String of the step definition
func (x *CheckNetworkAccess) String(res any) string {
	return fmt.Sprintf("Given we can access '%s' and receive status code of '%d'", x.url, x.status)
}

// Do the action
func (x *CheckNetworkAccess) Do(res any) any {
	return b.Do(
		b.When_we_make_http_call(x.url),
		b.Then_http_response_is(x.status),
	)
}
