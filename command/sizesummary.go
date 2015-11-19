package command

import (
	"fmt"
	"os"
)

func SizeSummary(opts *GeneralOptions, args []string) {
	if len(args) < 1 {
		fmt.Printf("No archive name(s) identified\n")
		os.Exit(-2)
	}

	s := opts.Stash

	// seen = map[filename] list of archivenames
	seen := map[string][]string{}
	// size = map[archivename] of slice (size of #archives) of int64
	total := map[string][]int64{}

	// First map all files back to the archives they come from
	for _, archivename := range args {
		a := s.GetArchive(archivename)
		if a == nil {
			fmt.Printf("No archive found by name %s\n", archivename)
			os.Exit(-2)
		}
		for _, filename := range a.GetFilelist() {
			if _, ok := seen[filename]; !ok {
				seen[filename] = []string{}
			}
			seen[filename] = append(seen[filename], archivename)
		}
		total[archivename] = make([]int64, len(args)+1)
	}

	// Now start adding in the appropriate slot
	for filename, archivelist := range seen {
		filesize := s.GetFileSize(filename)
		if filesize < 0 {
			fmt.Printf("File not found: %s\n", filename)
			continue
		}
		for _, archivename := range archivelist {
			total[archivename][0] += filesize
			total[archivename][len(archivelist)] += filesize
		}
	}

	header := fmt.Sprintf("%-30s\t%-14s%-14s", "Archive", "Total", "Unique")
	for x := 1; x < len(args); x += 1 {
		header += fmt.Sprintf("SHARE-%-8d", x)
	}
	fmt.Printf("%s\n", header)

	for archivename, vals := range total {
		line := fmt.Sprintf("%-30s\t", archivename)
		for _, v := range vals {
			line += fmt.Sprintf("%-14d", v)
		}
		fmt.Printf("%s\n", line)
	}
}
