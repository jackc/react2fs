package main

import (
	fsnotify "gopkg.in/fsnotify.v0"
	"os"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	Events  chan fsnotify.Event
	Errors  chan error
	quit    chan bool
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
	return w.watcher.Add(name)
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
		case event.Op != 0:
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

func (w *Watcher) Close() error {
	w.quit <- true

	return w.watcher.Close()
}
