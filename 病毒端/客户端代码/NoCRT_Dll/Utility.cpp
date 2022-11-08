#include "Utility.h"

size_t MyStrLenW(const wchar_t* s) {//长度不含终止符
	const wchar_t* o = s;
	while (*s) {
		s++;
	}
	return s - o;
}

size_t MyStrLen(const char* s) {//长度不含终止符
	const char* o = s;
	while (*s) {
		s++;
	}
	return s - o;
}

bool MyStrCmpW(const wchar_t* str1, const wchar_t* str2) {
	while (*str1 && *str2) {
		if (*str1 != *str2)
			return false;
		++str1; ++str2;
	}
	return *str1 == *str2;
}

void MyStrAddW(wchar_t* d, const wchar_t* s) {
	while (*d)
		d++;
	while (*d = *s) {
		d++; s++;
	}
}

void MyStrCpyW(wchar_t* d, const wchar_t* s) {
	while (*d = *s) {
		d++; s++;
	}
}

bool PipeCmd(wchar_t* pszCmd, wchar_t* CurrentDirectory, char* pszResultBuffer, DWORD dwResultBufferSize, bool wait){
	HANDLE hReadPipe = NULL;
	HANDLE hWritePipe = NULL;
	SECURITY_ATTRIBUTES securityAttributes = { 0 };
	BOOL bRet = FALSE;
	STARTUPINFOW si = { 0 };
	PROCESS_INFORMATION pi = { 0 };

	// 设定管道的安全属性
	securityAttributes.bInheritHandle = TRUE;
	securityAttributes.nLength = sizeof(securityAttributes);
	securityAttributes.lpSecurityDescriptor = NULL;
	// 创建匿名管道
	bRet = CreatePipe(&hReadPipe, &hWritePipe, &securityAttributes, 0);
	if (FALSE == bRet)
	{
		return FALSE;
	}
	// 设置新进程参数
	si.cb = sizeof(si);
	si.dwFlags = STARTF_USESHOWWINDOW | STARTF_USESTDHANDLES;
	si.wShowWindow = SW_HIDE;
	si.hStdError = hWritePipe;
	si.hStdOutput = hWritePipe;
	// 创建新进程执行命令, 将执行结果写入匿名管道中
	bRet = CreateProcessW(NULL, pszCmd, NULL, NULL, TRUE, 0, NULL, CurrentDirectory, &si, &pi);
	if (FALSE == bRet)
	{
		return false;
	}
	// 等待命令执行结束
	WaitForSingleObject(pi.hThread, INFINITE);
	WaitForSingleObject(pi.hProcess, INFINITE);
	// 从匿名管道中读取结果到输出缓冲区
	if(wait)
		ReadFile(hReadPipe, pszResultBuffer, dwResultBufferSize, NULL, NULL);
	// 关闭句柄, 释放内存
	CloseHandle(pi.hThread);
	CloseHandle(pi.hProcess);
	CloseHandle(hWritePipe);
	CloseHandle(hReadPipe);

	return TRUE;
}