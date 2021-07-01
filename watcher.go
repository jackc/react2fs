package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	Events  chan fsnotify.Event
	Errors  chan error
	quit    chan bool
	Include *regexp.Regexp
	Exclude *regexp.Regexp
	paths   []string
}

func NewWatcher() (*Watcher, error) {
	var w Watcher
	var err error

	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w.Events = make(chan fsnotify.Event)
	w.Errors = make(chan error)
	w.quit = make(chan bool)

	go w.watch()

	return &w, err
}

func (w *Watcher) Add(name string) error {
	w.paths = append(w.paths, name)
	return w.add(name)
}

func (w *Watcher) add(name string) error {
	err := w.watcher.Add(name)
	if err != nil {
		return err
	}

	stat, err := os.Stat(name)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return nil
	}

	fileInfos, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}

	for _, fi := range fileInfos {
		if fi.IsDir() {
			err = w.Add(filepath.Join(name, fi.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Watcher) watch() {
	for {
		select {
		case event := <-w.watcher.Events:
			if w.isMatchingFile(event.Name) {
				if event.Op == fsnotify.Create {
					stat, err := os.Stat(event.Name)
					if err != nil {
						w.Errors <- err
						continue
					}
					if stat.IsDir() {
						err = w.add(event.Name)
						if err != nil {
							w.Errors <- err
						}
					}
				}

				if event.Op != 0 && event.Op != fsnotify.Chmod {
					select {
					case w.Events <- event:
					case <-w.quit:
						return
					}
				}
			}

		case err := <-w.watcher.Errors:
			select {
			case w.Errors <- err:
			case <-w.quit:
				return
			}
		case <-w.quit:
			return
		}
	}
}

func (w *Watcher) isMatchingFile(name string) bool {
	return w.isWatchedPath(name) && w.isIncludedFile(name) && !w.isExcludedFile(name)
}

// isWatchedPath checks that the file event comes from a path that is being watched. In theory that shouldn't be
// necessary as only events from watched paths should be received. But at least on MacOS that is not true when
// switching branches in Git.
func (w *Watcher) isWatchedPath(name string) bool {
	for _, p := range w.paths {
		if strings.HasPrefix(name, p) {
			return true
		}
	}

	return false
}

func (w *Watcher) isIncludedFile(name string) bool {
	return w.Include == nil || w.Include.MatchString(name)
}

func (w *Watcher) isExcludedFile(name string) bool {
	return w.Exclude != nil && w.Exclude.MatchString(name)
}

func (w *Watcher) Close() error {
	w.quit <- true

	return w.watcher.Close()
}
