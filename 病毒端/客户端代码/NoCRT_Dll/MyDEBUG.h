#pragma once

//#define UseMyDebug

#ifdef UseMyDebug
#include <stdarg.h>
#include <Windows.h>

typedef int (*nsf)(char* buf, const char* format, ...);
#define MyDebug(format,...)  {nsf sf=(nsf)GetProcAddress(GetModuleHandle(TEXT("ntdll.dll")), "sprintf");\
						char output[1024]={0};/*Max 65534*/\
						sf((char*)(output),(char*)(__FILE__##":%d"##" NOCRTDLL: "##format),__LINE__,__VA_ARGS__);\
						OutputDebugStringA(output);}
#else

#define MyDebug(format,...)

#endif