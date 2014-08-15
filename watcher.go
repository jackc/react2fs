package main

import (
	fsnotify "gopkg.in/fsnotify.v0"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	Events  chan fsnotify.Event
	Errors  chan error
	quit    chan bool
	Include *regexp.Regexp
	Exclude *regexp.Regexp
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
		var event fsnotify.Event
		var err error

		select {
		case event = <-w.watcher.Events:
			if event.Op == fsnotify.Create {
				stat, err := os.Stat(event.Name)
				if err != nil {
					w.Errors <- err
					continue
				}
				if stat.IsDir() {
					err = w.Add(event.Name)
					if err != nil {
						w.Errors <- err
					}
				}
			}
		case err = <-w.watcher.Errors:
		case <-w.quit:
			return
		}

		switch {
		case event.Op != 0 && w.isMatchingFile(event.Name):
			select {
			case w.Events <- event:
			case <-w.quit:
				return
			}
		case err != nil:
			select {
			case w.Errors <- err:
			case <-w.quit:
				return
			}
		}
	}
}

func (w *Watcher) isMatchingFile(name string) bool {
	return w.isIncludedFile(name) && !w.isExcludedFile(name)
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
