package main

import (
	"errors"
)

// Mode is the mode the application runs with.
type Mode int

const (
	// NO_MODE is the initial zero-value for the mode.
	NO_MODE Mode = iota
	// MODE_SNAPSHOT is the mode for snapshot, without allowing for overwriting
	// in case of conflict.
	MODE_SNAPSHOT
	// MODE_SNAPSHOT_OVERWRITE is the mode for snapshot, allowing for
	// overwriting on conflicts.
	MODE_SNAPSHOT_OVERWRITE
	// MODE_LIST is the mode for listing profiles.
	MODE_LIST
	// MODE_APPLY is the mode for applying a profile.
	MODE_APPLY
	// MODE_DELETE is the mode for deleting a profile.
	MODE_DELETE
)

var (
	errNoName       = errors.New("The --name parameter is required in this mode. See the help (--help) for more information")
	errNoMode       = errors.New("You need to supply a mode (either --snapshot, --list, --apply or--delete)")
	errTooManyModes = errors.New("Only one mode (either --snapshot, --list, --apply or --delete) is allowed at a time")
)

// GetModeFromFlags guesses the mode to run the profile manager in (i.e. the
// operation to do) from the command line flags.
// Returns errNoName if no profile name was supplied in a mode in which it's
// required, errNoMode if no mode was supplied, or errTooManyModes if more than
// one mode was supplied.
func GetModeFromFlags() (mode Mode, err error) {
	mode = NO_MODE

	// Check for each mode if we need to apply it and if we can do it if needed.
	if mode, err = setMode(*snapshot, MODE_SNAPSHOT, mode); err != nil {
		return
	}
	if mode, err = setMode(*list, MODE_LIST, mode); err != nil {
		return
	}
	if mode, err = setMode(*apply, MODE_APPLY, mode); err != nil {
		return
	}
	if mode, err = setMode(*delete, MODE_DELETE, mode); err != nil {
		return
	}

	// Check if a mode was set.
	if mode == NO_MODE {
		err = errNoMode
		return
	}

	// Check if we the profile name is required.
	if mode != MODE_LIST {
		if len(*profilename) == 0 {
			err = errNoName
			return
		}
	}

	// Check if we can overwrite in snapshot mode.
	if mode == MODE_SNAPSHOT && *overwriteSnapshot {
		mode = MODE_SNAPSHOT_OVERWRITE
	}

	return
}

// setMode checks whether the current mode needs to be changed according to a
// command line flag, and if so checks whether we can do it (i.e. if the current
// mode is NO_MODE).
// Returns errTooManyModes if a mode is already set..
func setMode(flag bool, desiredMode, currentMode Mode) (Mode, error) {
	// Check if we need to care about this mode.
	if !flag {
		return currentMode, nil
	}

	// Check if we can change the mode.
	if currentMode != NO_MODE {
		return NO_MODE, errTooManyModes
	}

	return desiredMode, nil
}
