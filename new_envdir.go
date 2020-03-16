package goroutine

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

//ReadDir ...
func ReadDir(dir string) (map[string]string, error) {
	env := map[string]string{}
	final := env

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return env, err
	}

	for _, file := range files {
		key := file.Name()

		file, err := os.Open(path.Join(dir, key))
		if err != nil {
			continue
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		line, isTooLong, err := reader.ReadLine()
		if err != nil {
			continue
		}
		if isTooLong {
			continue
		}
		value := string(line)
		if strings.IndexRune(value, '=') >= 0 {
			continue
		}
		final[key] = value
	}

	return final, nil
}

//RunCmd ...
func RunCmd(cmd []string, env map[string]string) int {
	cmdExe := exec.Command(cmd[0], cmd[1:]...)
	cmdExe.Env = os.Environ()
	for k, v := range env {
		cmdExe.Env = append(
			cmdExe.Env,
			k+"="+v,
		)
	}

	cmdExe.Stdin = os.Stdin
	cmdExe.Stdout = os.Stdout
	cmdExe.Stderr = os.Stderr

	if err := cmdExe.Run(); err != nil {
		if exitError, isExit := err.(*exec.ExitError); isExit {
			if statusCode, isOk := exitError.Sys().(syscall.WaitStatus); isOk {
				return int(statusCode.ExitStatus())
			}
		}
		log.Fatalf("Command execution error: %v", err)
		return 1
	}

	return 0
}
