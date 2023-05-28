---
code: null
id: MDU6SXNzdWUyMzExMTg3MzI=
databaseId: 231118732
number: 55
url: https://github.com/kamilsk/retry/issues/55
title: "timeout at db example in go 1.7"
issueType: null
labels:
  - "type:critical"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-05-24T17:43:20Z
updatedAt: 2017-05-25T10:03:37Z
lastEditedAt: null
closedAt: 2017-05-25T10:03:37Z
---

# timeout at db example in go 1.7

```
goroutine 499 [running]:
panic(0x65b040, 0xc4200c7980)
	/usr/local/go/src/runtime/panic.go:500 +0x1a1
testing.startAlarm.func1()
	/usr/local/go/src/testing/testing.go:918 +0x10b
created by time.goFunc
	/usr/local/go/src/time/sleep.go:154 +0x44

goroutine 1 [chan send]:
github.com/kamilsk/retry_test.Example_dbConnectionRestore()
	/go/src/github.com/kamilsk/retry/example_db_test.go:88 +0x274
testing.runExample(0x6c485f, 0x1b, 0x6e2070, 0x6c40d2, 0x19, 0x0, 0x1)
	/usr/local/go/src/testing/example.go:114 +0x1a0
testing.RunExamples(0x6e20a0, 0x7f2e60, 0x3, 0x3, 0x410601)
	/usr/local/go/src/testing/example.go:38 +0x15e
testing.(*M).Run(0xc420045ee8, 0x6c6bed)
	/usr/local/go/src/testing/testing.go:744 +0xbd
main.main()
	github.com/kamilsk/retry/_test/_testmain.go:124 +0x1bd

goroutine 17 [syscall, locked to thread]:
runtime.goexit()
	/usr/local/go/src/runtime/asm_amd64.s:2086 +0x1

goroutine 13 [syscall]:
syscall.Syscall(0x0, 0x3, 0xc4200ee400, 0x200, 0x0, 0x5f8, 0x0)
	/usr/local/go/src/syscall/asm_linux_amd64.s:18 +0x5
syscall.read(0x3, 0xc4200ee400, 0x200, 0x200, 0x7fd422219000, 0x0, 0x0)
	/usr/local/go/src/syscall/zsyscall_linux_amd64.go:783 +0x55
syscall.Read(0x3, 0xc4200ee400, 0x200, 0x200, 0x0, 0xc420026800, 0xc420072d00)
	/usr/local/go/src/syscall/syscall_unix.go:161 +0x49
os.(*File).read(0xc4200240a0, 0xc4200ee400, 0x200, 0x200, 0x4a885b, 0x65b240, 0x200)
	/usr/local/go/src/os/file_unix.go:228 +0x4f
os.(*File).Read(0xc4200240a0, 0xc4200ee400, 0x200, 0x200, 0x6, 0x0, 0xc42001ce20)
	/usr/local/go/src/os/file.go:101 +0x59
bytes.(*Buffer).ReadFrom(0xc42000e3f0, 0x7d73c0, 0xc4200240a0, 0xc42001ce88, 0x1, 0x70)
	/usr/local/go/src/bytes/buffer.go:176 +0x155
io.copyBuffer(0x7d6900, 0xc42000e3f0, 0x7d73c0, 0xc4200240a0, 0x0, 0x0, 0x0, 0x69a1a0, 0x474e01, 0xc42000e3f0)
	/usr/local/go/src/io/io.go:384 +0x323
io.Copy(0x7d6900, 0xc42000e3f0, 0x7d73c0, 0xc4200240a0, 0xc420013650, 0x8, 0x0)
	/usr/local/go/src/io/io.go:360 +0x68
testing.runExample.func1(0xc4200240a0, 0xc420014c60)
	/usr/local/go/src/testing/example.go:68 +0xa0
created by testing.runExample
	/usr/local/go/src/testing/example.go:75 +0xfe

goroutine 24 [IO wait]:
net.runtime_pollWait(0x7fd422179f00, 0x72, 0x0)
	/usr/local/go/src/runtime/netpoll.go:160 +0x59
net.(*pollDesc).wait(0xc4200cc0d0, 0x72, 0xc42001c5e0, 0xc420012130)
	/usr/local/go/src/net/fd_poll_runtime.go:73 +0x38
net.(*pollDesc).waitRead(0xc4200cc0d0, 0x7d88c0, 0xc420012130)
	/usr/local/go/src/net/fd_poll_runtime.go:78 +0x34
net.(*netFD).accept(0xc4200cc070, 0x0, 0x7d7480, 0xc4200ba720)
	/usr/local/go/src/net/fd_unix.go:419 +0x238
net.(*TCPListener).accept(0xc4200ce020, 0x434b9b, 0xc42001c690, 0x51b29d)
	/usr/local/go/src/net/tcpsock_posix.go:132 +0x2e
net.(*TCPListener).Accept(0xc4200ce020, 0x6e2280, 0xc420078200, 0x7db440, 0xc4200170b0)
	/usr/local/go/src/net/tcpsock.go:222 +0x49
net/http.(*Server).Serve(0xc4200c2100, 0x7dabc0, 0xc4200ce020, 0x0, 0x0)
	/usr/local/go/src/net/http/server.go:2273 +0x1ce
net/http/httptest.(*Server).goServe.func1(0xc4200ca360)
	/usr/local/go/src/net/http/httptest/server.go:235 +0x6d
created by net/http/httptest.(*Server).goServe
	/usr/local/go/src/net/http/httptest/server.go:236 +0x5c

goroutine 14 [chan receive]:
database/sql.(*DB).connectionOpener(0xc420094420)
	/usr/local/go/src/database/sql/sql.go:730 +0x4a
created by database/sql.Open
	/usr/local/go/src/database/sql/sql.go:493 +0x1e9

goroutine 28 [IO wait]:
net.runtime_pollWait(0x7fd422179e40, 0x72, 0x7)
	/usr/local/go/src/runtime/netpoll.go:160 +0x59
net.(*pollDesc).wait(0xc4200cc1b0, 0x72, 0xc4200359d0, 0xc420012130)
	/usr/local/go/src/net/fd_poll_runtime.go:73 +0x38
net.(*pollDesc).waitRead(0xc4200cc1b0, 0x7d88c0, 0xc420012130)
	/usr/local/go/src/net/fd_poll_runtime.go:78 +0x34
net.(*netFD).Read(0xc4200cc150, 0xc4200fe000, 0x1000, 0x1000, 0x0, 0x7d88c0, 0xc420012130)
	/usr/local/go/src/net/fd_unix.go:243 +0x1a1
net.(*conn).Read(0xc4200ce040, 0xc4200fe000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
	/usr/local/go/src/net/net.go:173 +0x70
net/http.(*persistConn).Read(0xc4200f4000, 0xc4200fe000, 0x1000, 0x1000, 0x530af0, 0xc420035b58, 0x404ead)
	/usr/local/go/src/net/http/transport.go:1261 +0x154
bufio.(*Reader).fill(0xc4200ca900)
	/usr/local/go/src/bufio/bufio.go:97 +0x10c
bufio.(*Reader).Peek(0xc4200ca900, 0x1, 0xc420035bbd, 0x1, 0x0, 0xc420014ae0, 0x0)
	/usr/local/go/src/bufio/bufio.go:129 +0x62
net/http.(*persistConn).readLoop(0xc4200f4000)
	/usr/local/go/src/net/http/transport.go:1418 +0x1a1
created by net/http.(*Transport).dialConn
	/usr/local/go/src/net/http/transport.go:1062 +0x4e9

goroutine 29 [select]:
net/http.(*persistConn).writeLoop(0xc4200f4000)
	/usr/local/go/src/net/http/transport.go:1646 +0x3bd
created by net/http.(*Transport).dialConn
	/usr/local/go/src/net/http/transport.go:1063 +0x50e

goroutine 12 [IO wait]:
net.runtime_pollWait(0x7fd422179d80, 0x72, 0x8)
	/usr/local/go/src/runtime/netpoll.go:160 +0x59
net.(*pollDesc).wait(0xc42000e220, 0x72, 0xc4200367b0, 0xc420012130)
	/usr/local/go/src/net/fd_poll_runtime.go:73 +0x38
net.(*pollDesc).waitRead(0xc42000e220, 0x7d88c0, 0xc420012130)
	/usr/local/go/src/net/fd_poll_runtime.go:78 +0x34
net.(*netFD).Read(0xc42000e1c0, 0xc42007f000, 0x1000, 0x1000, 0x0, 0x7d88c0, 0xc420012130)
	/usr/local/go/src/net/fd_unix.go:243 +0x1a1
net.(*conn).Read(0xc420024050, 0xc42007f000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
	/usr/local/go/src/net/net.go:173 +0x70
net/http.(*connReader).Read(0xc4200ba740, 0xc42007f000, 0x1000, 0x1000, 0x50d669, 0x7ee0e0, 0x0)
	/usr/local/go/src/net/http/server.go:586 +0x144
bufio.(*Reader).fill(0xc420014840)
	/usr/local/go/src/bufio/bufio.go:97 +0x10c
bufio.(*Reader).ReadSlice(0xc420014840, 0xa, 0x0, 0x1e, 0x0, 0x33, 0x0)
	/usr/local/go/src/bufio/bufio.go:330 +0xb5
bufio.(*Reader).ReadLine(0xc420014840, 0xc4200be5a0, 0xf0, 0xf0, 0x6b45a0, 0x4936c3, 0x7f6408)
	/usr/local/go/src/bufio/bufio.go:359 +0x37
net/textproto.(*Reader).readLineSlice(0xc420017170, 0xc420036aa8, 0xc420036aa8, 0x410e68, 0xf0, 0x6b45a0)
	/usr/local/go/src/net/textproto/reader.go:55 +0x5e
net/textproto.(*Reader).ReadLine(0xc420017170, 0xc4200be5a0, 0xc, 0x401f73, 0xc420036c78)
	/usr/local/go/src/net/textproto/reader.go:36 +0x2f
net/http.readRequest(0xc420014840, 0xc420036c00, 0xc4200be5a0, 0x0, 0x0)
	/usr/local/go/src/net/http/request.go:793 +0xa5
net/http.(*conn).readRequest(0xc420078200, 0x7db380, 0xc420050980, 0x0, 0x0, 0x0)
	/usr/local/go/src/net/http/server.go:765 +0x10d
net/http.(*conn).serve(0xc420078200, 0x7db380, 0xc420050980)
	/usr/local/go/src/net/http/server.go:1532 +0x3d3
created by net/http.(*Server).Serve
	/usr/local/go/src/net/http/server.go:2293 +0x44d

goroutine 15 [sleep]:
time.Sleep(0xf4240)
	/usr/local/go/src/runtime/time.go:59 +0xe1
github.com/kamilsk/retry_test.Example_dbConnectionRestore.func2(0xc420094420, 0x7db3c0, 0xc4200124a8, 0xc420014d20, 0x1, 0xf4240)
	/go/src/github.com/kamilsk/retry/example_db_test.go:84 +0x257
created by github.com/kamilsk/retry_test.Example_dbConnectionRestore
	/go/src/github.com/kamilsk/retry/example_db_test.go:86 +0x1f5
FAIL	github.com/kamilsk/retry	1.882s
```

<!-- 2017-05-24T17:50Z https://github.com/kamilsk/retry/issues/55#issuecomment-303800892
https://travis-ci.org/kamilsk/retry/jobs/235744082
-->
