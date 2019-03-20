package profilemanager

import (
	"github.com/babolivier/matrix-notification-profile-manager/types"
)

type action int

const (
	ACTION_ADD action = iota
	ACTION_DELETE
)

type deltaRule struct {
	rule   types.PushRule
	action action
}
