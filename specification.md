# 设计思路

大部分代码均使用 TDD 思想构建。

为了解析一个 `.ini` 配置文件，从顺序上说，我们要能够读入文件，然后读取配置文件中的每一行，依据每一行的不同类型去解析。于是我们首先编写了读入文件的测试 `readFromFile`。

但是，在单元测试中去测试本地 IO 相比一般的测试来说更繁琐，并且引入了外部的不确定性。所以我们选择使用依赖注入，将读入文件时使用的 `os.File` 替换为 `io.Reader`，并把后续代码抽离出来。此时，我们就可以使用其他实现了 `io.Reader` 的方法作为真实文件的替代，例如 `strings.newReader`。

这部分完成后，我们考虑如何解析配置文件。最简单的任务，便是解析一行键值对。于是我们编写了 `ParseLine` 测试，要求其能读取类似 `app_mode = development` 的字符串。当这一测试通过后，我们再逐渐增加其他的情形，例如首尾有空格的情况。就这样在 test fail 和 test pass 的多个周期后，我们完成了多种类型的“行”的读取，包括键值对、纯注释、节标题、混合注释等。

当完成了行的读取，我们就可以对整个配置文件进行读取，只需要把读入的配置文件按行分割即可。我们同样编写了测试，此时我们的测例是一个完整的 `.ini` 文件，而非仅仅是其中的一行。

为了支持不同的注释符号，我们更改了相应的测试，并把代码中写死的注释符号抽离出来，成为函数的参数。为了在不同的系统平台下支持不同的注释符号默认值，我们编写了 GitHub Actions Workflow，可在不同的系统上自动测试代码。

此时，我们完成了配置文件的一次性读取功能。而为了支持自动重载，我们引入了 `fsnotify` 的代码（由于不能使用第三方包，且手动支持不同平台的 `WatchFile` 也是繁琐、不太具有意义的事情，所以我们只好把相应的代码复制了一下）。此部分的测试相对难以编写，偏向集成测试的范畴，故我们省略了相关的测试代码，由 fsnotify 和我们原有的解析配置的测试代码作为一定的保证。

## 单元测试结果

