package main

import (
	"time"

	"github.com/google/go-github/github"

	b "github.com/sofmon/behave"
)

func init() {

	rep := github.Repository{}

	b.Do(
		Given_we_can_access("https://google.com").With_status(200),
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
