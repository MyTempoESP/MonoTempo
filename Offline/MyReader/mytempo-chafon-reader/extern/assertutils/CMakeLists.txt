cmake_minimum_required(VERSION 3.15)

project(assertutils)

option(ASSERTUTILS_RELEASE "Less logging info and non verbosity" OFF)
configure_file(
	"${PROJECT_SOURCE_DIR}/include/assertutils/config/config.h.in"
	"${PROJECT_SOURCE_DIR}/include/assertutils/config/config.h"
)

add_subdirectory(src)

