package main

import (
	"testing"
	"time"
)

func TestProcessRestartExitedProcess(t *testing.T) {
	t.Parallel()

	wp, err := StartProcess([]string{"sleep", "0"})
	if err != nil {
		t.Fatal(err)
	}
	defer wp.Kill()

	pid := wp.process.Pid

	// wait for process to end -- not using wp.process.Wait() because that cleans
	// up the child process such that the Kill in Restart will fail.
	time.Sleep(time.Second)

	err = wp.Restart()
	if err != nil {
		t.Fatal(err)
	}

	if wp.process.Pid == pid {
		t.Fatal("Restart should have changed pid, but it didn't")
	}
}

func TestProcessRestartRunningProcess(t *testing.T) {
	t.Parallel()

	wp, err := StartProcess([]string{"sleep", "10"})
	if err != nil {
		t.Fatal(err)
	}
	defer wp.Kill()

	pid := wp.process.Pid

	err = wp.Restart()
	if err != nil {
		t.Fatal(err)
	}

	if wp.process.Pid == pid {
		t.Fatal("Restart should have changed pid, but it didn't")
	}
}
