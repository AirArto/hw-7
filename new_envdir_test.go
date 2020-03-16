package goroutine

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	toFile, err := os.OpenFile("testFromFile", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf("\n\t%s", err)
	}
	defer toFile.Close()

	buffer := []byte("123")
	_, err = toFile.Write(buffer)
}
