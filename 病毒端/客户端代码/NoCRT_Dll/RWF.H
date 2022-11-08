#pragma once
#include <windows.h>
#include "Utility.h"

void LongPath2W(const wchar_t* LongPath, wchar_t* Path);//转换为可以处理的长文件名路径,lstrlenW(Path)=_tcslen(LongPath)+lstrlenW(L"\\\\?\\")
BOOL SuperCreateDirectory(const wchar_t* Path);//创建多级目录,必须是C:\a\b\c这种格式,支持长文件名

class RWF{
	private:
		HANDLE fh;
	public:
		~RWF();

		BOOL openFile(
			LPCTSTR lpFileName,
			DWORD dwCreationDisposition = OPEN_EXISTING,
			DWORD dwDesiredAccess = GENERIC_READ | GENERIC_WRITE,
			DWORD dwShareMode = FILE_SHARE_READ,
			DWORD dwFlagsAndAttributes = FILE_ATTRIBUTE_NORMAL
		);

		BOOL close();

		DWORD read(//会清零缓冲区
			LPVOID lpBuffer,
			DWORD nNumberOfBytesToRead = 0//下同,只有使用HeapAlloc()申请的内存空间可以不填这个参数,并且将默认缓冲区大小为整个申请的空间的大小
		);

		DWORD write(//在当前文件指针处覆盖式写入数据,文件指针自动加上写入的数据的长度
			LPCVOID lpBuffer,
			DWORD nNumberOfBytesToWrite = 0
		);

		DWORD insert(//在当前文件指针处插入数据,文件指针自动加上插入的数据的长度,返回实际插入长度,失败返回0
			LPCVOID lpBuffer,
			DWORD nNumberOfBytesToWrite = 0
		);

		DWORD del(DWORD dwLength);//在当前指针处删除指定长度的数据,不设置指针,返回实际删除长度

		BOOL flushFileBuffers();//将文件立即写入磁盘

		DWORD setPointer(
			LONG lDistanceToMove,//如果欲设置的值超出文件范围则指针将被设置到文件末尾(=FileSize)
			DWORD dwMoveMethod = FILE_BEGIN
		);

		DWORD getPointer();

		DWORD getFileSize();
};