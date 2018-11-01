package main

import (
	"fmt"
	"time"

	b "github.com/sofmon/behave"
)

func main() {

	obj := struct {
		UserID int `json:"userId"`
	}{
		UserID: 1,
	}

	b.Do(
		b.HTTP.WhenWeCall("https://jsonplaceholder.typicode.com/todos/1").WithMethod("GET").WithHeader("Accepts", "application/json").WithJSONBody(obj),
		b.HTTP.ThenResponseIs(200).HavingHeader("Content-Type", "application/json; charset=utf-8"),
		b.JSON.ThenObjectMatches(&obj),
		GivenWeWaitFor(5*time.Second),
		WhenWeBrowseGoogle(),
	)
}

/* Example 1 */

// WaitAction - define a new custom action
type WaitAction struct {
	duration time.Duration
}

// GivenWeWaitFor - define a nice way to initialize the custom action
// you can do WaitAction{}, but sometimes it is better to initialize through function
func GivenWeWaitFor(d time.Duration) *WaitAction {
	return &WaitAction{d}
}

// String - implements the b.Action interface
func (x *WaitAction) String() string {
	return fmt.Sprintf("Given we wait for %v", x.duration)
}

// Do - implements the b.Action interface
func (x *WaitAction) Do(res b.Result) b.Result {
	time.Sleep(x.duration)
	return res // be nice and allow other to read the result
}

/* Example 2 */

// BrowseGoogleAction - define a new custom action
type BrowseGoogleAction struct{}

// WhenWeBrowseGoogle - define a nice way to initialize the custom action
func WhenWeBrowseGoogle() *BrowseGoogleAction {
	return &BrowseGoogleAction{}
}

// String - implements the b.Action interface
func (x *BrowseGoogleAction) String() string {
	return "When we browse google"
}

// Do - implements the b.Action interface
func (x *BrowseGoogleAction) Do(res b.Result) b.Result {
	return b.Do(
		b.HTTP.WhenWeCall("https://google.com").WithMethod("GET").WithHeader("Accepts", "text/html"),
		b.HTTP.ThenResponseIs(200).HavingHeader("Content-Type", "application/json; charset=utf-8"),
	)
}
