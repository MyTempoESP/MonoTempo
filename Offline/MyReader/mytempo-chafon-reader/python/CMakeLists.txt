# python doesn't require lib before name

find_package(Python REQUIRED COMPONENTS Development)

if (Python_FOUND)
	Python_add_library(pychafon MODULE pychafon.cpp)

	target_include_directories(chafon_reader_api PUBLIC "${Python_INCLUDE_DIRS}")

	set_target_properties(pychafon PROPERTIES PREFIX "")

	target_link_libraries(pychafon PRIVATE chafon_reader_api)

	install(TARGETS pychafon
		LIBRARY DESTINATION lib/pychafon
	)
else()
	message(FATAL_ERROR "Python development libraries not found!")
endif()

#install(TARGETS chafon_rfid
#LIBRARY DESTINATION lib/chafon_python
#)

