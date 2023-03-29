#ifndef _PVM_WRAPPER_C
#define _PVM_WRAPPER_C

#include <stdio.h>
#include "pvm_wrapper.h"

int pvm_catchout_stdout()
{
    return pvm_catchout(stdout);
}

void *ptr_at(void **ptr, int idx)
{
    return ptr[idx];
}

// pvm_packf static bindings
int pvm_packf_string(char *fmt, char *arg)
{
    return pvm_packf(fmt, arg);
}

// pvm_unpackf static bindings
int pvm_unpackf_string(char *fmt, char *arg)
{
    return pvm_unpackf(fmt, arg);
}

#endif