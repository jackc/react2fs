package main

import (
	"os"
	"os/exec"
)

type WatchedProcess struct {
	CmdPath  string
	Argv     []string
	procAttr *os.ProcAttr
	process  *os.Process
}

// StartWatchedProcess starts a new process from from argv.
func StartWatchedProcess(argv []string) (*WatchedProcess, error) {
	var wp WatchedProcess
	var err error

	wp.CmdPath, err = exec.LookPath(argv[0])
	if err != nil {
		return nil, err
	}

	wp.Argv = argv

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	wp.procAttr = &os.ProcAttr{
		Dir:   workingDir,
		Files: []*os.File{nil, os.Stdout, os.Stderr},
		Env:   os.Environ(),
	}

	err = wp.Restart()
	return &wp, err
}

// Restart restarts the watched process
func (wp *WatchedProcess) Restart() (err error) {
	if wp.process != nil {
		err = wp.process.Kill()
		if err != nil {
			return err
		}
		err = wp.process.Release()
		if err != nil {
			return err
		}
	}

	wp.process, err = os.StartProcess(wp.CmdPath, wp.Argv, wp.procAttr)
	return err
}
