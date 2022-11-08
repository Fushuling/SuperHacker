#pragma once
#include <Windows.h>

#define MyMemCpy(d,s,l) for(size_t mci=0;mci<(l);mci++){*((char*)(d)+mci)=*((char*)(s)+mci);}
size_t MyStrLenW(const wchar_t* s);
size_t MyStrLen(const char* s);
bool MyStrCmpW(const wchar_t* str1, const wchar_t* str2);
void MyStrAddW(wchar_t* d, const wchar_t* s);
void MyStrCpyW(wchar_t* d, const wchar_t* s);

#define MyAlloc(size) HeapAlloc(GetProcessHeap(), 0, (size_t)(size))
#define MyReAlloc(p,size) HeapReAlloc(GetProcessHeap(), 0, (void*)(p), (size_t)(size))
#define MyFree(p) HeapFree(GetProcessHeap(), 0, (void*)(p))
#define MyMemSize(p) HeapSize(GetProcessHeap(), 0, (void*)(p))//¿ÉÒÔ¿ÕÖ¸Õë

bool PipeCmd(wchar_t* pszCmd, wchar_t* CurrentDirectory, char* pszResultBuffer, DWORD dwResultBufferSize, bool wait);