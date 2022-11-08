#include "MyPacket.h"

bool MySocket::MySocketInit(HANDLE mu)
{
	this->pRecv = NULL;
	this->BodySize = 0;
	this->pSend = NULL;
	this->SendMemSize = 0;
	this->roff = 0;
	this->woff = 0;

	if (mu == NULL) {
		this->sock = INVALID_SOCKET;
		return false;
	}
	this->m = mu;

	WSADATA wsaData;
	if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
		this->sock = INVALID_SOCKET;
		return false;
	}

	this->sock = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
	if (this->sock == INVALID_SOCKET) {
		//WSAGetLastError()
		//freeaddrinfo(result);
		WSACleanup();
		return false;
	}
	return true;
}

MySocket* MySocket::New()
{
	MySocket* c = (MySocket*)MyAlloc(sizeof(MySocket));
	if (c == NULL)
		return NULL;
	c->pRecv = NULL;
	c->BodySize = 0;
	c->pSend = NULL;
	c->SendMemSize = 0;
	c->roff = 0;
	c->woff = 0;
	c->m = this->m;
	c->sock = this->sock;
	return c;
}

void MySocket::Free()
{
	WaitForSingleObject(this->m, INFINITE);
	if (this->sock != INVALID_SOCKET) {
		shutdown(this->sock, SD_BOTH);
		closesocket(this->sock);
		WSACleanup();
		this->sock = INVALID_SOCKET;
		MyFree(this->pRecv);
		MyFree(this->pSend);
		this->pRecv = NULL;
		this->pSend = NULL;
	}
	ReleaseMutex(this->m);
}

bool MySocket::Connect(const char* IP, u_short port)
{
	if (this->sock == INVALID_SOCKET) {
		if (!this->MySocketInit(this->m))
			return false;
	}

	sockaddr_in sockAddr;
	memset(&sockAddr, 0, sizeof(sockAddr));

	sockAddr.sin_family = PF_INET;
	sockAddr.sin_addr.s_addr = inet_addr(IP);
	sockAddr.sin_port = htons(port);

	if (connect(this->sock, (SOCKADDR*)&sockAddr, sizeof(SOCKADDR)) == SOCKET_ERROR) {
		//closesocket(this->sock);
		return false;
	}
	return true;
}

bool MySocket::WriteString(const wchar_t* str)
{
	return this->WriteBytes(str, (MyStrLenW(str) + 1) * sizeof(wchar_t));
	//size_t l = (MyStrLenW(str) + 1) * sizeof(wchar_t);//包含终止符的总字节数
	//if (this->pSend == NULL) {
	//	this->woff = 0;
	//	this->SendMemSize = l;
	//	this->pSend = MyAlloc(l);
	//}
	//else if (this->woff + l > this->SendMemSize) {
	//	this->SendMemSize += l;
	//	this->pSend = MyReAlloc(this->pSend, this->SendMemSize);
	//}
	//if (this->pSend == NULL) {
	//	MyDebug("WriteString: Alloc Memory failed");
	//	return false;
	//}
	//MyMemCpy((char*)this->pSend + this->woff, str, l);
	//this->woff += l;
	//return true;
}

bool MySocket::ReadString(wchar_t** str)
{
	void* buf; size_t l;
	if (!this->ReadBytes(&buf, &l))
		return false;
	if (*(wchar_t*)((char*)buf + l - sizeof(wchar_t)) != 0)
		return false;
	*str = (wchar_t*)buf;
	return true;
	//if (this->pRecv != NULL) {
	//	if (str == NULL) {
	//		MyDebug("Readstr failed: null pointer");
	//		return false;
	//	}
	//	if (this->roff + sizeof(INT32) >= this->BodySize) {
	//		MyDebug("Readstr failed: not enough length 1");
	//		return false;
	//	}

	//	INT32 stringLen = *(INT32*)((char*)this->pRecv + this->roff);
	//	this->roff += sizeof(INT32);

	//	if (stringLen <= 0) {
	//		MyDebug("Readstr failed: wrong length");
	//		return false;
	//	}
	//	if (this->roff + (UINT32)stringLen > this->BodySize) {
	//		MyDebug("Readstr failed: not enough length 2");
	//		return false;
	//	}
	//	if (*(wchar_t*)((char*)this->pRecv + this->roff + stringLen - sizeof(wchar_t))!=0) {
	//		MyDebug("Readstr failed: not end with \\0");
	//		return false;
	//	}

	//	*str = (wchar_t*)((char*)this->pRecv + this->roff);
	//	this->roff += stringLen;
	//	return true;
	//}
	//return false;
}

bool MySocket::WriteBytes(const void* buf, size_t l)
{
	if (this->pSend == NULL) {
		this->woff = sizeof(MySendPacketHead);
		this->SendMemSize = sizeof(MySendPacketHead) + l + sizeof(INT32);
		this->pSend = MyAlloc(this->SendMemSize);
		if (this->pSend == NULL) {
			MyDebug("WriteBytes: Alloc Memory failed");
			return false;
		}
	}
	else if (this->woff + l + sizeof(INT32) > this->SendMemSize) {
		this->SendMemSize = this->woff + l + sizeof(INT32);
		void* t = MyReAlloc(this->pSend, this->SendMemSize);
		if (t == NULL) {
			MyFree(this->pSend);
			this->pSend = NULL;
			MyDebug("WriteBytes: Alloc Memory failed");
			return false;
		}
		this->pSend = t;
	}

	*(INT32*)((char*)this->pSend + this->woff) = l;
	this->woff += sizeof(INT32);
	MyMemCpy((char*)this->pSend + this->woff, buf, l);
	this->woff += l;
	return true;
}

