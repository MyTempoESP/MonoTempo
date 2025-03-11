#ifndef __TAG_PROCESSOR__H__
#define __TAG_PROCESSOR__H__

#include "CFApi.h"
#include "../config/config.h"
//#include <cstdio>

void processTag(TagInfo *tag);

#ifdef __PYTHON_READER_MODULE__
#include <Python.h>
/* static PyObject * tagCallback; */
void setTagCallback(PyObject* c);
void destroyTagCallback(void);
#endif
#endif
