#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "libhuzhoucrypt.h"
int main(){
    char key[6]="12345";
    char buf[200]="hello world";
    GoSlice value;
    GoSlice gokey;
    value.data=buf;
    value.cap=strlen(buf);
    value.len=strlen(buf);
    gokey.data=key;
    gokey.cap=strlen(key);
    gokey.len=strlen(key);

    struct HZEncrypt_return ret;
    ret = HZEncrypt(gokey,value,32);
    printf("%s\n",ret.r0.p);
    return 0;
}
