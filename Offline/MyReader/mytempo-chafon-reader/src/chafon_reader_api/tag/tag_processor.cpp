#include "tag/tag_processor.hpp"

#ifdef __PYTHON_READER_MODULE__

static PyObject* tagCallback = NULL;

void setTagCallback(PyObject* c) {
	tagCallback = c;
}

void destroyTagCallback(void) {
	if (tagCallback == NULL) return;

	INFO("Found tag callback, destroying...");

	Py_XDECREF(tagCallback);

	tagCallback = NULL;
}

#endif

void processTag(TagInfo *tag) {
  char epc[255];

  int len = 0;
  for (int i = 0; i < 12; i++) {
    len += sprintf(epc + len, "%02X", tag->code[i]);
  }

#ifdef __PYTHON_READER_MODULE__
  if (tagCallback == NULL) return;

  PyObject* arglist;
  PyObject* result;

  //PyGILState_STATE state = PyGILState_Ensure();

  //arglist = Py_BuildValue("(s#)", epc, 24);
  arglist = Py_BuildValue("(s)", epc);
  result = PyObject_CallObject(tagCallback, arglist);

  //PyGILState_Release(state);

  Py_DECREF(arglist);

  if (result == NULL) {
	ERR( "Callback returned NULL!" );
  }

  Py_XDECREF(result);
#else
  WARNF("GOT EPC: %s", epc);
#endif
}
