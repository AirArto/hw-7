package goroutine

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

func main() {
	args := os.Args
	dir := args[0]
	cmd := args[1:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatalf("Env directory error: %v", err)
		os.Exit(1)
	}
	statusCode := RunCmd(cmd, env)

	os.Exit(statusCode)
}

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
			return env, err
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		line, isTooLong, err := reader.ReadLine()
		value := string(line)
		if err != nil {
			return env, err
		}
		if isTooLong {
			return env, errors.New("Too many data in env file " + key)
		}
		if strings.IndexRune(value, '=') >= 0 {
			return env, errors.New("Symbol \"=\" in env file " + key)
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
