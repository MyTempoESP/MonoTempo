#ifndef __ASSERTUTILS__H__
#define __ASSERTUTILS__H__

#include "error.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef enum STATUS {
	FAIL = 0,
	OK = 128
} STATUS;

STATUS assertz(int result, const char* msg);
STATUS assertz_quiet(int result);

// thank you fuckit.c =)
#define ASSERT(f, msg) do{\
				STATUS s = assertz_quiet(f);\
				if(s != OK) {\
					ERR(#f " failed!");\
					WARN("returning FAIL...");\
					return FAIL;\
				}\
			}while(0)

// takes a STATUS function
#define DEFER(f, label) do{\
				STATUS s = f;\
				if (s != OK) {\
					ERR(#f " got FAIL!");\
					WARN("deferring... " #label);\
					goto label;\
				}\
			}while(0)

#ifdef __cplusplus
}
#endif

#endif

