package main

import (
	"flag"
	"os"

	"github.com/babolivier/matrix-notification-profile-manager/config"
	"github.com/babolivier/matrix-notification-profile-manager/profilemanager"
	"github.com/babolivier/matrix-notification-profile-manager/profilemanager/matrix"

	"github.com/matrix-org/gomatrix"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var mode Mode
var profilename, configPath *string

// TODO: Don't store these booleans as global variables.
var snapshot, list, apply, delete, overwriteSnapshot *bool

func init() {
	// Available modes.
	snapshot = flag.Bool("snapshot", false, "Take a snapshot of the current notification settings and save them as a profile")
	list = flag.Bool("list", false, "List the existing profiles. Doesn't need to be provided a profile name")
	apply = flag.Bool("apply", false, "Apply the provided profile as the current notification profile for the account")
	delete = flag.Bool("delete", false, "Delete the provided profile from the profiles database")

	// Additional settings.
	profilename = flag.String("name", "", "If used with --snapshot, name of the profile to create. Otherwise, name of the profile to apply or delete.")
	overwriteSnapshot = flag.Bool("overwrite", false, "Optional. If used with --snapshot and if a profile already exists with that name, overwrites it with the new snapshot. Useless otherwise.")

	// Config settings.
	configPath = flag.String("config", "config.yaml", "The path to the configuration file.")

	flag.Parse()

	// Guess the mode to run the profile manager in (i.e. the operation to do)
	// from the command line flags.
	var err error
	if mode, err = GetModeFromFlags(); err != nil {
		log.Error(errors.Wrap(err, "Error when parsing the commant-line arguments"))
		// TODO: Make flag.Usage nicer (with separate sections for modes and other settings).
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	// Load the configuration.
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Couldn't read configuration"))
	}

	// Initialise the Matrix client.
	cli, err := gomatrix.NewClient(cfg.Matrix.HSURL, cfg.Matrix.UserID, cfg.Matrix.AccessToken)
	if err != nil {
		log.Fatal(errors.Wrap(err, "init matrix client failed"))
	}

	// Act according to the mode.
	switch mode {
	case MODE_SNAPSHOT, MODE_SNAPSHOT_OVERWRITE:
		if err = profilemanager.SnapshotSettings(cli, *profilename, mode == MODE_SNAPSHOT_OVERWRITE); err != nil {
			log.Fatal(errors.Wrap(err, "Snapshot failed"))
		}
	case MODE_LIST:
		// Listing doesn't require any profile manipulation, so we do it from
		// the profilemanager/matrix package directly.
		profiles, err := matrix.GetProfiles(cli)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Retrieving profiles failed"))
		}

		log.Info("--- Profiles ---")

		for name := range profiles {
			log.Infof("\t- %s", name)
		}
	case MODE_APPLY:
		if err = profilemanager.ApplyProfile(cli, *profilename); err != nil {
			log.Fatal(errors.Wrap(err, "Applying profile failed"))
		}
	case MODE_DELETE:
		if err = profilemanager.DeleteProfile(cli, *profilename); err != nil {
			log.Fatal(errors.Wrap(err, "Deleting profile failed"))
		}
	}
}
