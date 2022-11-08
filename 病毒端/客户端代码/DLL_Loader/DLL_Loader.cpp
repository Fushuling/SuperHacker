#include <Windows.h>
using namespace std;
int main()
{
	LoadLibraryA("NoCRT_Dll.dll");
	system("pause");
}