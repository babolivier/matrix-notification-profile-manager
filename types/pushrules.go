package types

// PushRules implements the response structure from
// https://matrix.org/speculator/spec/HEAD/client_server/unstable.html#get-matrix-client-r0-pushrules
type PushRules struct {
	Global Ruleset `json:"global"`
}

// Ruleset implements the Ruleset structure from
// https://matrix.org/speculator/spec/HEAD/client_server/unstable.html#get-matrix-client-r0-pushrules
type Ruleset struct {
	Content   []PushRule `json:"content"`
	Override  []PushRule `json:"override"`
	Room      []PushRule `json:"room"`
	Sender    []PushRule `json:"sender"`
	Underride []PushRule `json:"underride"`
}

// PushRule implements the PushRule structure from
// https://matrix.org/speculator/spec/HEAD/client_server/unstable.html#get-matrix-client-r0-pushrules
// It also has a Scope and a Kind field that are ignored by the JSON processing
// and helps with the addition and deletion of rules.
type PushRule struct {
	Actions    interface{}     `json:"actions"`
	Default    bool            `json:"default"`
	Enabled    bool            `json:"enabled"`
	RuleID     string          `json:"rule_id"`
	Conditions []PushCondition `json:"conditions"`
	Pattern    string          `json:"pattern"`
	Scope      string          `json:"-"`
	Kind       string          `json:"-"`
}

// PushCondition implements the PushCondition structure from
// https://matrix.org/speculator/spec/HEAD/client_server/unstable.html#get-matrix-client-r0-pushrules
type PushCondition struct {
	Kind    string `json:"kind"`
	Key     string `json:"key"`
	Pattern string `json:"pattern"`
	Is      string `json:"is"`
}

// ToRulesMap converts a PushRules instance into a map associating a rule ID to
// the matching rule (the rule having the custom Scope and Kind fields set).
// This makes a PushRules payload much easier to process.
func (pr *PushRules) ToRulesMap() map[string]PushRule {
	rulesMap := make(map[string]PushRule, 0)

	// TODO: Find a way to do this in a more efficient way.
	for _, rule := range pr.Global.Content {
		rule.Scope = "global"
		rule.Kind = "content"
		rulesMap[rule.RuleID] = rule
	}

	for _, rule := range pr.Global.Override {
		rule.Scope = "global"
		rule.Kind = "override"
		rulesMap[rule.RuleID] = rule
	}

	for _, rule := range pr.Global.Room {
		rule.Scope = "global"
		rule.Kind = "room"
		rulesMap[rule.RuleID] = rule
	}

	for _, rule := range pr.Global.Sender {
		rule.Scope = "global"
		rule.Kind = "sender"
		rulesMap[rule.RuleID] = rule
	}

	for _, rule := range pr.Global.Underride {
		rule.Scope = "global"
		rule.Kind = "underride"
		rulesMap[rule.RuleID] = rule
	}

	return rulesMap
}
