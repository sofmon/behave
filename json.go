package behave

// JSON entry point for all actions
var JSON = JSONAnchor{}

// JSONAnchor for actions
type JSONAnchor struct {
}

// ThenObjectMatches sample object
func (h *JSONAnchor) ThenObjectMatches(v interface{}) *JSONMatch {
	return (&JSONMatch{}).HavingMatchWith(v)
}
