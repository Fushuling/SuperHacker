#pragma once
#include <WinSock2.h>
#pragma comment(lib, "ws2_32.lib")
//warning C4996: 'inet_addr': Use inet_pton() or InetPton() instead or define _WINSOCK_DEPRECATED_NO_WARNINGS to disable deprecated API warnings

#include <Windows.h>
#include "MyDEBUG.h"
#include "Utility.h"

#define MaxPacketSize (50*1024*1024)//50MB
#define PacketMagic 0xcafefafa

typedef struct  MyRecvPacketHead
{
	UINT32 Magic;
	INT32 PacketSize;//Go�ж����Ϊuint32,һ�����ݰ�����СΪ2GB
};

typedef struct  MySendPacketHead
{
	UINT32 Sign;
	UINT32 Magic;
	INT32 PacketSize;
};

class MySocket {
private:
	SOCKET sock;

	void* pRecv;
	size_t BodySize;//���body��С

	void* pSend;
	size_t SendMemSize;//��д����ڴ��С

	UINT32 roff;
	UINT32 woff;

	HANDLE m;//�ⲿ���ɵĻ�����,��ֻ����
public:
	//MySocket()
	bool MySocketInit(HANDLE mu);
	MySocket* New();//���Ƴ�һ������ͬsocket�ͻ���������,Ҫ�����ͷ�
	//~MySocket();
	void Free();

	bool Connect(const char* IP, u_short port);
	bool WriteString(const wchar_t* str);
	bool ReadString(wchar_t** str);

	bool WriteBytes(const void* buf,size_t l);
	bool ReadBytes(void** buf, size_t* l);

	bool WriteUint32(UINT32 v);
	bool ReadUint32(UINT32* p);

	bool Send(UINT32 Sign);
	int Recv();
 };