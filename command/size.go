package command

import (
	"fmt"
	"os"
)

func Size(opts *GeneralOptions, args []string) {
	opts.Stash.Load()

	if len(args) < 1 {
		fmt.Printf("No archive name identified\n")
		os.Exit(-2)
	}
	archivename := args[0]

	s := opts.Stash

	a := s.GetArchive(archivename)
	if a == nil {
		fmt.Printf("No archive found by name %s\n", archivename)
		os.Exit(-2)
	}

	for _, filename := range a.GetFilelist() {
		fmt.Printf("%-20d\t%s\n", s.GetFileSize(filename), filename)
	}
}
