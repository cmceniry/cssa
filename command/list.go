package command

import (
	"sort"
	"fmt"
	"os"
)

func List(opts *GeneralOptions, args []string) {
	opts.Stash.Load()

	archivenames := opts.Stash.GetArchivelist()
	sort.Strings(archivenames)
	for _, a := range archivenames {
		fmt.Printf("%s\n", a)
	}
	os.Exit(0)
}
