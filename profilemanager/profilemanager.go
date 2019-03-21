// Package profilemanager implements util functions for managing push
// notification settings profiles for Matrix users.
package profilemanager

import (
	"github.com/babolivier/matrix-notification-profile-manager/profilemanager/matrix"

	"github.com/matrix-org/gomatrix"
	log "github.com/sirupsen/logrus"
)

// EnableDebugLogging enable debug logs.
func EnableDebugLogging() {
	log.SetLevel(log.DebugLevel)
}

// SnapshotSettings takes a snapshot of the current notifications settings for
// the user and saves it as a profile in the user's account data.
// If a profile already exists with this name, returns matrix.ErrProfileExists
// unless overwrite is true, in which case it overwrites the existing profile
// for this name.
// Returns an error if either retrieving the settings or saving the profile
// failed.
func SnapshotSettings(cli *gomatrix.Client, name string, overwrite bool) error {
	log.Debugf("Snapshotting %s", name)

	// Retrieve current settings.
	profile, err := matrix.GetPushRules(cli)
	if err != nil {
		return err
	}

	log.Debugf("Got user's notification settings")

	// Save settings as a new profile.
	if err = matrix.SaveProfile(cli, name, profile, overwrite); err != nil {
		return err
	}

	log.Debugf("Saved profile %s", name)

	return nil
}

// ApplyProfile retrieves and apply a notification profile.
// Returns an error if retrieving the profile or applying it failed.
func ApplyProfile(cli *gomatrix.Client, name string) error {
	log.Debugf("Applying profile %s", name)

	// Retrieve the profiles.
	profiles, err := matrix.GetProfiles(cli)
	if err != nil {
		return err
	}
	log.Debugf("Retrieved %d profiles", len(profiles))

	profileRules := profiles[name]
	profileRulesMap := profileRules.ToRulesMap()

	// Retrieve the current rules to resolve the delta that needs to be applied.
	currentRules, err := matrix.GetPushRules(cli)
	if err != nil {
		return err
	}
	log.Debug("Retrieved current rules")
	currentRulesMap := currentRules.ToRulesMap()

	// Resolve the delta, i.e. determine which rules need to be added and which
	// rules need to be deleted in order for the notification settings to match
	// the profile.
	delta := make([]deltaRule, 0)
	// If a rule in the profile isn't in the current rules, we must add it
	resolveDelta(profileRulesMap, currentRulesMap, ACTION_ADD, &delta)
	// If a rule in the current rules isn't in the profile, we must delete it
	resolveDelta(currentRulesMap, profileRulesMap, ACTION_DELETE, &delta)

	log.Debugf("Got %d rules in delta")

	// Apply the delta (i.e. add and delete the necessary rules).
	if err = applyDelta(cli, delta); err != nil {
		return err
	}

	return nil
}

// DeleteProfile deletes a notification profile from the user's account data.
// Returns an error if either the deletion failed or no profile exists for this
// name (in which case it returns matrix.ErrProfileNotExists).
func DeleteProfile(cli *gomatrix.Client, name string) error {
	log.Debugf("Deleting profile %s", name)

	// Delete the profile.
	if err := matrix.DeleteProfile(cli, name); err != nil {
		return err
	}

	return nil
}
