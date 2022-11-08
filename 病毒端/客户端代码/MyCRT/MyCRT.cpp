#include "MyCRT.h"

#pragma function(memset)
void* __cdecl memset(void* pTarget, int value, size_t cbTarget) {
    unsigned char* p = static_cast<unsigned char*>(pTarget);
    while (cbTarget-- > 0) {
        *p++ = static_cast<unsigned char>(value);
    }
    return pTarget;
}