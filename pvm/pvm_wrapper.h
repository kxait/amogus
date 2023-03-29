#ifndef _PVM_WRAPPER_H
#define _PVM_WRAPPER_H

#include <pvm3.h>

int pvm_catchout_stdout();

void *ptr_at(void **ptr, int idx);

int pvm_packf_string(char *fmt, char *arg);

int pvm_unpackf_string(char *fmt, char *arg);

#endif