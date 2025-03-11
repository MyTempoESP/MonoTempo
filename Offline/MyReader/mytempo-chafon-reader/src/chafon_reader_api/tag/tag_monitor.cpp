#include "tag/tag_monitor.hpp"

#define TAG_TIMEOUT 10000
constexpr int TAG_MAXTRIES = (TAG_TIMEOUT / 1000) * 3;

TagMonitor::TagMonitor() {
}

TagMonitor::~TagMonitor() {
	INFO("Destroying tag monitor...");

#ifdef __PYTHON_READER_MODULE__
	destroyTagCallback();
#else
	if (monitoring) {
		monitoring = false;
	}

	waitTermination();
#endif
}

#ifndef __PYTHON_READER_MODULE__
void TagMonitor::stopMonitoring(void) {
	monitoring = false;
}

void TagMonitor::startMonitoring(int reader_handle) {
	monitoring = true;

	monitorThread = make_thread(reader_handle);
}

void TagMonitor::waitTermination(void) {
	if (monitorThread.joinable() == true) {
		monitorThread.join();
	}

	return;
}
#endif

int GetTagUiiWrapper(int h, TagInfo* t, int timeout) {
  int ret = 0;
#ifdef PYTHON_BUILD
  Py_BEGIN_ALLOW_THREADS
#endif

  ret = GetTagUii(h, t, timeout);

#ifdef PYTHON_BUILD
  Py_END_ALLOW_THREADS
#endif
  return ret;
}

void /*_Noreturn*/ TagMonitor::monitorTags(int reader_handle) {
  TagInfo tag;

loop:
  while (GetTagUiiWrapper(reader_handle, &tag, TAG_TIMEOUT) != 0)
    ;

  processTag(&tag);

  goto loop;
}

void TagMonitor::monitorNTags(int reader_handle, int n) {
  if (n == 0) {
	monitorTags(reader_handle);
	return;
  }

  TagInfo tag;

  auto tries = 0;
  auto i = 0;

  tries = 0;
loop:
  while (GetTagUiiWrapper(reader_handle, &tag, TAG_TIMEOUT) != 0) {
    if (tries++ >= TAG_MAXTRIES)
      return;
  }

  Py_BEGIN_ALLOW_THREADS

  processTag(&tag);

  Py_END_ALLOW_THREADS

  if (i++ >= n)
    return;

  goto loop;
}

