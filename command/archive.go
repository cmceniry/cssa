package command

import (
	"github.com/cmceniry/cssa/util"
	"github.com/cmceniry/cssa/stash"
	"time"
	"flag"
	"fmt"
	"os"
)

func Archive(opts *GeneralOptions, args []string) {
	archiveFlags := flag.NewFlagSet("archive", flag.ExitOnError)
	archiveFlags.Parse(args)
	if len(archiveFlags.Args()) < 1 {
		fmt.Printf("No snapshot name identified\n")
		os.Exit(-2)
	}
	snapshot := archiveFlags.Arg(0)

	ss := opts.Stash

	// archive flow
	//   Gather requested snapshot info
	//   Compare to inventory of snapshots
	//   Copy data from local to archive
	//   Update inventory

	m, warns, err := util.GetLocalSnapShot(opts.CstarConfigfile, snapshot)
	switch err {
	case nil:
	case util.SnapshotNotFound:
		fmt.Printf("Snapshot not found\n")
		os.Exit(-1)
	default:
		panic(err)
	}
	if len(warns) != 0 {
		fmt.Printf("Issues processing:\n")
		for _, w := range warns {
			fmt.Printf("  %s\n", w)
		}
		os.Exit(-1)
	}

	a := ss.CreateArchive(fmt.Sprintf("%s-%d", snapshot, time.Now().Unix()))
	for _, f := range m.Files {
		if err := a.AddFile(f); err != nil {
			switch err {
			case stash.ConflictFileExists:
				fmt.Printf("Conflicting archive file %s\n")
				a.SetInvalid()
			case stash.DuplicateFileExists:
				fmt.Printf("Speed up: %s\n", f.Filename)
			default:
				fmt.Printf("Unable to add file %s: %s\n", f.Filename, err)
				a.SetInvalid()
			}
			continue
		}
		fmt.Printf("Added file %s\n", f.Filename)
	}
	err = a.Commit()
	if err != nil {
		fmt.Printf("Error committing archive: %s\n", err)
		os.Exit(-1)
	}

	os.Exit(0)
}
