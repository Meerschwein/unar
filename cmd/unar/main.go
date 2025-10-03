package main

import (
	"flag"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/meerschwein/unar/internal/archives"
)

func init() {
	log.SetFlags(0) // no timestamp

	fmts := slices.Collect(maps.Keys(archives.FormatArchives))
	slices.Sort(fmts)
	flag.StringVar(&format, "f", "", "archive format of the file\npossible values: "+strings.Join(fmts, ", "))

	flag.Usage = func() {
		log.Println("Usage: unar [options] archive")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename = flag.Arg(0)
}

var (
	format   string
	filename string
)

func main() {
	var archFsFn archives.ArchiveFsFn
	var dstPath string

	if format != "" {
		found := false
		archFsFn, found = archives.FormatArchives[format]
		if !found {
			log.Fatal("unknown archive format ", format)
		}
		dstPath = filepath.Base(strings.TrimSuffix(filename, filepath.Ext(filename)))
	} else {
		for suffix, fn := range archives.SuffixArchives {
			if strings.HasSuffix(filename, suffix) {
				archFsFn = fn
				dstPath = filepath.Base(strings.TrimSuffix(filename, suffix))
				goto found
			}
		}

		log.Println("unknown file extension")
		flag.Usage()
		os.Exit(1)

	found:
	}

	_, err := os.Open(dstPath)
	if err == nil {
		log.Fatal(`directory "`, dstPath, `" already exists`)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fs, err := archFsFn(file, info.Size())
	if err != nil {
		log.Fatal(err)
	}

	err = os.CopyFS(dstPath, fs)
	if err != nil {
		log.Fatal(err)
	}
}
