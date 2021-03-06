// Package simpleini provides a simple way to read and reload `.ini` configuration files.
package simpleini

import (
	"io"
	"io/ioutil"
	"os"
)

// Listener interface defines a listener function.
type Listener interface {
	Listen(Configuration)
}

// Watch parses the given configuration file and watching its changes.
// listener will be invoked when configuration file changes.
func Watch(filename string, listener Listener) (Configuration, error) {
	// Read file first
	iniConf, err := parseFromFile(filename)
	if err != nil {
		return nil, err
	}

	// Start watch
	watchFileChanges(filename, func() {
		iniConf, err := parseFromFile(filename)
		if err != nil {
			return
		}
		listener.Listen(iniConf)
	})

	return iniConf, nil
}

func parseFromFile(filename string) (Configuration, error) {
	iniContent, err := readFromFile(filename)
	if err != nil {
		return nil, err
	}

	iniConf, err := parseFromString(iniContent)
	if err != nil {
		return nil, err
	}
	return iniConf, nil
}

func watchFileChanges(filename string, callback func()) error {
	watcher, err := fsnNewWatcher()
	if err != nil {
		return err
	}

	if err = watcher.Add(filename); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnWrite == fsnWrite {
					callback()
				}
			case err, ok := <-watcher.Errors:
				if err != nil {
					panic(err)
				}
				if !ok {
					return
				}
			}
		}
	}()

	return nil
}

func readFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return "", err
	}
	return readFromReader(file)
}

func readFromReader(reader io.Reader) (string, error) {
	if content, err := ioutil.ReadAll(reader); err == nil {
		return string(content), nil
	} else {
		return "", err
	}
}
