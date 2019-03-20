package types

type PushRules struct {
	Global Ruleset `json:"global"`
}

type Ruleset struct {
	Content   []PushRule `json:"content"`
	Override  []PushRule `json:"override"`
	Room      []PushRule `json:"room"`
	Sender    []PushRule `json:"sender"`
	Underride []PushRule `json:"underride"`
}

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

type PushCondition struct {
	Kind    string `json:"kind"`
	Key     string `json:"key"`
	Pattern string `json:"pattern"`
	Is      string `json:"is"`
}

func (pr *PushRules) ToRulesMap() map[string]PushRule {
	internalRules := make(map[string]PushRule, 0)

	// TODO: this is ugly, surely we can make it look better
	for _, rule := range pr.Global.Content {
		rule.Scope = "global"
		rule.Kind = "content"
		internalRules[rule.RuleID] = rule
	}

	for _, rule := range pr.Global.Override {
		rule.Scope = "global"
		rule.Kind = "override"
		internalRules[rule.RuleID] = rule
	}

	for _, rule := range pr.Global.Room {
		rule.Scope = "global"
		rule.Kind = "room"
		internalRules[rule.RuleID] = rule
	}

	for _, rule := range pr.Global.Sender {
		rule.Scope = "global"
		rule.Kind = "sender"
		internalRules[rule.RuleID] = rule
	}

	for _, rule := range pr.Global.Underride {
		rule.Scope = "global"
		rule.Kind = "underride"
		internalRules[rule.RuleID] = rule
	}

	return internalRules
}
