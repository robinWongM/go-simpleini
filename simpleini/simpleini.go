package simpleini

import (
	"io"
	"io/ioutil"
	"os"
)

func Watch(filename string, listener func(Configuration)) (Configuration, error) {
	// Read file first
	iniConf, err := ParseFromFile(filename)
	if err != nil {
		return nil, err
	}

	// Start watch
	watchFileChanges(filename, func() {
		iniConf, err := ParseFromFile(filename)
		if err != nil {
			return
		}
		listener(iniConf)
	})

	return iniConf, nil
}

func ParseFromFile(filename string) (Configuration, error) {
	iniContent, err := ReadFromFile(filename)
	if err != nil {
		return nil, err
	}

	iniConf, err := ParseFromString(iniContent)
	if err != nil {
		return nil, err
	}
	return iniConf, nil
}

func watchFileChanges(filename string, listener func()) error {
	watcher, err := NewWatcher()
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
					listener()
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

func ReadFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return "", err
	}
	return ReadFromReader(file)
}

func ReadFromReader(reader io.Reader) (string, error) {
	if content, err := ioutil.ReadAll(reader); err == nil {
		return string(content), nil
	} else {
		return "", err
	}
}
