package profilemanager

import (
	"github.com/babolivier/matrix-notification-profile-manager/profilemanager/matrix"
	"github.com/babolivier/matrix-notification-profile-manager/types"

	"github.com/matrix-org/gomatrix"
)

func resolveDelta(srcMap, refMap map[string]types.PushRule, act action, delta *[]deltaRule) {
	for ruleID := range srcMap {
		if _, ok := refMap[ruleID]; !ok {
			*delta = append(*delta, deltaRule{
				action: act,
				rule:   srcMap[ruleID],
			})
		}
	}
}

func applyDelta(cli *gomatrix.Client, delta []deltaRule) (err error) {
	for _, deltaRule := range delta {
		if deltaRule.action == ACTION_ADD {
			if err = matrix.AddPushRule(cli, deltaRule.rule); err != nil {
				return err
			}
		} else if deltaRule.action == ACTION_DELETE {
			if err = matrix.DeletePushRule(cli, deltaRule.rule); err != nil {
				return err
			}
		}
	}

	return nil
}
