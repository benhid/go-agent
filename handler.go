package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func handler(data []byte) agentExecRes {
	parts := strings.Split(os.Getenv("fprocess"), " ")

	targetCmd := exec.Command(parts[0], parts[1:]...)

	// Set stdout and stderr.
	var execStdOut, execStdErr bytes.Buffer

	targetCmd.Stdout = &execStdOut
	targetCmd.Stderr = &execStdErr

	stdin, _ := targetCmd.StdinPipe()

	go func() {
		defer stdin.Close()
		// Values written to stdin are passed to targetCmd's standard input.
		_, _ = stdin.Write(data)
	}()

	log.WithField("targetCmd", targetCmd.String()).Info("Running command")

	start := time.Now()
	err := targetCmd.Run()
	elapsed := time.Since(start)

	if err != nil {
		log.WithError(err).Error("Error running command: failed to run")

		return agentExecRes{
			Message:      "Failed to run",
			Error:        err.Error(),
			StdOut:       execStdOut.String(),
			StdErr:       execStdErr.String(),
			ExecDuration: elapsed.Microseconds(),
			MemUsage:     targetCmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss,
		}
	}

	return agentExecRes{
		Message:      "Success",
		StdOut:       execStdOut.String(),
		StdErr:       execStdErr.String(),
		ExecDuration: elapsed.Microseconds(),
		MemUsage:     targetCmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss,
	}
}
