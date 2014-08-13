package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var options struct {
	dir     string
	pattern string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options] command\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&options.dir, "dir", ".", "directories to watch (separate multiple directories with commas)")
	flag.StringVar(&options.pattern, "pattern", ".*", "only watch files matching this regexp")
	flag.Parse()

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

	for _, p := range strings.Split(options.pattern, ",") {
		re, err := regexp.Compile(p)
		if err != nil {
			log.Fatal("Invalid pattern:", err)
		}
		watcher.Patterns = append(watcher.Patterns, re)
	}

	dirs := strings.Split(options.dir, ",")
	for _, d := range dirs {
		err = watcher.Add(d)
		if err != nil {
			log.Fatal(err)
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
