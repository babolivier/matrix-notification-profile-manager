package main

type Mode int

const (
	NO_MODE Mode = iota
	MODE_SNAPSHOT
	MODE_SNAPSHOT_OVERWRITE
	MODE_LIST
	MODE_APPLY
	MODE_DELETE
)

func GetModeFromFlags() (mode Mode, err error) {
	mode = NO_MODE

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

	// The profile name is required in this condition.
	if mode != MODE_LIST {
		if len(*profilename) == 0 {
			err = errNoName
			return
		}
	}

	if mode == MODE_SNAPSHOT && *overwriteSnapshot {
		mode = MODE_SNAPSHOT_OVERWRITE
	}

	if mode == NO_MODE {
		err = errNoMode
		return
	}

	return
}

func setMode(flag bool, desiredMode, currentMode Mode) (Mode, error) {
	if !flag {
		return currentMode, nil
	}

	if currentMode != NO_MODE {
		return NO_MODE, errTooManyModes
	}

	return desiredMode, nil
}
