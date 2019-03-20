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

var (
	errNoName       = errors.New("The --name parameter is required. See the help (--help) for more information")
	errNoMode       = errors.New("You need to supply a mode (either --snapshot, --list, --apply or--delete)")
	errTooManyModes = errors.New("Only one mode (either --snapshot, --list, --apply or --delete) is allowed at a time")
)

var mode Mode
var profilename, configPath *string

func init() {
	snapshot = flag.Bool("snapshot", false, "Take a snapshot of the current notification settings and save them as a profile")
	list = flag.Bool("list", false, "List the existing profiles")
	apply = flag.Bool("apply", false, "Apply the provided profile as the current notification profile for the account")
	delete = flag.Bool("delete", false, "Delete the provided profile from the profiles database")

	profilename = flag.String("name", "", "Required. If used with --snapshot, name of the profile to create. Otherwise, name of the profile to apply or delete.")
	overwriteSnapshot = flag.Bool("overwrite", false, "Optional. If used with --snapshot and if a profile already exists with that name, overwrites it with the new snapshot. Useless otherwise.")

	configPath = flag.String("config", "config.yaml", "The path to the configuration file.")

	flag.Parse()

	var err error
	if mode, err = GetModeFromFlags(); err != nil {
		log.Error(errors.Wrap(err, "Error when parsing the commant-line arguments"))
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Couldn't read configuration"))
	}

	cli, err := gomatrix.NewClient(cfg.Matrix.HSURL, cfg.Matrix.UserID, cfg.Matrix.AccessToken)
	if err != nil {
		log.Fatal(errors.Wrap(err, "init matrix client failed"))
	}

	if *snapshot {

	}

	switch mode {
	case MODE_SNAPSHOT, MODE_SNAPSHOT_OVERWRITE:
		if err = profilemanager.SnapshotProfile(cli, *profilename, mode == MODE_SNAPSHOT_OVERWRITE); err != nil {
			log.Fatal(errors.Wrap(err, "Snapshot failed"))
		}
	case MODE_LIST:
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
