// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package simpleini

import (
	"os"
	"testing"
	"time"
)

func TestEventStringWithValue(t *testing.T) {
	for opMask, expectedString := range map[fsnOp]string{
		fsnChmod | fsnCreate: `"/usr/someFile": CREATE|CHMOD`,
		fsnRename:         `"/usr/someFile": RENAME`,
		fsnRemove:         `"/usr/someFile": REMOVE`,
		fsnWrite | fsnChmod:  `"/usr/someFile": WRITE|CHMOD`,
	} {
		event := fsnEvent{Name: "/usr/someFile", Op: opMask}
		if event.String() != expectedString {
			t.Fatalf("Expected %s, got: %v", expectedString, event.String())
		}

	}
}

func TestEventOpStringWithValue(t *testing.T) {
	expectedOpString := "WRITE|CHMOD"
	event := fsnEvent{Name: "someFile", Op: fsnWrite | fsnChmod}
	if event.Op.String() != expectedOpString {
		t.Fatalf("Expected %s, got: %v", expectedOpString, event.Op.String())
	}
}

func TestEventOpStringWithNoValue(t *testing.T) {
	expectedOpString := ""
	event := fsnEvent{Name: "testFile", Op: 0}
	if event.Op.String() != expectedOpString {
		t.Fatalf("Expected %s, got: %v", expectedOpString, event.Op.String())
	}
}

// TestWatcherClose tests that the goroutine started by creating the watcher can be
// signalled to return at any time, even if there is no goroutine listening on the events
// or errors channels.
func TestWatcherClose(t *testing.T) {
	t.Parallel()

	name := tempMkFile(t, "")
	w := newWatcher(t)
	err := w.Add(name)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
	// Allow the watcher to receive the event.
	time.Sleep(time.Millisecond * 100)

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}
}
