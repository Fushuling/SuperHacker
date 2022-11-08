package main

//sign封包编号
const (
	HEARTBEAT uint32 = iota
	CLIENT_PACKET uint32 = iota
)

//封包功能编号
const (
	UPLOAD uint32 = iota
	DOWNLOAD uint32 = iota
	LSF uint32 = iota
	CD uint32 = iota
	CMD uint32 = iota
	SUICIDE uint32 = iota
)

const (
	SUCCESS = iota
	FAIL
)