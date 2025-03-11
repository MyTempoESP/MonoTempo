#include "error.h"

#include <stdio.h>
#include <stdarg.h>

// TODO: Centralize?
// general function f(color, msg, file, line) which handles each case

void log_ok(const char* msg, const char* fn, const char* file, int line) {
	printf("\033[32;1mGOOD\t\033[0m\033[1m%.10s\t:%d [\033[32;1m%s\033[0m\033[1m] => \033[32;1m'%s'\033[0m\n", file, line, fn, msg);
}

void log_status(const char* msg, const char* fn, const char* file, int line) {
	printf("\033[1mINFO\t\033[0m\033[1m%.10s\t:%d [\033[32;1m%s\033[0m\033[1m] => \033[36;1m'%s'\033[0m\n", file, line, fn, msg);
}

void log_okf(const char* fn, const char* file, int line, const char* msg, ...) {
	printf("\033[32;1mGOOD\t\033[0m\033[1m%.10s\t:%d [\033[32;1m%s\033[0m\033[1m] => \033[32;1m'", file, line, fn);

	va_list args;
	va_start(args, msg);
	vfprintf(stderr, msg, args);
	va_end(args);

	puts("\033[0m'");
}

void log_statusf(const char* fn, const char* file, int line, const char* msg, ...) {
	printf("\033[36;1mINFO\t\033[0m\033[1m%.10s\t:%d [\033[32;1m%s\033[0m\033[1m] => \033[36;1m'", file, line, fn);

	va_list args;
	va_start(args, msg);
	vprintf(msg, args);
	va_end(args);

	puts("\033[0m");

}

#ifndef ASSERTUTILS_RELEASE
void log_err(const char* msg, const char* fn, const char* file, int line) {
	fprintf(stderr, "\033[31;1mFAIL\t\033[0m\033[1m%.10s\t:%d [\033[32;1m%s\033[0m\033[1m] => \033[31;1m'%s'\033[0m\n", file, line, fn, msg);
}

void log_warn(const char* msg, const char* fn, const char* file, int line) {
	printf("\033[33;1mWARN\t\033[0m\033[1m%.10s\t:%d [\033[32;1m%s\033[0m\033[1m] => \033[33;1m'%s'\033[0m\n", file, line, fn, msg);
}

void log_errf(const char* fn, const char* file, int line, const char* msg, ...) {
	printf("\033[31;1mFAIL\t\033[0m\033[1m%.10s\t:%d [\033[32;1m%s\033[0m\033[1m] => \033[31;1m'", file, line, fn);

	va_list args;
	va_start(args, msg);
	vfprintf(stderr, msg, args);
	va_end(args);

	puts("\033[0m'");
}

void log_warnf(const char* fn, const char* file, int line, const char* msg, ...) {
	printf("\033[33;1mWARN\t\033[0m\033[1m%.10s\t:%d [\033[32;1m%s\033[0m\033[1m] => \033[33;1m'", file, line, fn);

	va_list args;
	va_start(args, msg);
	vprintf(msg, args);
	va_end(args);

	puts("\033[0m");

}
#endif

