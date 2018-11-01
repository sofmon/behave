package behave

// HTTP entry point for all actions
var HTTP = HTTPAnchor{}

// HTTPAnchor for actions
type HTTPAnchor struct {
}

// WhenWeCall an HTTP service
func (h *HTTPAnchor) WhenWeCall(url string) *HTTPAction {
	return (&HTTPAction{}).WithURL(url)
}

// ThenResponseIs matching
func (h *HTTPAnchor) ThenResponseIs(status int) *HTTPActionCheck {
	return (&HTTPActionCheck{}).HavingStatus(status)
}
