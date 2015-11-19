package command

import (
	"github.com/cmceniry/cssa/stash"
	"fmt"
	"os"
)

func InitializeStash(opts *GeneralOptions, args []string) {
	// initialize an empty stash
	s, err := stash.CreateHollowStash(opts.CssaConfigfile)
	if err != nil {
		fmt.Printf("Error initializing stash: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("%#v\n", s)
}
