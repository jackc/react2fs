package main

import (
	"os"
	"os/exec"
)

type Process struct {
	CmdPath  string
	Argv     []string
	procAttr *os.ProcAttr
	process  *os.Process
}

// StartProcess starts a new process from from argv.
func StartProcess(argv []string) (*Process, error) {
	var wp Process
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
func (wp *Process) Restart() (err error) {
	err = wp.Kill()
	if err != nil {
		return err
	}

	wp.process, err = os.StartProcess(wp.CmdPath, wp.Argv, wp.procAttr)
	return err
}

// Kill stops the process without restarting it
func (wp *Process) Kill() (err error) {
	if wp.process != nil {
		err = wp.process.Kill()
		if err != nil {
			return err
		}

		err = wp.process.Release()
		if err != nil {
			return err
		}

		wp.process = nil
	}

	return nil
}
