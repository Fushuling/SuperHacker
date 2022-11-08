package MyOperatePacket4Server

import (
	"MyOperatePacket4Server/MyBuf"
	"encoding/binary"
	"errors"
	"golang.org/x/text/encoding/unicode"
	"io"
	"reflect"
	"strconv"
	"unsafe"
)
const maxsize uint32=50*1024*1024//50MB
const magicNum uint32=0xcafefafa
//数据包格式: 魔数+长度+数据(支持3种数据类型int32,string,[]byte)
//数据格式: 字符串(string)和字节数组([]byte)格式均为: 长度+内容(string含终止符)

type MyRecvPatket struct {
	recvReady chan error
	couldRead chan error
	first bool

	reader io.Reader
	magic uint32
	length uint32
	data MyBuf.Buffer
}
type MySendPatket struct {
	writer io.Writer
	magic uint32
	data MyBuf.Buffer
}

func Bytes2string(b []byte) string{
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

func String2bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func  InitSendPacket(w io.Writer) *MySendPatket{
	return &MySendPatket{writer: w,magic: magicNum}
}

func (p *MySendPatket) SendPacket() error{//发送完毕后清空数据包内容
	e:=binary.Write(p.writer,binary.LittleEndian,p.magic)//写魔数
	if e!=nil{return e}
	e=binary.Write(p.writer,binary.LittleEndian,uint32(p.data.Len()))//写长度
	if e!=nil{return e}
	_,e=p.writer.Write(p.data.Bytes())//写内容
	if e!=nil{return e}

	if p.data.Cap()>int(maxsize) {//之前传输了很大的文件时,让GC回收这部分内存
		p.data=*MyBuf.NewBuffer(make([]byte,0,1024))
	}else{
		p.data.Reset()
	}
	return e
}

func (p *MySendPatket) WriteUint32(v uint32) error{
	e:=binary.Write(&p.data,binary.LittleEndian,v)
	return e
}

func (p *MySendPatket) WriteBytes(b []byte) (uint32,error){
	e:=binary.Write(&p.data,binary.LittleEndian,uint32(len(b)))

	if e!=nil{return 0,e}

	n,e:=p.data.Write(b)
	return uint32(n),e
}

func (p *MySendPatket) WriteString(s string) (uint32,error){
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()//UTF8转UTF16
	utf16,e:=encoder.Bytes(String2bytes(s))
	utf16=append(utf16, 0,0)//终止符

	e=binary.Write(&p.data,binary.LittleEndian,uint32(len(utf16)))

	if e!=nil{return 0,e}

	n,e:=p.data.Write(utf16)
	return uint32(n),e
}

//封包的recv和read操作都是协程安全的
//其设计为一个协程(1)验证并接收封包,另一个协程(2)读取封包内容
//服务器上的每个TCP数据接收连接都对应这样一个协程(1)和协程(2)
//本质就是协程(1)提供数据给协程(2);协程(2)读取用户指令,然后生成封包发送给客户端,接着读取协程(1)提供的数据,处理后将执行结果输出
//为防止覆盖当前封包数据,协程(2)只有调用ReadNextPacket()明确指出读完当前封包内容后,协程(1)才会继续接收封包

func  InitRecvPacket(r io.Reader) *MyRecvPatket{
	return &MyRecvPatket{reader:r,magic: magicNum,length: 0,recvReady: make(chan error),couldRead: make(chan error),first: true}
}

func (p *MyRecvPatket) RecvPacket() error{//接收之前清空数据包内容

	if p.first==true{
		p.first=false//第一次没有过程正在读取数据包内容,直接跳过阻塞
	}else{
		e:=<-p.recvReady//阻塞直到得到可以读下一个数据包的通知
		if e!=nil{return e}
	}

	var v uint32
	e:=binary.Read(p.reader, binary.LittleEndian, &v)//读魔数

	if e!=nil && e!=io.EOF{p.couldRead<-e/*通知ReadyToRead()解析封包头出现错误*/;return e}
	if v!=magicNum{
		ne:=errors.New("wrong packet format: magicNum: "+strconv.FormatUint(uint64(v),16))
		p.couldRead<-ne
		return ne
	}

	e=binary.Read(p.reader, binary.LittleEndian, &v)//读长度
	if e!=nil && e!=io.EOF{p.couldRead<-e;return e}

	//简易长度验证,恶意封包可能会导致panic
	if int(v)<0{
		ne:=errors.New("wrong packet format: packet length")
		p.couldRead<-ne
		return ne
	}

	if p.data.Cap()>int(maxsize) && v<maxsize {//之前传输了很大的文件时,让GC回收这部分内存
		p.data=*MyBuf.NewBuffer(make([]byte,0,v))
	}else{
		p.data.Reset()
	}

	_,e=p.data.ReadFull(p.reader,int(v))//读内容,保证这个数据包的内容一定是全部都存储在内存中了

	if e!=nil && e!=io.EOF {
		p.couldRead<-e
		return e
	}

	p.couldRead<-nil
	return nil
}

//注意: 任何Read方法如果返回error!=nil则应该立即停止继续读取,开始错误处理比如打印错误然后让协程返回;其后发生错误后又调用Finish()只会导致阻塞

//ReadUint32 一定要有错误处理,应对数据包格式错误的情况
func (p *MyRecvPatket) ReadUint32() (uint32,error){
	var v uint32
	e:=binary.Read(&p.data, binary.LittleEndian, &v)
	if e!=nil {p.recvReady<-e/*通知RecvPacket()解析封包内容出现错误*/;return 0,e}
	return v,e
}

//ReadBytes 一定要有错误处理,应对数据包格式错误的情况
func (p *MyRecvPatket) ReadBytes() ([]byte,error) {
	var v uint32
	e:=binary.Read(&p.data, binary.LittleEndian, &v)

	if e!=nil {p.recvReady<-e;return nil,e}
	if v > uint32(p.data.Len()){
		ne:=errors.New("wrong packet format: bytes length")
		p.recvReady<-ne
		return nil,ne
	}

	b:=p.data.Next(int(v))

	return b,e
}

// ReadString 一定要有错误处理,应对数据包格式错误的情况
func (p *MyRecvPatket) ReadString() (string,error){
	var v uint32
	e:=binary.Read(&p.data, binary.LittleEndian, &v)

	if e!=nil {p.recvReady<-e;return "",e}
	if v > uint32(p.data.Len()) {
		ne:=errors.New("wrong packet format: string length")
		p.recvReady<-ne
		return "",ne
	}

	s:=p.data.Next(int(v))//终止符也要一起读出来

	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()//UTF16转UTF8
	utf8, e := decoder.Bytes(s[:len(s)-2])//去掉终止符
	if e!=nil {p.recvReady<-e;return "",e}

	return Bytes2string(utf8),e
}

// ReadyToRead 读取数据包的内容前调用,此函数会阻塞直到有内容可以读取
// 一定要有错误处理,应对连接断开等情况
func (p *MyRecvPatket) ReadyToRead() error{
	return <-p.couldRead
}

// ReadNextPacket 读完数据包的内容后调用,来通知RecvPacket()可以开始接收下一个数据包了,此函数会阻塞直到有内容可以读取
// 一定要有错误处理,应对连接断开等情况
func (p *MyRecvPatket) ReadNextPacket() error {
	p.recvReady<-nil
	return <-p.couldRead
}

// Finish 读完所有需要读取的数据包后调用,调用时确保RecvPacket()没有处于阻塞状态
func (p *MyRecvPatket) Finish() {
	p.first=true
}

//Stop 停止所有接收封包的行为,释放资源
func (p *MyRecvPatket) Stop(e error) {
	select {
	case p.recvReady<-e:
	default:
	}
	select {
	case p.couldRead<-e:
	default:
	}
}