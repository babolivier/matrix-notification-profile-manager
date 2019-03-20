package profilemanager

import (
	"github.com/babolivier/matrix-notification-profile-manager/profilemanager/matrix"

	"github.com/matrix-org/gomatrix"
	log "github.com/sirupsen/logrus"
)

func SnapshotProfile(cli *gomatrix.Client, name string, overwrite bool) error {
	profile, err := matrix.GetPushRules(cli)
	if err != nil {
		return err
	}

	if err = matrix.SaveProfile(cli, name, profile, overwrite); err != nil {
		return err
	}

	log.Printf("Profile \"%s\" saved from current notification settings", name)
	return nil
}

func ApplyProfile(cli *gomatrix.Client, name string) error {
	profiles, err := matrix.GetProfiles(cli)
	if err != nil {
		return err
	}

	profileRules := profiles[name]
	profileRulesMap := profileRules.ToRulesMap()

	currentRules, err := matrix.GetPushRules(cli)
	if err != nil {
		return err
	}
	currentRulesMap := currentRules.ToRulesMap()

	delta := make([]deltaRule, 0)
	// If a rule in the profile isn't in the current rules, we must add it
	resolveDelta(profileRulesMap, currentRulesMap, ACTION_ADD, &delta)
	// If a rule in the current rules isn't in the profile, we must delete it
	resolveDelta(currentRulesMap, profileRulesMap, ACTION_DELETE, &delta)

	if err = applyDelta(cli, delta); err != nil {
		return err
	}

	log.Printf("Profile \"%s\" applied", name)

	return nil
}

func DeleteProfile(cli *gomatrix.Client, name string) error {
	if err := matrix.DeleteProfile(cli, name); err != nil {
		return err
	}

	log.Printf("Profile \"%s\" deleted", name)

	return nil
}
