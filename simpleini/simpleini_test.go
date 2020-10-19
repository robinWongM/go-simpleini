package simpleini

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func createTempFile(t *testing.T, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func TestReadFromFile(t *testing.T) {
	expected := `[section]
test1=test2
test3=test4
`

	file, removeFile := createTempFile(t, expected)
	defer removeFile()

	file.Seek(0, 0)

	got, err := ReadFromReader(file)
	if err != nil {
		t.Errorf("ReadFromFile error, %e", err)
	} else if got != expected {
		t.Errorf("ReadFromFile failed, got %v, expected %v", got, expected)
	}

}

func TestReadFromReader(t *testing.T) {
	expected := `[section]
test1=test2
test3=test4
`
	reader := strings.NewReader(expected)

	got, err := ReadFromReader(reader)
	if err != nil {
		t.Errorf("ReadFromReader error, %e", err)
	} else if got != expected {
		t.Errorf("ReadFromReader failed, got %v, expected %v", got, expected)
	}

}
