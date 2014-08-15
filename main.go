package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const Version = "0.1.0"

var options struct {
	dir     string
	include string
	exclude string
	version bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options] command\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&options.dir, "dir", ".", "directories to watch (separate multiple directories with commas)")
	flag.StringVar(&options.include, "include", "", "only watch files matching this regexp")
	flag.StringVar(&options.exclude, "exclude", "", "don't watch files matching this regexp")
	flag.BoolVar(&options.version, "version", false, "print version and exit")
	flag.Parse()

	if options.version {
		fmt.Printf("react2fs v%v\n", Version)
		os.Exit(0)
	}

	cmd := flag.Args()
	if len(cmd) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	watcher, err := NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	includes := strings.Split(options.include, ",")
	for _, glob := range includes {
		glob := strings.TrimSpace(glob)
		watcher.Includes = append(watcher.Includes, glob)
	}

	excludes := strings.Split(options.exclude, ",")
	for _, glob := range excludes {
		glob := strings.TrimSpace(glob)
		watcher.Excludes = append(watcher.Excludes, glob)
	}

	dirs := strings.Split(options.dir, ",")
	for _, d := range dirs {
		err = watcher.Add(d)
		if err != nil {
			log.Fatal("Unable to watch directory:", err)
		}
	}

	var nowRunning *Process
	nowRunning, err = StartProcess(cmd)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Println(event)
			err = nowRunning.Restart()
			if err != nil {
				log.Println("error:", err)
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}
