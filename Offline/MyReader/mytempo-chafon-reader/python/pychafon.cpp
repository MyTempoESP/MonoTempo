#include "config/config.h"
#include "reader/reader.hpp"
#include "tag/tag_processor.hpp" /* setTagCallback */

#include "pychafon.hpp"

Reader* reader;

static PyObject *
pychafon_make_reader(PyObject* self, PyObject* args) {
	INFO( "Creating new reader..." );

	if (reader != NULL) {
		PyErr_SetString(PyExc_RuntimeError, "Connection already open (try pychafon.close) ...");
		return NULL;
	}

	reader = new Reader();

	DevicePara param;

	DEFER(reader->openAndFetchParams(&param), quit);

	Py_RETURN_NONE;
quit:
	PyErr_SetString(PyExc_RuntimeError, "Couldn't open connection to reader");

	delete reader;

	reader = NULL;
	return NULL;
}

static PyObject *
pychafon_close_reader(PyObject* self, PyObject* args) {
	if (reader == NULL) {
		PyErr_SetString(PyExc_RuntimeError, "No readers currently open...");
		return NULL;
	}

	delete reader;

	reader = NULL;

	Py_RETURN_NONE;
}

static PyObject *
pychafon_set_tag_callback(PyObject* self, PyObject* args) {
	PyObject* result;
	PyObject* temp;

	if (PyArg_ParseTuple(args, "O:set_callback", &temp)) {
		if (!PyCallable_Check(temp)) {
			PyErr_SetString(PyExc_TypeError, "Parameter must be callable");
			return NULL;
		}

		Py_XINCREF(temp);

		destroyTagCallback();

		INFO("Setting tag callback...");
		setTagCallback(temp);

		Py_RETURN_NONE;
	}

	return result;
}

static PyObject *
pychafon_start_read(PyObject* self, PyObject* args) {
	if (reader == NULL) {
		PyErr_SetString(PyExc_RuntimeError, "No readers currently open...");
		return NULL;
	}

	int n = 0;

	if (!PyArg_ParseTuple(args, "i", &n))
		return NULL;

	int err = 0;

	Py_BEGIN_ALLOW_THREADS
	if (reader->startInventory() == FAIL) {
		err = 1;
	}
	Py_END_ALLOW_THREADS

	if (err == 0) {
		reader->read_n(n);
		goto quit;
	}


quit:
	Py_RETURN_NONE;
}

static PyMethodDef ChafonMethods[] = {
	{"read", pychafon_start_read, METH_VARARGS, "Fetch N tags from the rfid internal buffer"},
	{"open", pychafon_make_reader, METH_VARARGS, "Open connection to reader interface"},
	{"close", pychafon_close_reader, METH_VARARGS, "Close connection to reader"},
	{"setTagCallback", pychafon_set_tag_callback, METH_VARARGS, "Set a function to be called on tag read"},

	{NULL, NULL, 0, NULL}
};

static struct PyModuleDef pychafon = {
	PyModuleDef_HEAD_INIT,
	"pychafon",
	NULL,
	-1,

	ChafonMethods
};

#ifdef __cplusplus
extern "C" {
#endif

PyMODINIT_FUNC
PyInit_pychafon(void)
{
	return PyModule_Create(&pychafon);
}

#ifdef __cplusplus
}
#endif

