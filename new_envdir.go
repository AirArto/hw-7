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
			newKey := key + "="
			final[newKey] = ""
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
		final[key] = trimValue
	}

	return final, nil
}

//RunCmd ...
func RunCmd(cmd []string, env map[string]string) int {
	cmdExe := exec.Command(cmd[0], cmd[1:]...)
	cmdExe.Env = os.Environ()
	newEnv := []string{}
	for _, value := range cmdExe.Env {
		key := strings.Split(value, "=")[0]
		if _, ok := env[key+"="]; !ok {
			newEnv = append(
				newEnv,
				value,
			)
		}
	}
	cmdExe.Env = newEnv
	for key, value := range env {
		if strings.IndexRune(key, '=') == -1 {
			cmdExe.Env = append(
				cmdExe.Env,
				key+"="+value,
			)
		}
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