bool MySocket::ReadBytes(void** buf, size_t* l)
{
	if (this->pRecv != NULL) {
		if (buf == NULL || l == NULL) {
			MyDebug("ReadBytes failed: null pointer");
			return false;
		}
		if (this->roff + sizeof(INT32) >= this->BodySize) {
			MyDebug("ReadBytes failed: not enough length 1");
			return false;
		}
		INT32 byteLen = *(INT32*)((char*)this->pRecv + this->roff);
		this->roff += sizeof(INT32);
		
		if (byteLen<=0){
			MyDebug("ReadBytes failed: wrong length");
			return false;
		}
		if (this->roff + (UINT32)byteLen > this->BodySize) {
			MyDebug("ReadBytes failed: not enough length 2");
			return false;
		}

		*buf = (char*)this->pRecv + this->roff;
		*l = byteLen;

		this->roff += byteLen;
		return true;
	}
	return false;
}

bool MySocket::WriteUint32(UINT32 v)
{
	if (this->pSend == NULL) {
		this->woff = sizeof(MySendPacketHead);
		this->SendMemSize = sizeof(MySendPacketHead) + sizeof(UINT32);
		this->pSend = MyAlloc(this->SendMemSize);
		if (this->pSend == NULL) {
			MyDebug("WriteUint32: Alloc Memory failed");
			return false;
		}
	}
	else if (this->woff + sizeof(UINT32) > this->SendMemSize) {
		this->SendMemSize = this->woff + sizeof(UINT32);
		void* t = MyReAlloc(this->pSend, this->SendMemSize);
		if (t == NULL) {
			MyFree(this->pSend);
			this->pSend = NULL;
			MyDebug("WriteUint32: Alloc Memory failed");
			return false;
		}
		this->pSend = t;
	}
	

	*(UINT32*)((char*)this->pSend + this->woff) = v;
	this->woff += sizeof(UINT32);
	return true;
}

bool MySocket::ReadUint32(UINT32* p)
{
	if (this->pRecv != NULL) {
		if (p == NULL) {
			MyDebug("ReadUint32 failed: null pointer");
			return false;
		}
		if (this->roff + sizeof(UINT32) > this->BodySize) {
			MyDebug("ReadUint32 failed: not enough length");
			return false;
		}

		*p = *(UINT32*)((char*)this->pRecv + this->roff);
		this->roff += sizeof(UINT32);
		return true;
	}
	return false;
}

bool MySocket::Send(UINT32 Sign)
{
	if (this->sock != INVALID_SOCKET) {
		if (this->pSend == NULL) {
			this->woff = sizeof(MySendPacketHead);
			this->SendMemSize = sizeof(MySendPacketHead);
			this->pSend = MyAlloc(this->SendMemSize);
			if (this->pSend == NULL) {
				return false;
			}
		}
		MySendPacketHead sph;
		sph.Sign = Sign;
		sph.Magic = PacketMagic;
		sph.PacketSize = this->woff - sizeof(MySendPacketHead);
		*(MySendPacketHead*)this->pSend = sph;

		WaitForSingleObject(this->m, INFINITE);
		int c = 0; int sendLen = this->woff > sizeof(MySendPacketHead) ? this->woff : sizeof(UINT32);//当封包没有内容时只发送Sign,如心跳包
		do
		{
			int r = send(this->sock, (const char*)this->pSend, sendLen, 0);
			if (r == SOCKET_ERROR) {
				MyDebug("Recv Error %d", WSAGetLastError());
				return false;
			}
			c += r;
		} while (c < sendLen);
		ReleaseMutex(this->m);
		this->woff = sizeof(MySendPacketHead);
		return true;
	}
	return false;
}

int MySocket::Recv()
{
	if (this->sock != INVALID_SOCKET) {
		this->roff = 0;
		this->BodySize = 0;

		MyRecvPacketHead RecvHead;
		int r,c=0;
		do
		{
			r = recv(this->sock, (char*)&RecvHead + c, sizeof(MyRecvPacketHead) - c, NULL);
			if (r == SOCKET_ERROR) {//连接出错
				MyDebug("Recv RecvBody: Failed");
				return 0;
			}
			else if (r == 0) {//连接关闭
				MyDebug("Recv RecvBody: conn close");
				return 2;
			}
			c += r;
		} while (c < sizeof(MyRecvPacketHead));

		if (RecvHead.Magic != PacketMagic) {
			MyDebug("Wrong PacketMagic");
			return 0;
		}
		
		if (RecvHead.PacketSize <= 0) {
			MyDebug("Wrong PacketSize");
			return 0;
		}

		if (this->pRecv == NULL) {
			this->pRecv = MyAlloc(RecvHead.PacketSize);
			if (this->pRecv == NULL) {
				MyDebug("Recv: Alloc Memory failed");
				return 0;
			}
		}else{
			size_t hs = MyMemSize(this->pRecv);
			if (RecvHead.PacketSize > hs || hs > MaxPacketSize) {
				void* t = MyReAlloc(this->pRecv, RecvHead.PacketSize);
				if (t == NULL) {
					MyFree(this->pRecv);
					this->pRecv = NULL;
					MyDebug("Recv: Alloc Memory failed");
					return 0;
				}
				this->pRecv = t;
			}
		}
		

		this->BodySize = RecvHead.PacketSize;
		c = 0;
		do
		{
			r = recv(this->sock, (char*)this->pRecv + c, RecvHead.PacketSize - c, NULL);
			if (r == SOCKET_ERROR) {//连接出错
				MyDebug("Recv RecvBody Failed");
				return 0;
			}
			else if (r == 0) {//连接关闭
				MyDebug("Recv RecvBody: conn close");
				return 2;
			}
			c += r;
		} while (c < RecvHead.PacketSize);
		return 1;
	}
	return 0;
}
