find_package(Python REQUIRED COMPONENTS Development)
if (NOT Python_FOUND)
	message(FATAL_ERROR "Python development libraries not found!")
endif()

