package profilemanager

import (
	"github.com/babolivier/matrix-notification-profile-manager/profilemanager/types"
)

// Action is the action to take for a given rule.
type Action int

const (
	// ACTION_ADD accompanies a rule that needs to be added to the user's
	// settings.
	ACTION_ADD Action = iota
	// ACTION_DELETE accompanies a rule that needs to be deleted from the user's
	// settings.
	ACTION_DELETE
)

// deltaRule is how a rule is stored in the overall delta.
type deltaRule struct {
	rule   types.PushRule
	action Action
}
