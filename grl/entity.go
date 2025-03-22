package grl

type GRuleEntity struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	When        string   `json:"when,omitempty"`
	Then        []string `json:"then,omitempty"`
	Salience    string   `json:"salience,omitempty"`
}
