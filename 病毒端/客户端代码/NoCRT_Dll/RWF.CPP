#include "RWF.H"
//#include <iostream>
//using namespace std;
void LongPath2W(const wchar_t* LongPath, wchar_t* Path) {//转换为可以处理的长文件名路径,lstrlenW(Path)=_tcslen(LongPath)+lstrlenW(L"\\\\?\\")
	wchar_t prefix[] = L"\\\\?\\";
	size_t i = 0;
	for (; i < (sizeof(prefix) / sizeof(wchar_t)) - 1; i++) {
		Path[i] = prefix[i];
	}
	size_t l = i;
	for (; i < l + MyStrLenW(LongPath) + 1; i++) {
		Path[i] = LongPath[i - l];
	}
}

BOOL SuperCreateDirectory(const wchar_t* Path) {//创建多级目录,必须是C:\a\b\c这种格式,支持长文件名
	const size_t l = sizeof(L"\\\\?\\") - 1;// MyStrLenW(L"\\\\?\\");
	wchar_t* p = (wchar_t*)MyAlloc((l + MyStrLenW(Path) + 1) * sizeof(wchar_t));
	if (p == 0)
		return false;
	LongPath2W(Path, p);
	for (size_t i = l + 3; i < MyStrLenW(p) + 1;) {
		if (p[i] == L'\\') {
			p[i] = 0;
			CreateDirectoryW(p, NULL);
			p[i] = L'\\';
			i++;
			/*while (p[i] == L'\\') {
				i++;
			}*/
			if (p[i] == L'\\') {
				MyFree(p);
				return false;
			}
			continue;
		}
		else if (p[i] == 0) {
			CreateDirectoryW(p, NULL);
			break;
		}
		i++;
	}
	MyFree(p);
	return true;
}

RWF::~RWF()
{
	if(this->fh!=0)
		CloseHandle(this->fh);
}

BOOL RWF::openFile(
			LPCTSTR lpFileName,
			DWORD dwCreationDisposition,
			DWORD dwDesiredAccess,
			DWORD dwShareMode,
			DWORD dwFlagsAndAttributes
			)
{
	this->fh=CreateFile(lpFileName,dwDesiredAccess,dwShareMode,NULL,dwCreationDisposition,dwFlagsAndAttributes,NULL);
	return INVALID_HANDLE_VALUE!=this->fh;
}

BOOL RWF::close()
{
	BOOL b = CloseHandle(this->fh);
	this->fh = 0;
	return b;
}

DWORD RWF::read(
			LPVOID lpBuffer,
			DWORD nNumberOfBytesToRead
			)
{
	//cout<<"V1"<<nNumberOfBytesToRead;
	if (nNumberOfBytesToRead == 0)
		nNumberOfBytesToRead = HeapSize(GetProcessHeap(), 0, lpBuffer);

	DWORD lpNumberOfBytesRead;
	ReadFile(this->fh,lpBuffer,nNumberOfBytesToRead,&lpNumberOfBytesRead,NULL);

	return lpNumberOfBytesRead;
}

DWORD RWF::write(
			LPCVOID lpBuffer,
			DWORD nNumberOfBytesToWrite
			)
{
	if (nNumberOfBytesToWrite == 0)
		nNumberOfBytesToWrite = HeapSize(GetProcessHeap(), 0, lpBuffer);
	
	DWORD lpNumberOfBytesWritten;
	WriteFile(this->fh,lpBuffer,nNumberOfBytesToWrite,&lpNumberOfBytesWritten,NULL);

	return lpNumberOfBytesWritten;
}

DWORD RWF::insert(
				LPCVOID lpBuffer,
				DWORD nNumberOfBytesToWrite
				)
{
	if (nNumberOfBytesToWrite == 0)
		nNumberOfBytesToWrite = HeapSize(GetProcessHeap(), 0, lpBuffer);
	
	DWORD lpNumberOfBytesWritten;
	DWORD pointer=this->getPointer();
	DWORD s=this->getFileSize()-pointer;
	LPTSTR p=0;

	if (s>0){
		p = (LPTSTR)HeapAlloc(GetProcessHeap(), 0, s);
			if(p==NULL)
				return 0;
	}
	else
	{//s==0 就是在文件末尾
		WriteFile(this->fh,lpBuffer,nNumberOfBytesToWrite,&lpNumberOfBytesWritten,NULL);
		return lpNumberOfBytesWritten;
	}
	
	ReadFile(this->fh,p,s,NULL,NULL);//保存
	
	SetFilePointer(this->fh,pointer,NULL,FILE_BEGIN);
	WriteFile(this->fh,lpBuffer,nNumberOfBytesToWrite,&lpNumberOfBytesWritten,NULL);//写入
	
	WriteFile(this->fh,p,s,NULL,NULL);
	HeapFree(GetProcessHeap(), 0, p);
	
	SetFilePointer(this->fh,pointer+nNumberOfBytesToWrite,NULL,FILE_BEGIN);
	
	return lpNumberOfBytesWritten;
}

DWORD RWF::del(DWORD dwLength)
{
	DWORD filesize=this->getFileSize();
	DWORD pointer=this->getPointer();
	
	if(pointer+dwLength>=filesize){
		SetEndOfFile(this->fh);
		return filesize-pointer;
	}
	
	//s>0
	DWORD s=filesize-SetFilePointer(this->fh,dwLength,NULL,FILE_CURRENT);
	DWORD lpNumberOfBytesWritten;
	

	LPTSTR p = (LPTSTR)HeapAlloc(GetProcessHeap(), 0, s);
	if(p==NULL)
	{
		SetFilePointer(this->fh,pointer,NULL,FILE_BEGIN);
		return 0;
	}
	
	ReadFile(this->fh,p,s,NULL,NULL);//保存
	
	SetFilePointer(this->fh,pointer,NULL,FILE_BEGIN);
	WriteFile(this->fh,p,s,&lpNumberOfBytesWritten,NULL);//写入
	HeapFree(GetProcessHeap(), 0, p);
	
	SetFilePointer(this->fh,filesize-dwLength,NULL,FILE_BEGIN);
	SetEndOfFile(this->fh);
	
	SetFilePointer(this->fh,pointer,NULL,FILE_BEGIN);
	
	return dwLength;
}

BOOL RWF::flushFileBuffers()
{
	return FlushFileBuffers(this->fh);
}

DWORD RWF::setPointer(
			LONG lDistanceToMove,
			DWORD dwMoveMethod
			)
{
	DWORD pointer=SetFilePointer(this->fh,lDistanceToMove,NULL,dwMoveMethod);
	//cout<<endl<<"mark:"<<pointerB<<endl;
	if(pointer > this->getFileSize())
		return SetFilePointer(this->fh,0,NULL,FILE_END);
	return pointer;
}

DWORD RWF::getPointer()
{
	return SetFilePointer(this->fh,0,NULL,FILE_CURRENT);
}

DWORD RWF::getFileSize()
{
	return GetFileSize(this->fh,NULL);
}