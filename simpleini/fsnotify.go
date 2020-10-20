// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package simpleini

import (
	"bytes"
	"errors"
	"fmt"
)

// fsnEvent represents a single file system notification.
type fsnEvent struct {
	Name string // Relative path to the file or directory.
	Op   fsnOp  // File operation that triggered the event.
}

// fsnOp describes a set of file operations.
type fsnOp uint32

// These are the generalized file operations that can trigger a notification.
const (
	fsnCreate fsnOp = 1 << iota
	fsnWrite
	fsnRemove
	fsnRename
	fsnChmod
)

func (op fsnOp) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if op&fsnCreate == fsnCreate {
		buffer.WriteString("|CREATE")
	}
	if op&fsnRemove == fsnRemove {
		buffer.WriteString("|REMOVE")
	}
	if op&fsnWrite == fsnWrite {
		buffer.WriteString("|WRITE")
	}
	if op&fsnRename == fsnRename {
		buffer.WriteString("|RENAME")
	}
	if op&fsnChmod == fsnChmod {
		buffer.WriteString("|CHMOD")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()[1:] // Strip leading pipe
}

// String returns a string representation of the event in the form
// "file: REMOVE|WRITE|..."
func (e fsnEvent) String() string {
	return fmt.Sprintf("%q: %s", e.Name, e.Op.String())
}

// Errors that returned by fsnotify
var (
	ErrFsnotifyEventOverflow = errors.New("fsnotify queue overflow")
)
