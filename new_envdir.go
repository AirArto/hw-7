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

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return env, err
	}

	for _, file := range files {
		key := file.Name()
		if strings.IndexRune(key, '=') >= 0 {
			continue
		}

		file, err := os.Open(path.Join(dir, key))
		if err != nil {
			continue
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			continue
		}

		if fi.Size() == 0 {
			env[key] = ""
			continue
		}

		reader := bufio.NewReader(file)
		line, isTooLong, err := reader.ReadLine()
		if err != nil {
			continue
		}
		if isTooLong {
			continue
		}
		value := string(line)
		trimValue := strings.TrimSpace(value)
		env[key] = trimValue
	}

	return env, nil
}

//RunCmd ...
func RunCmd(cmd []string, env map[string]string) int {
	cmdExe := exec.Command(cmd[0], cmd[1:]...)
	for key, value := range env {
		if value != "" {
			os.Setenv(key, value)
		} else {
			os.Unsetenv(key)
		}
	}
	cmdExe.Env = os.Environ()
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
