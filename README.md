# Matrix notification profile manager

[![GoDoc](https://godoc.org/github.com/babolivier/matrix-notification-profile-manager/profilemanager?status.svg)](https://godoc.org/github.com/babolivier/matrix-notification-profile-manager/profilemanager)

[Matrix](https://matrix.org) offers amazing granularity when it comes to
configuring notification settings. However, this can result in a lot of
operations when setting the overall client settings to a desired configuration.

Matrix notification profile manager provides both a command line interface
and a Go package that allows one to manage notification settings profiles. It
works by snapshotting notification settings at a given time, saving that
snapshot as a profile, and switching from one profile to another by adding or
deleting the necessary rules so that the current user's settings match the
desired profile.

## Command line interface

The CLI can be retrieved and built by running:

```
go get github.com/babolivier/matrix-notification-profile-manager
```

Updating it is as simple as running:

```
go get -u github.com/babolivier/matrix-notification-profile-manager
```

Then you can use it in one of the available modes:

* Take a snapshot of the current notification settings and save it as a profile:

```
matrix-notification-profile-manager --snapshot --name PROFILE_NAME
```

If a profile already exists with the name `PROFILE_NAME`, you can overwrite it
by adding `--overwrite` to your command line.

* List the existing profiles:

```
matrix-notification-profile-manager --list
```

* Apply a profile:

```
matrix-notification-profile-manager --apply --name PROFILE_NAME
```

* Delete a profile:

```
matrix-notification-profile-manager --delete --name PROFILE_NAME
```

**Note**: Every mode support debug logging whenever the command line includes
`--debug`.

This interface needs a configuration file which path can be provided via the
`--config` flag (which defaults to `./config.yaml`). An example configuration
file is provided on this repository as
[`config.sample.yaml`](/config.sample.yaml).

## Go package

Most of the profile management logic lies in the `profilemanager` subpackage, so
using it can be done this way:

```go
package main

import (
	"github.com/babolivier/matrix-notification-profile-manager/profilemanager"

	"github.com/matrix-org/gomatrix"
)

func main() {
	cli, err := gomatrix.NewClient("https://matrix.example.com", "@alice:example.com", "ACCESS_TOKEN")
	if err != nil {
		panic(err)
	}

	// Save the current notification as a new profile.
	// If the last parameter (`overwrite`) is false, SnapshotSettings will
	// return an error if a profile already exists with the name `PROFILE_NAME`.
	if err = profilemanager.SnapshotSettings(cli, "PROFILE_NAME", true); err != nil {
		panic(err)
	}

	// Apply the profile named PROFILE_NAME to the user's notification settings.
	if err = profilemanager.ApplyProfile(cli, "PROFILE_NAME"); err != nil {
		panic(err)
	}

	// Delete the profile named PROFILE_NAME.
	if err = profilemanager.DeleteProfile(cli, "PROFILE_NAME"); err != nil {
		panic(err)
	}
}
```

The full documentation is available [here (on
godoc)](https://godoc.org/github.com/babolivier/matrix-notification-profile-manager/profilemanager).
More Matrix-specific functions that don't always need extra processing (e.g.
profile retrieval, etc.) are documented [here (on
godoc)](https://godoc.org/github.com/babolivier/matrix-notification-profile-manager/profilemanager/matrix).

## Profile storage

In order to make it interoperable between programs using the Go package, it
saves the profiles in the user's account data on the user's Matrix homeserver.
The type of the data is `bzh.abolivier.profiles.push`, and its structure is:

```json
{
	"default": {
		...
	},
	"profile1": {
		...
	}
}
```

Here, `default` and `profile1` are two profiles' names, and the data associated
follow the structure for the response format of [this Matrix
endpoint](https://matrix.org/docs/spec/client_server/r0.4.0.html#get-matrix-client-r0-pushrules).

## What's next

Here's a quick preview on the next improvements I'd like to make on this work:

* move the push rules management at the Matrix level to [gomatrix](https://github.com/matrix-org/gomatrix)
* implement renaming a profile
* ...
