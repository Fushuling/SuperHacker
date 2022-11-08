package MyOperatePacket4Server

import (
	"encoding/binary"
	"io"
)
/*
type MyRecvSignPacket struct {
	r io.Reader
	sign uint32//根据标志确定封包类型: 心跳包/客户端数据包
}
type MySendSignPacket struct {
	w io.Writer
	sign uint32//根据标志确定封包类型: 心跳包/客户端数据包
}

func  InitRecvSignPacket(r io.Reader) MyRecvSignPacket{
	return MyRecvSignPacket{r:r}
}

func  InitSendSignPacket(w io.Writer) MySendSignPacket{
	return MySendSignPacket{w: w}
}
*/

func ReadSign(r io.Reader) (uint32,error){
	var v uint32
	e:=binary.Read(r, binary.LittleEndian, &v)
	return v,e
}