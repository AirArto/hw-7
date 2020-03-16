package goroutine

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	toFile, err := os.OpenFile("testFromFile", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf("\n\t%s", err)
	}

	buffer := []byte("123")
	_, err = toFile.Write(buffer)
	toFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("\n\t%s", err)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		t.Errorf("\n\t%s", err)
	}

	env, err := ReadDir(absDir)
	if err != nil {
		log.Fatalf("Env directory error: %v", err)
	}
	cmd := []string{"./otus"}

	statusCode := RunCmd(cmd, env)
	if statusCode != 123 {
		log.Fatalf("Execution error")
	}

	toFile, err = os.OpenFile("testFromFile", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf("\n\t%s", err)
	}

	buffer = []byte("12")
	_, err = toFile.Write(buffer)
	toFile.Close()

	env, err = ReadDir(absDir)
	if err != nil {
		log.Fatalf("Env directory error: %v", err)
	}
	cmd = []string{"./otus"}

	statusCode = RunCmd(cmd, env)
	if statusCode != 12 {
		log.Fatalf("Execution error")
	}
}
