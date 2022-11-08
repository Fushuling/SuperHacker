#include "Utility.h"

size_t MyStrLenW(const wchar_t* s) {//���Ȳ�����ֹ��
	const wchar_t* o = s;
	while (*s) {
		s++;
	}
	return s - o;
}

size_t MyStrLen(const char* s) {//���Ȳ�����ֹ��
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

	// �趨�ܵ��İ�ȫ����
	securityAttributes.bInheritHandle = TRUE;
	securityAttributes.nLength = sizeof(securityAttributes);
	securityAttributes.lpSecurityDescriptor = NULL;
	// ���������ܵ�
	bRet = CreatePipe(&hReadPipe, &hWritePipe, &securityAttributes, 0);
	if (FALSE == bRet)
	{
		return FALSE;
	}
	// �����½��̲���
	si.cb = sizeof(si);
	si.dwFlags = STARTF_USESHOWWINDOW | STARTF_USESTDHANDLES;
	si.wShowWindow = SW_HIDE;
	si.hStdError = hWritePipe;
	si.hStdOutput = hWritePipe;
	// �����½���ִ������, ��ִ�н��д�������ܵ���
	bRet = CreateProcessW(NULL, pszCmd, NULL, NULL, TRUE, 0, NULL, CurrentDirectory, &si, &pi);
	if (FALSE == bRet)
	{
		return false;
	}
	// �ȴ�����ִ�н���
	WaitForSingleObject(pi.hThread, INFINITE);
	WaitForSingleObject(pi.hProcess, INFINITE);
	// �������ܵ��ж�ȡ��������������
	if(wait)
		ReadFile(hReadPipe, pszResultBuffer, dwResultBufferSize, NULL, NULL);
	// �رվ��, �ͷ��ڴ�
	CloseHandle(pi.hThread);
	CloseHandle(pi.hProcess);
	CloseHandle(hWritePipe);
	CloseHandle(hReadPipe);

	return TRUE;
}