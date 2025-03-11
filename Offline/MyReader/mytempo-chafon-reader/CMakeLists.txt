cmake_minimum_required(VERSION 3.15)

project(chafon_reader LANGUAGES C CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED True)

option(PRODUCTION "Build with less logging, use in conjuction with -DCMAKE_BUILD_TYPE=Release" OFF)
if (PRODUCTION)
	option(ASSERTUTILS_RELEASE "" ON)
endif()

configure_file(
	"${CMAKE_SOURCE_DIR}/include/chafon_reader/config/definitions.h.in"
	"${CMAKE_SOURCE_DIR}/include/chafon_reader/config/definitions.h"
)

set(CMAKE_MODULE_PATH "${CMAKE_MODULE_PATH};${CMAKE_CURRENT_SOURCE_DIR}/cmake")
include(submodules)

add_subdirectory(extern/assertutils)

option(BUILD_APP, "Build an app instead of a python module" OFF)

if (NOT BUILD_APP)
	set(PYTHON_BUILD 1)
endif()

add_subdirectory(src)

if (BUILD_APP)
	add_subdirectory(app)
else()
	add_subdirectory(python)
endif()

