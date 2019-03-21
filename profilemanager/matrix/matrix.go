// Package matrix implements Matrix-specific operations and logic for managing
// push notification settings profiles for Matrix users.
package matrix

import (
	"errors"
	"net/http"

	"github.com/babolivier/matrix-notification-profile-manager/profilemanager/types"

	"github.com/matrix-org/gomatrix"
	log "github.com/sirupsen/logrus"
)

// ACCOUNT_DATA_TYPE is the type of the data that gets added to the user's
// account data.
const ACCOUNT_DATA_TYPE = "bzh.abolivier.profiles.push"

var (
	// ErrProfileExists is returned when trying to create a profile with a name
	// that already has a profile attached to it.
	ErrProfileExists = errors.New("A profile already exists with this name, use another name or the -overwrite flag to overwrite it")
	// ErrProfileNotExists is returned when trying to delete a profile which
	// doesn't exist.
	ErrProfileNotExists = errors.New("No profile exists with this name")
)

// GetPushRules retrieves the notifications settings for the user.
// Returns an error if the retrieval failed.
func GetPushRules(cli *gomatrix.Client) (rules types.PushRules, err error) {
	url := cli.BuildURL("pushrules/")
	_, err = cli.MakeRequest("GET", url, nil, &rules)
	return
}

// AddPushRule adds a rule to the user's notification settings.
// Returns an error if the addition failed.
func AddPushRule(cli *gomatrix.Client, rule types.PushRule) error {
	url := cli.BuildURL("pushrules", rule.Scope, rule.Kind, rule.RuleID)
	_, err := cli.MakeRequest("PUT", url, rule, nil)
	return err
}

// DeletePushRule deletes a rule from the user's notification settings.
// Returns an error if the deletion failed.
func DeletePushRule(cli *gomatrix.Client, rule types.PushRule) error {
	url := cli.BuildURL("pushrules", rule.Scope, rule.Kind, rule.RuleID)
	_, err := cli.MakeRequest("DELETE", url, nil, nil)
	return err
}

// GetProfiles retrieves the list of profile for the user's account data.
// Returns an error if the retrieval failed.
func GetProfiles(cli *gomatrix.Client) (profiles map[string]types.PushRules, err error) {
	log.Debug("Retrieving user profiles")
	url := cli.BuildURL("user", cli.UserID, "account_data", ACCOUNT_DATA_TYPE)
	_, err = cli.MakeRequest("GET", url, nil, &profiles)
	log.Debugf("Got %d profiles", len(profiles))
	return
}

// SaveProfile saves a profile to the user's account data.
// Returns ErrProfileExists if a profile already exists for the provided name,
// unless overwrite is true, in which case it overwrites the profile the name is
// attached to.
// Returns an error if there was an issue talking with the Matrix homeserver.
func SaveProfile(cli *gomatrix.Client, name string, profile types.PushRules, overwrite bool) error {
	// Retrieve the list of profiles so we can edit it.
	profiles, err := GetProfiles(cli)
	if err != nil {
		// If the error results from the homeserver responding with 404 Not
		// Found, it means that the user doesn't have any profile stored, in
		// which case we just need to save the one we have.
		httpErr, ok := err.(gomatrix.HTTPError)
		if !ok || httpErr.Code != http.StatusNotFound {
			return err
		}

		// We didn't retrieve any profiles list from the homeserver, so we need
		// to initialise the map.
		log.Debug("Got 404 Not Found from the homeserver, continuing wth an empty map")
		profiles = make(map[string]types.PushRules)
	}

	// Check whether a profile already exists for this name and we can overwrite
	// it.
	if _, ok := profiles[name]; ok && !overwrite {
		return ErrProfileExists
	} else if ok {
		log.Debugf("Overwriting profile %s", name)
	}

	// Add the profile to the profiles list.
	profiles[name] = profile

	// Save the profiles list.
	return saveProfiles(cli, profiles)
}

// DeleteProfile deletes a profile from the user's account data.
// Returns ErrProfileNotExists if there's no profile with this name in the
// user's account data.
// Returns an error if there was an issue talking with the Matrix homeserver.
func DeleteProfile(cli *gomatrix.Client, name string) error {
	// Retrieve the list of profiles so we can edit it.
	// We don't check whether the error, if any, comes from the homeserver
	// responding with a 404 Not Found, because the right thing to do is to
	// return an error in this case. Though we could check for it and return
	// ErrProfileNotExists for a better UX (TODO?).
	profiles, err := GetProfiles(cli)
	if err != nil {
		return err
	}

	// Check if the profile exists.
	if _, ok := profiles[name]; !ok {
		log.Debugf("Trying to delete non-existing profile %s", name)
		return ErrProfileNotExists
	}

	// Delete the profile from the list.
	delete(profiles, name)

	// Save the updated list.
	return saveProfiles(cli, profiles)
}

// saveProfiles saves a profiles list to the user's account data.
// Returns an error if the operation failed.
func saveProfiles(cli *gomatrix.Client, profiles map[string]types.PushRules) error {
	log.Debugf("Saving the new profiles list containing %d profiles", len(profiles))
	url := cli.BuildURL("user", cli.UserID, "account_data", ACCOUNT_DATA_TYPE)
	_, err := cli.MakeRequest("PUT", url, profiles, nil)
	return err
}
