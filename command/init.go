package command

import (
	"fmt"
	"os"
)

func InitializeStash(opts *GeneralOptions, args []string) {
	// initialize an empty stash
	err := opts.Stash.CreateNew()
	if err != nil {
		fmt.Printf("Error initializing stash: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("%#v\n", opts.Stash)
}
