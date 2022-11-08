module TCPServer

go 1.17

replace MyOperatePacket4Server => ../MyOperatePacket4Server

require (
	MyOperatePacket4Server v0.0.0-00010101000000-000000000000
	github.com/jessevdk/go-flags v1.5.0
)

require (
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
	golang.org/x/text v0.3.7 // indirect
)
