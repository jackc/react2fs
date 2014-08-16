package main

import (
	"errors"
	"os"
	"os/exec"
	"time"
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
func (wp *Process) Kill() error {
	if wp.process != nil {
		process := wp.process
		wp.process = nil

		// Wait to check error of Kill because process may already be dead. We only
		// care about the error if the process never finished and we were unable to
		// kill it.
		killErr := process.Kill()

		waitDone := make(chan bool)
		waitErr := make(chan error)

		go func() {
			_, err := process.Wait()
			if err != nil {
				waitErr <- err
			} else {
				waitDone <- true
			}
		}()

		select {
		case <-waitDone:
			return nil
		case err := <-waitErr:
			return err
		case <-time.After(10 * time.Second):
			if killErr != nil {
				return killErr
			}
			return errors.New("Timeout waiting for process to terminate")
		}
	}

	return nil
}
