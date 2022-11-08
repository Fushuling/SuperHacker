# SuperHacker

## 🎉项目简介

👊由苟安懋和张城赫共同制作的远控程序，包含远控端（服务端）（苟安懋），病毒（客户端）（张城赫）。

## 👏启动方式

🐱远控端：打开TCPServer，运行TCPServer.exe，进入远控端主页面后使用-h查看帮助，默认监听端口为9999

🐀病毒端：打开客户端代码，运行启动客户端.exe

## 👍功能实现及使用说明

注：后面苟安懋用golang写的几个功能张城赫出于对病毒端大小的严苛要求并没有直接加进去，我们本来准备用插件形式外带在病毒端，但后面因为时间问题没有完成。所以如果想要在远控端使用这几个功能只能用文件浏览功能打开文件夹然后用命令执行功能执行文件夹里的exe文件😥😥

- **能够向客户端发送shell指令并且在客户端执行之后接受返回的指令执行结果**
  lsc （列出所有病毒端序号和ip）
  c 1（选择病毒端，c+序号或者ip）
  cd D：\ （默认文件路径为空路径，需要主动选择文件夹；空路径时使用lsf可以查看目标有哪些盘）
  cmd -w   （进入命令执行功能，如不需要回显则不用添加-w参数，没有输出的命令在添加了-w的情况下会让客户端阻塞，使用wait命令可以限制最大阻塞时间，否则默认为无限阻塞）
  cmd /c(或powershell)+所需指令，如：cmd /c ping 127.0.0.1 (含有中文会乱码，仅仅影响显示其他不会影响)

- **采用模块化程序设计思想**

  主程序与分支功能分离
  
- **服务器可以与多个客户端通信**

   进入远控端主页面，使用lsc，查看所有与远控端连接的客户端ip，使用c 序号选中目标病毒端

- **实现心跳功能（定期检查客户端存活）**

   使用lsc，添加-l=时间间隔(s)，只会显示在从此时刻开始计算，往前推这么长的时间间隔内发送过心跳包的客户端，没有发送过的客户端会被认定为断开了连接。(注意：默认心跳包为5s发一次，客户端有掉线重连机制)

- **文件浏览**

   选中病毒端后打开文件夹，输入lsf (空路径时使用lsf可以查看目标有哪些盘) ：

- **文件的上传及下载**

   download:  download 客户端中待下载的文件路径 下载到的路径;带有空格的路径必须用引号括起来,可以是完整或相对路径

   upload:  upload 待上传的文件路径 上传到客户端的文件路径;带有空格的路径必须用引号括起来,可以是完整或相对路径

- **病毒端程序自杀**

   在远控端服务界面选中病毒端后输入suicide删除当前病毒端。

   客户端其实是一个dll不是那个exe，exe只是加载dll用的， suicide实现的是在内存中把dll卸载了，可以对抗内存取证，防止被反溯源 

- **浏览器密码抓取**

  在病毒端/Passwd中，用命令执行功能执行main.exe，启动后可以将数据输出在命令行中，并保存在一个txt文件中。

  目前只实现了谷歌浏览器账号密码的抓取，其他浏览器的密码的破解有点搞不懂，解密脚本参考了HackBrowserData。

- **调用数据库来实现相关信息的储存**

  参考[golang保存图片到数据库](https://blog.csdn.net/benben_2015/article/details/79223120)，通过调用"gopkg.in/mgo.v2"，实现了golang将数据上传到mongodb的功能。

  在远控端文件夹中打开dbUpload，启动main.exe，在本地浏览器中打开127.0.0.1:8000/entrance，选择图片并且上传，可以将图片的ID和imgurl上传至数据库中。（默认为本地mongodb，test数据库，test表，如果有账号密码验证需要自己去源码修改）（k8s上可能没法启动，这完全就是一个新的程序，可以在本地测试）

- **自定义通讯协议**

  由于最初苟安懋的病毒端与远控端都是在golang环境下编写的，运用了golang的原生库，该原生库的数据传输格式很难用CPP直接解析，所以我们使用TCP协议，制定了常用数据类型的传输和打包格式，实现了一个通讯协议来进行双端的交互，同时增强了交互流量的隐蔽性，具体内容在MyOperatePacket4Server和MyPacket.cpp中。

- **键盘记录**

  在病毒端/Keylogger中，用命令执行功能执行main.exe。

  调用了 `github.com/kindlyfire/go-keylogger` 库，可以实现对键盘输入的记录，并将数据存储到一个txt文件中

- **病毒端截屏**

  在病毒端\Screenshoot中，用命令执行功能执行main.exe。

  调用了 `github.com/kbinani/screenshot` ，可以实现截屏功能，然后在病毒端所在位置生成一个截屏图片

- **屏幕监控**（只能本地监控本地）

  病毒端\Screenshare中，用命令执行功能执行main.exe，在本地打开127.0.0.1:8080
  
  录屏实现参考了[govnc](https://github.com/maka00/govnc)
  
  在最初实现了截屏功能后，我们想到能不能实现更有趣的功能，比如调用病毒端摄像头之类的，然后我们在外网上找到了相关的文章，不过要想调用病毒端的摄像头需要的权限很高，而且很容易被杀毒软件杀死，所以这个想法就此作废。后来我们退而求其次想到了能不能实现对病毒端屏幕的实时监控，就用上面那个截屏功能，通过不断的截图然后放在缓存区内，可以在本地的一个网页里实时监控屏幕，但由于这样cpu占用会比较高，对于远程传输不利，所以最后我们也只做到本地监控本地，不能把数据传输出去。不过或许可以把该网页做成外网可以访问的公网网页，这样就能在其他地方实时监控病毒端屏幕了
  
- **客户端特色**
客户端核心逻辑实现在DLL中，exe只是加载DLL的外壳，使用DLL形式的木马比传统EXE形式的有更多优势，比如不需要实现进程隐藏，隐蔽性较好，因为可以让DLL加载到任意进程中，进程只是外壳。同时不用处理或者经过一些处理后可以和很多技术相结合，比如各种DLL注入技术，DLL劫持技术，DLL捆绑到EXE中等等，更容易实现权限维持。同时DLL编写环境较为苛刻，没有使用C++ Runtime Library，DLLMain就是DLL的直接入口点，很多C标准库的函数都不能使用，需要自己实现，在编译过程中就有很多坑，唯一能使用的外部就是WINDOWS的API，所以编写起来较为困难，但是编译出来的DLL只有8KB，比一些Hello World的程序还要小。
