package oscmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func Run(rawCmd string, args ...string) {
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

func RunForResult(rawCmd string, args ...string) (string, error) {
	output, err := exec.Command(rawCmd, args...).Output()
	return string(output), err
}
