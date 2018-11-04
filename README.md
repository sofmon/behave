# Go behave

`behave` is a package helping to define, lock and check microservice behavier. It is inspired from `gherkin` but it takes advantage of [Go](https://golang.org)'s strong type syntaxis, compile time verification, package initialization and [dep](https://github.com/golang/dep) dependency manager.


# Base features

## Current features

### Make and check HTTP calls
- `When_we_make_http_call` - Creates an HTTP call for specific url. Default method is 'GET'.
  - `With_url` - Changes the url of the web call.
  - `With_method` - Defines the HTTP method.
  - `With_header` - Sets an HTTP method. Can be called multiple times for different headers.
  - `With_json_body` - Sets the body of the HTTP call as a JSON object, strongly defined as Go struct.
  - `With_string_body` - Sets the body of the HTTP call to a string. This have priority over  `With_json_body`.
- `Then_http_response_is` - Checks that an HTTP response have specific status code.
  - `Having_status` - Changes the expected HTTP response status code.
  - `Having_header` - Checks that HTTP response have a specific header with specific value.

HTTP result passed on can be used as JSON result (see below).

### Check JSON objects
- `Then_result_is_json` - Checks the result from the privies operation is a valid JSON string/object.
  - `Having_match_with` - Validates that the JSON object matches another JSON object provided as a Go struct. Only field within the matching object are checked.
  - `Also_extracted_to` - Extracts the JSON object into a provided Go struct.

### Custom checks
- `Then_check_that` - Performs a custom check written as a `func()bool`.
  - `Also_that` - Add another custom check to the behave step.

## Example
``` Go
package main

import (
	"time"

	"github.com/google/go-github/github"
	
	b "github.com/sofmon/behave"
)

func init() {

	rep := github.Repository{}

	b.Do(
		b.When_we_make_http_call("https://api.github.com/repos/sofmon/behave").
			With_method("GET").
			With_header("Accepts", "application/json"),
		b.Then_http_response_is(200).
			Having_header("Content-Type", "application/json; charset=utf-8"),
		b.Then_result_is_json().
			Having_match_with(
				&struct {
					URL string `json:"html_url"`
				}{
					URL: "https://github.com/sofmon/behave",
				},
			).
			Also_extracted_to(&rep),
		b.Then_check_that(
			"repository was created before 2018-11-04",
			func() bool {
				return rep.CreatedAt.Before(time.Date(2018, 11, 4, 0, 0, 0, 0, time.UTC))
			},
		).
			Also_that(
				"owner is sofmon",
				func() bool {
					return rep.Owner.Login != nil && *rep.Owner.Login == "sofmon"
				},
			),
	)
}

func main() {}
```

Output
```
$ go run .

    (1) When we do HTTP call like:
      GET /repos/sofmon/behave HTTP/1.1
      Host: api.github.com
      User-Agent: Go-http-client/1.1
      Accepts: application/json
      Accept-Encoding: gzip



    (2) Then HTTP response status code must be '200', having header 'Content-Type=application/json; charset=utf-8'

    (3) Then JSON object is having match with object like:
      {"html_url":"https://github.com/sofmon/behave"}

    (4) Then we check that:
      owner is sofmon
      repository was created before 2018-11-04
```

# Custom feature

## Define custom step/feature

``` Go
// CheckNetworkAccess as an example of custom action
type CheckNetworkAccess struct {
	url    string
	status int
}

/* entry point */

// Given_we_can_access as example of custom action
func Given_we_can_access(url string) *CheckNetworkAccess {
	return (&CheckNetworkAccess{}).With_url(url).With_status(200)
}

/* define modifications */

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

/* implementing behave action */

// String of the step definition
func (x *CheckNetworkAccess) String() string {
	return fmt.Sprintf("Given we can access '%s' and receive status code of '%d'", x.url, x.status)
}

// Do the action
func (x *CheckNetworkAccess) Do(res b.Result) b.Result {
	return b.Do(
		b.When_we_make_http_call(x.url),
		b.Then_http_response_is(x.status),
	)
}
```

## Usage
``` Go
b.Do(
  Given_we_can_access("https://google.com").With_status(200),
		...
)
```

Output:

```
$ go run .

    (1) Given we can access 'https://google.com' and receive status code of '200'

        (1) When we do HTTP call like:
          GET / HTTP/1.1
          Host: google.com
          User-Agent: Go-http-client/1.1
          Accept-Encoding: gzip



        (2) Then HTTP response status code must be '200'

    ...
```

# Environment behavier lock

At the moment of your commit, the service code expect specific behavier from the surrounding services, the environment. For example if your `vouchers` service, deposit money, you expect the `account` and `profile` services to behave in specific way. On the other hand the `account` service itself expects specific `securities` service behaver.

```
  ┏━━━━━━━━━━┓   ┏━━━━━━━━━┓
  ┃ vouchers ┃───┃ profile ┃
  ┗━━━━━━━━━━┛   ┗━━━━━━━━━┛
     │
     │  ┏━━━━━━━━━┓   ┏━━━━━━━━━━━━┓
     └──┃ account ┃───┃ securities ┃
        ┗━━━━━━━━━┛   ┗━━━━━━━━━━━━┛
```

The combined behavior of those services defines the environment in which the `vouchers` service is guaranteed, or at least tested to perform correctly. 

This is where [dep](https://github.com/golang/dep) helps. When you write all your tests. You will put them in dedicated package for that service. Then all test execution will be define in the `init()` function. This way at the moment a behave test uses the bahave features form another package, it will trigger the behaver tests for that dependent service.

Doing "`$ deb init`" or "`$ dep ensure --update`" you would download all related services behaver and lock them in the vendor folder.

Now the expected behavier of the environment is locked within the specific commit.
