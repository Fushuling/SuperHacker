package main

import (
	"MyOperatePacket"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)
var wait chan int
func recv(reader io.Reader){
	r:= MyOperatePacket4Server.InitRecvPacket(reader)
	go myread(r)
	s,_:= MyOperatePacket4Server.ReadSign(reader)
	fmt.Println("Sign ",s)
	r.RecvPacket()//对应r.ReadyToRead()->r.ReadNextPacket()

	s,_= MyOperatePacket4Server.ReadSign(reader)
	fmt.Println("Sign ",s)
	r.RecvPacket()//对应r.ReadNextPacket()->r.Finish()
}

func myread(r *MyOperatePacket4Server.MyRecvPatket){
	//只是测试,实际使用必须错误处理
	r.ReadyToRead()
	fmt.Println(r.ReadUint32())
	fmt.Println(r.ReadString())
	fmt.Println(r.ReadString())
	b,_:=r.ReadBytes()
	fmt.Println(string(b))
	r.ReadNextPacket()
	fmt.Println(r.ReadUint32())
	fmt.Println(r.ReadString())
	fmt.Println(r.ReadString())
	b,_=r.ReadBytes()
	fmt.Println(string(b))
	r.Finish()

	wait<-0
}
func main(){
	wait=make(chan int)
	buf:=bytes.NewBuffer(nil)

	//写入SignPacket
	var v uint32=666
	binary.Write(buf, binary.LittleEndian, &v)

	s:= MyOperatePacket4Server.InitSendPacket(buf)
	s.WriteUint32(123)
	s.WriteString("hello ")
	s.WriteString("my ")
	s.WriteBytes([]byte("friend"))
	s.SendPacket()//第一次发送封包

	//写入SignPacket
	v=777
	binary.Write(buf, binary.LittleEndian, &v)

	s.WriteUint32(555)
	s.WriteString("hi ")
	s.WriteString("our ")
	s.WriteBytes([]byte("partner"))
	s.SendPacket()//第二次发送封包
	go recv(buf)

	<-wait
}