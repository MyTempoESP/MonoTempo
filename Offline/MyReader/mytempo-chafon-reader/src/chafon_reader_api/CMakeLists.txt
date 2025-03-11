cmake_minimum_required(VERSION 3.15)

add_library(chafon_reader_api SHARED
	reader/reader.cpp
	tag/tag_processor.cpp
	tag/tag_monitor.cpp

	"${chafon_reader_SOURCE_DIR}/include/chafon_reader/reader/reader.hpp"
	"${chafon_reader_SOURCE_DIR}/include/chafon_reader/tag/tag_monitor.hpp"
	"${chafon_reader_SOURCE_DIR}/include/chafon_reader/tag/tag_processor.hpp"
)

# ourselves
target_include_directories(chafon_reader_api PUBLIC "${chafon_reader_SOURCE_DIR}/include/chafon_reader")

# including assertutils
target_include_directories(chafon_reader_api PUBLIC "${chafon_reader_SOURCE_DIR}/extern/assertutils/include/assertutils")
# linking assertutils
target_link_libraries(chafon_reader_api assertutils)

# including chafon
target_include_directories(chafon_reader_api PUBLIC "${chafon_reader_SOURCE_DIR}/extern/chafon/include")

# linking against proprietary chafon shared libraries
add_library(chafon_shared_api SHARED IMPORTED)
set_target_properties(chafon_shared_api
	PROPERTIES IMPORTED_LOCATION "${chafon_reader_SOURCE_DIR}/extern/chafon/lib/libCFApi.so"
)

add_library(chafon_shared_hid SHARED IMPORTED)
set_target_properties(chafon_shared_hid 
	PROPERTIES IMPORTED_LOCATION "${chafon_reader_SOURCE_DIR}/extern/chafon/lib/libhid.so"
)

target_link_libraries(chafon_reader_api chafon_shared_api chafon_shared_hid)

# install chafon proprietary libraries
install(IMPORTED_RUNTIME_ARTIFACTS chafon_shared_hid
	LIBRARY DESTINATION lib/chafon
)

install(IMPORTED_RUNTIME_ARTIFACTS chafon_shared_api
	LIBRARY DESTINATION lib/chafon
)

# Install api library
install(TARGETS chafon_reader_api
	LIBRARY DESTINATION lib/chafon_rfid
	PUBLIC_HEADER DESTINATION include/chafon_rfid
)

