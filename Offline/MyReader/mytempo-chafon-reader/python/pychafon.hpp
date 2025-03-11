#ifndef __CHAFON_RFID__H__
#define __CHAFON_RFID__H__

#include <Python.h>

#define PY_SSIZE_T_CLEAN

static PyObject* pychafon_fetch_tags(PyObject* self, PyObject* args);
static PyObject* pychafon_make_reader(PyObject* self, PyObject* args);

#ifdef __cplusplus
extern "C" {
#endif

PyMODINIT_FUNC PyInit_pychafon(void);

#ifdef __cplusplus
}
#endif

#endif
