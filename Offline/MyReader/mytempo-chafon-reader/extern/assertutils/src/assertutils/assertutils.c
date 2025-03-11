#include "assertutils.h"

STATUS assertz(int result, const char* msg) {
	if (result != 0) {
		ERR( msg );
		return FAIL;
	}

	return OK;
}

STATUS assertz_quiet(int result) {
	if (result != 0) return FAIL;
	return OK;
}

