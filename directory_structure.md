| +————病毒端
|                  +————客户端代码   
|                                    +————DLL_Loader	客户端DLL加载器代码
|                                    +————MyCRT                实现必要的CRT函数,不然不能编译
|                                    +————NoCRT_Dll          客户端核心逻辑实现
|                                                    +————DLL_Loader    存放项目编译后成果
|                                                    +————dllmain.cpp    客户端启动入口
|                                                    +————MyPacket.h/cpp    C++版自定义通讯协议实现
|                                                    +————number.h    封包编号(功能编号)、执行结果定义
|                                                    +————RWF.h/cpp    实现了文件操作类
|                                                    +————Utility.h/cpp   功能函数定义
|                                                    +————MyDEBUG.h    调试使用,启用后在debview工具有输出
|                                                    +————pch.h/cpp/framework.h    预编译头文件
|                                                    +————NoCRT_Dll.sln    点击启动VS2019项目
|                  +————Keylogger       键盘记录功能
|                  +————Passwd           密码抓取功能
|                  +————Screenshare   屏幕监控功能
|                  +————Screenshoot    截屏功能
| +————远控端
|                  +————dbUpload 上传数据到数据库功能
|                  +————MyOperatePacket4Server 自定义通信协议包
|                  +————TCPServer  远控端主入口
|                                  +————.idea
|                                  +————ClientInfoList.go   对服务器当前连接的所有客户端进行管理
|                                  +————DefineCommandLine.go  正式运行之前对功能的载入
|                                  +————main.go  主函数
|                                  +————Number.go  封包编号(功能编号)
|                                  +————TCPServer.exe 远控端exe
|                                  +————Utility.go 功能函数编写
|   
|
|
|
|
|
|
|
|
|
|
|