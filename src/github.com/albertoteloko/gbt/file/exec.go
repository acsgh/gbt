package file

import (
	"os/exec"
	"io/ioutil"
	"os"
)

func Exec(cmd string, args ...string) (string, string, error) {
	return ExecWithEnv(cmd, args, os.Environ())
}

func ExecWithEnv(cmd string, args []string, env []string) (string, string, error) {
	command := exec.Command(cmd, args...)
	command.Env = env

	cmdOut, _ := command.StdoutPipe()
	cmdErr, _ := command.StderrPipe()

	err := command.Start()
	if err != nil {
		return "", "", err
	}

	stdOutput, err := ioutil.ReadAll(cmdOut)
	errOutput, err := ioutil.ReadAll(cmdErr)

	return string(stdOutput), string(errOutput), err
}
