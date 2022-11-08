package main

import (
	"MyOperatePacket4Server"
	"fmt"
	"net"
	"sync"
	"time"
)

var clientList clientList_//全局连接表: 存储当前连接到服务器的所有客户端//记得线程安全要加锁

var currentClient *clientInfo//存储当前选中的客户端

type clientInfo struct{
	num uint//没有使用lsc时均为0
	c net.Conn
	addr string
	lastTime time.Time//读写心跳要加锁
	currentPath string//存储当前用户所处路径
	wait time.Duration//等待一个命令执行的最长时间

	m sync.RWMutex
	r *MyOperatePacket4Server.MyRecvPatket
	s *MyOperatePacket4Server.MySendPatket
}

func (c *clientInfo) ReadTime() time.Time{
	c.m.RLock()
	t:=c.lastTime
	c.m.RUnlock()
	return t
}
func (c *clientInfo) SetTime(){
	c.m.Lock()
	c.lastTime=time.Now()
	c.m.Unlock()
}

//ReadConnByNum n>=1
func (c *clientList_) ReadConnByNum(n uint) *clientInfo{
	c.m.RLock()
	for _,v:=range c.d{
		if v.num==n{
			c.m.RUnlock()
			return v
		}
	}
	c.m.RUnlock()
	return nil
}
func (c *clientInfo) SetNum(n uint){
	c.m.Lock()
	c.num=n
	c.m.Unlock()
}

type clientList_ struct{
	d map[string]*clientInfo
	m sync.RWMutex
}

func (c *clientList_) ReadConn(addr string) *clientInfo{
	c.m.RLock()
	r:=c.d[addr]
	c.m.RUnlock()
	return r
}
func (c *clientList_) WriteConn(addr string,ci *clientInfo){
	c.m.Lock()
	c.d[addr]=ci
	c.m.Unlock()
}
func (c *clientList_) DelConn(addr string,err error){
	c.m.Lock()
	if c.d[addr]==nil{return}
	p:=c.d[addr]
	e:=p.c.Close()
	if e!=nil{
		fmt.Println("Close connect failed, err: ",e)
	}
	p.r.Stop(err)
	if currentClient==p{currentClient=nil}//清除当前已经选中的客户端
	delete(c.d,addr)
	c.m.Unlock()
}
