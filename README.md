# SuperHacker

## 🎉项目简介

## 👏启动方式

## 👍功能实现及使用说明

- **能够向客户端发送shell指令并且在客户端执行之后接受返回的指令执行结果**

- **采用模块化程序设计思想**
- **服务器可以与多个客户端通信**

- **实现心跳功能（定期检查客户端存活）**

- **文件浏览**

- **文件的上传及下载**

- **病毒端程序自杀**

- **浏览器密码抓取**

- **调用数据库来实现相关信息的储存**

- **自定义通讯协议**

- **键盘记录**

- **病毒端截屏**

- **屏幕监控**（只能本地监控本地）

  
- **客户端特色**
客户端核心逻辑实现在DLL中，exe只是加载DLL的外壳，使用DLL形式的木马比传统EXE形式的有更多优势，比如不需要实现进程隐藏，隐蔽性较好，因为可以让DLL加载到任意进程中，进程只是外壳。同时不用处理或者经过一些处理后可以和很多技术相结合，比如各种DLL注入技术，DLL劫持技术，DLL捆绑到EXE中等等，更容易实现权限维持。同时DLL编写环境较为苛刻，没有使用C++ Runtime Library，DLLMain就是DLL的直接入口点，很多C标准库的函数都不能使用，需要自己实现，在编译过程中就有很多坑，唯一能使用的外部就是WINDOWS的API，所以编写起来较为困难，但是编译出来的DLL只有8KB，比一些Hello World的程序还要小。
