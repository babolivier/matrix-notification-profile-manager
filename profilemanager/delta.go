package profilemanager

import (
	"github.com/babolivier/matrix-notification-profile-manager/profilemanager/matrix"
	"github.com/babolivier/matrix-notification-profile-manager/types"

	"github.com/matrix-org/gomatrix"
)

// resolveDelta checks which rules from srcMap don't appear in refMap and adds
// them to the overall delta with the necessary action to take (i.e. add or
// delete).
func resolveDelta(srcMap, refMap map[string]types.PushRule, act Action, delta *[]deltaRule) {
	// Iterate over the source map.
	for ruleID := range srcMap {
		// Check if the rule is missing from the reference map.
		if _, ok := refMap[ruleID]; !ok {
			// If the rule is missing from the reference map, add it to the delta.
			*delta = append(*delta, deltaRule{
				action: act,
				rule:   srcMap[ruleID],
			})
		}
	}
}

// applyDelta adds and deletes the rules from the overall delta according to the
// action that accompanies each rule.
// Returns an error if either an addition or a deletion failed. One error stops
// the whole process.
func applyDelta(cli *gomatrix.Client, delta []deltaRule) (err error) {
	// Iterate over the delta
	for _, deltaRule := range delta {
		// Act according to the action.
		if deltaRule.action == ACTION_ADD {
			// Add the rule to the user's settings.
			if err = matrix.AddPushRule(cli, deltaRule.rule); err != nil {
				return err
			}
		} else if deltaRule.action == ACTION_DELETE {
			// Delete the rule from the user's settings.
			if err = matrix.DeletePushRule(cli, deltaRule.rule); err != nil {
				return err
			}
		}
	}

	return nil
}
