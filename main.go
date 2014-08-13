package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	fsnotify "gopkg.in/fsnotify.v0"
)

var options struct {
	dir string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options] command\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&options.dir, "dir", ".", "directory to watch")
	flag.Parse()

	cmd := flag.Args()
	if len(cmd) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(options.dir)
	if err != nil {
		log.Fatal(err)
	}

	var nowRunning *WatchedProcess
	nowRunning, err = StartWatchedProcess(cmd)
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

			if event.Op == fsnotify.Create {
				stat, err := os.Stat(event.Name)
				if err != nil {
					log.Println("error:", err)
				}
				if stat.IsDir() {
					err = watcher.Add(event.Name)
					if err != nil {
						log.Println("error:", err)
					}
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}

}