```
=== RUN   TestPollerWithBadFd
--- PASS: TestPollerWithBadFd (0.00s)
=== RUN   TestPollerWithData
--- PASS: TestPollerWithData (0.00s)
=== RUN   TestPollerWithWakeup
--- PASS: TestPollerWithWakeup (0.00s)
=== RUN   TestPollerWithClose
--- PASS: TestPollerWithClose (0.00s)
=== RUN   TestPollerWithWakeupAndData
--- PASS: TestPollerWithWakeupAndData (0.00s)
=== RUN   TestPollerConcurrent
--- PASS: TestPollerConcurrent (0.15s)
=== RUN   TestFsnotifyMultipleOperations
    fsnotify_integration_test.go:104: event received: "/tmp/fsnotify506916073/TestFsnotifySeq.testfile": CREATE
    fsnotify_integration_test.go:104: event received: "/tmp/fsnotify506916073/TestFsnotifySeq.testfile": WRITE
    fsnotify_integration_test.go:104: event received: "/tmp/fsnotify506916073/TestFsnotifySeq.testfile": RENAME
    fsnotify_integration_test.go:104: event received: "/tmp/fsnotify506916073/TestFsnotifySeq.testfile": CREATE
    fsnotify_integration_test.go:180: calling Close()
    fsnotify_integration_test.go:182: waiting for the event channel to become closed...
    fsnotify_integration_test.go:185: event channel closed
--- PASS: TestFsnotifyMultipleOperations (0.67s)
=== RUN   TestFsnotifyMultipleCreates
    fsnotify_integration_test.go:217: event received: "/tmp/fsnotify799360771/TestFsnotifySeq.testfile": CREATE
    fsnotify_integration_test.go:217: event received: "/tmp/fsnotify799360771/TestFsnotifySeq.testfile": WRITE
    fsnotify_integration_test.go:217: event received: "/tmp/fsnotify799360771/TestFsnotifySeq.testfile": REMOVE
    fsnotify_integration_test.go:217: event received: "/tmp/fsnotify799360771/TestFsnotifySeq.testfile": CREATE
    fsnotify_integration_test.go:217: event received: "/tmp/fsnotify799360771/TestFsnotifySeq.testfile": WRITE
    fsnotify_integration_test.go:217: event received: "/tmp/fsnotify799360771/TestFsnotifySeq.testfile": WRITE
    fsnotify_integration_test.go:306: calling Close()
    fsnotify_integration_test.go:308: waiting for the event channel to become closed...
    fsnotify_integration_test.go:311: event channel closed
--- PASS: TestFsnotifyMultipleCreates (0.77s)
=== RUN   TestFsnotifyDirOnly
    fsnotify_integration_test.go:356: event received: "/tmp/fsnotify111290502/TestFsnotifyDirOnly.testfile": CREATE
    fsnotify_integration_test.go:356: event received: "/tmp/fsnotify111290502/TestFsnotifyDirOnly.testfile": WRITE
    fsnotify_integration_test.go:356: event received: "/tmp/fsnotify111290502/TestFsnotifyDirOnly.testfile": REMOVE
    fsnotify_integration_test.go:356: event received: "/tmp/fsnotify111290502/TestFsnotifyEventsExisting.testfile": REMOVE
    fsnotify_integration_test.go:408: calling Close()
    fsnotify_integration_test.go:410: waiting for the event channel to become closed...
    fsnotify_integration_test.go:413: event channel closed
--- PASS: TestFsnotifyDirOnly (0.57s)
=== RUN   TestFsnotifyDeleteWatchedDir
    fsnotify_integration_test.go:458: event received: "/tmp/fsnotify810836525/TestFsnotifyEventsExisting.testfile": REMOVE
    fsnotify_integration_test.go:458: event received: "/tmp/fsnotify810836525/TestFsnotifyEventsExisting.testfile": REMOVE
    fsnotify_integration_test.go:458: event received: "/tmp/fsnotify810836525": REMOVE
--- PASS: TestFsnotifyDeleteWatchedDir (0.50s)
=== RUN   TestFsnotifySubDir
    fsnotify_integration_test.go:504: event received: "/tmp/fsnotify402150824/sub": CREATE
    fsnotify_integration_test.go:504: event received: "/tmp/fsnotify402150824/TestFsnotifyFile1.testfile": CREATE
    fsnotify_integration_test.go:504: event received: "/tmp/fsnotify402150824/sub": REMOVE
    fsnotify_integration_test.go:504: event received: "/tmp/fsnotify402150824/TestFsnotifyFile1.testfile": REMOVE
    fsnotify_integration_test.go:561: calling Close()
    fsnotify_integration_test.go:563: waiting for the event channel to become closed...
    fsnotify_integration_test.go:566: event channel closed
--- PASS: TestFsnotifySubDir (0.71s)
=== RUN   TestFsnotifyRename
    fsnotify_integration_test.go:602: event received: "/tmp/fsnotify820805351/TestFsnotifyEvents.testfile": CREATE
    fsnotify_integration_test.go:602: event received: "/tmp/fsnotify820805351/TestFsnotifyEvents.testfile": WRITE
    fsnotify_integration_test.go:602: event received: "/tmp/fsnotify820805351/TestFsnotifyEvents.testfile": RENAME
    fsnotify_integration_test.go:602: event received: "/tmp/fsnotify820805351/TestFsnotifyEvents.testfileRenamed": CREATE
    fsnotify_integration_test.go:602: event received: "/tmp/fsnotify820805351/TestFsnotifyEvents.testfile": RENAME
    fsnotify_integration_test.go:637: calling Close()
    fsnotify_integration_test.go:639: waiting for the event channel to become closed...
    fsnotify_integration_test.go:642: event channel closed
--- PASS: TestFsnotifyRename (0.51s)
=== RUN   TestFsnotifyRenameToCreate
    fsnotify_integration_test.go:684: event received: "/tmp/fsnotify469083418/TestFsnotifyEvents.testfileRenamed": CREATE
    fsnotify_integration_test.go:713: calling Close()
    fsnotify_integration_test.go:715: waiting for the event channel to become closed...
    fsnotify_integration_test.go:718: event channel closed
--- PASS: TestFsnotifyRenameToCreate (0.51s)
=== RUN   TestFsnotifyRenameToOverwrite
    fsnotify_integration_test.go:772: event received: "/tmp/fsnotify705774940/TestFsnotifyEvents.testfileRenamed": CREATE
    fsnotify_integration_test.go:801: calling Close()
    fsnotify_integration_test.go:803: waiting for the event channel to become closed...
    fsnotify_integration_test.go:806: event channel closed
--- PASS: TestFsnotifyRenameToOverwrite (0.51s)
=== RUN   TestRemovalOfWatch
    fsnotify_integration_test.go:844: No event received, as expected.
--- PASS: TestRemovalOfWatch (0.60s)
=== RUN   TestFsnotifyAttrib
    fsnotify_integration_test.go:901: event received: "/tmp/fsnotify405292917/TestFsnotifyAttrib.testfile": CHMOD
    fsnotify_integration_test.go:901: event received: "/tmp/fsnotify405292917/TestFsnotifyAttrib.testfile": WRITE
    fsnotify_integration_test.go:901: event received: "/tmp/fsnotify405292917/TestFsnotifyAttrib.testfile": CHMOD
    fsnotify_integration_test.go:979: calling Close()
    fsnotify_integration_test.go:981: waiting for the event channel to become closed...
    fsnotify_integration_test.go:984: event channel closed
--- PASS: TestFsnotifyAttrib (1.52s)
=== RUN   TestFsnotifyClose
--- PASS: TestFsnotifyClose (0.05s)
=== RUN   TestFsnotifyFakeSymlink
    fsnotify_integration_test.go:1053: Created bogus symlink
    fsnotify_integration_test.go:1039: event received: "/tmp/fsnotify032815983/zzznew": CREATE
    fsnotify_integration_test.go:1075: calling Close()
--- PASS: TestFsnotifyFakeSymlink (0.51s)
=== RUN   TestCyclicSymlink
--- PASS: TestCyclicSymlink (0.51s)
=== RUN   TestConcurrentRemovalOfWatch
    fsnotify_integration_test.go:1131: regression test for race only present on darwin
--- SKIP: TestConcurrentRemovalOfWatch (0.00s)
=== RUN   TestClose
--- PASS: TestClose (0.01s)
=== RUN   TestRemoveWithClose
--- PASS: TestRemoveWithClose (0.01s)
=== RUN   TestEventStringWithValue
--- PASS: TestEventStringWithValue (0.00s)
=== RUN   TestEventOpStringWithValue
--- PASS: TestEventOpStringWithValue (0.00s)
=== RUN   TestEventOpStringWithNoValue
--- PASS: TestEventOpStringWithNoValue (0.00s)
=== RUN   TestWatcherClose
=== PAUSE TestWatcherClose
=== RUN   TestParseFromString
--- PASS: TestParseFromString (0.00s)
=== RUN   TestParseFromStringWithCommentDelimiter
--- PASS: TestParseFromStringWithCommentDelimiter (0.00s)
=== RUN   TestParseLineUnix
=== RUN   TestParseLineUnix/key_value_pair_line
=== RUN   TestParseLineUnix/key_value_pair_line_without_spaces
=== RUN   TestParseLineUnix/value_contains_=
=== RUN   TestParseLineUnix/comment
=== RUN   TestParseLineUnix/comment_contains_key_pair_line
=== RUN   TestParseLineUnix/key_pair_line_with_comment
=== RUN   TestParseLineUnix/section
=== RUN   TestParseLineUnix/comment_contains_section
=== RUN   TestParseLineUnix/section_with_spaces
--- PASS: TestParseLineUnix (0.00s)
    --- PASS: TestParseLineUnix/key_value_pair_line (0.00s)
    --- PASS: TestParseLineUnix/key_value_pair_line_without_spaces (0.00s)
    --- PASS: TestParseLineUnix/value_contains_= (0.00s)
    --- PASS: TestParseLineUnix/comment (0.00s)
    --- PASS: TestParseLineUnix/comment_contains_key_pair_line (0.00s)
    --- PASS: TestParseLineUnix/key_pair_line_with_comment (0.00s)
    --- PASS: TestParseLineUnix/section (0.00s)
    --- PASS: TestParseLineUnix/comment_contains_section (0.00s)
    --- PASS: TestParseLineUnix/section_with_spaces (0.00s)
=== RUN   TestParseLineWindows
=== RUN   TestParseLineWindows/key_value_pair_line
=== RUN   TestParseLineWindows/key_value_pair_line_without_spaces
=== RUN   TestParseLineWindows/value_contains_=
=== RUN   TestParseLineWindows/comment
=== RUN   TestParseLineWindows/comment_contains_key_pair_line
=== RUN   TestParseLineWindows/key_pair_line_with_comment
=== RUN   TestParseLineWindows/section
=== RUN   TestParseLineWindows/comment_contains_section
=== RUN   TestParseLineWindows/section_with_spaces
--- PASS: TestParseLineWindows (0.00s)
    --- PASS: TestParseLineWindows/key_value_pair_line (0.00s)
    --- PASS: TestParseLineWindows/key_value_pair_line_without_spaces (0.00s)
    --- PASS: TestParseLineWindows/value_contains_= (0.00s)
    --- PASS: TestParseLineWindows/comment (0.00s)
    --- PASS: TestParseLineWindows/comment_contains_key_pair_line (0.00s)
    --- PASS: TestParseLineWindows/key_pair_line_with_comment (0.00s)
    --- PASS: TestParseLineWindows/section (0.00s)
    --- PASS: TestParseLineWindows/comment_contains_section (0.00s)
    --- PASS: TestParseLineWindows/section_with_spaces (0.00s)
=== RUN   TestReadFromFile
--- PASS: TestReadFromFile (0.00s)
=== RUN   TestReadFromReader
--- PASS: TestReadFromReader (0.00s)
=== CONT  TestWatcherClose
--- PASS: TestWatcherClose (0.10s)
PASS
ok      github.com/robinWongM/go-simpleini/simpleini    8.216s
```