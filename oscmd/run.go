package oscmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type Runner struct {
}

func (runnder Runner) Run(rawCmd string, args ...string) {
	cmd := exec.Command(rawCmd, args...)

	var stdoutBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)

	var stderrBuf bytes.Buffer
	stdErrIn, _ := cmd.StderrPipe()
	stdErr := io.MultiWriter(os.Stderr, &stderrBuf)

	cmd.Start()

	go func() {
		io.Copy(stdout, stdoutIn)
	}()

	go func() {
		io.Copy(stdErr, stdErrIn)
	}()

	cmd.Wait()
}

func (runnder Runner) RunForResult(rawCmd string, args ...string) string {
	cmd := exec.Command(rawCmd, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String()
	}

	return out.String()
}
