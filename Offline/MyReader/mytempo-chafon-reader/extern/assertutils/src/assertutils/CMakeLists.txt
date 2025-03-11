add_library(assertutils SHARED
	error.c
	assertutils.c

	"${assertutils_SOURCE_DIR}/include/assertutils/error.h"
	"${assertutils_SOURCE_DIR}/include/assertutils/assertutils.h"
)

target_include_directories(assertutils PUBLIC "${assertutils_SOURCE_DIR}/include/assertutils")

install(TARGETS assertutils
	LIBRARY DESTINATION lib/assertutils
	PUBLIC_HEADER DESTINATION include/assertutils
)

