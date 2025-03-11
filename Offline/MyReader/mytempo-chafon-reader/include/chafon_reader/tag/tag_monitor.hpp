#ifndef __TAG_MONITOR__H__
#define __TAG_MONITOR__H__

#include "../config/config.h"
#include "CFApi.h"

#include "tag_processor.hpp"

#ifndef __PYTHON_READER_MODULE__
#include <atomic>
#include <thread>
#endif

class TagMonitor {
public:
  void monitorNTags(int reader_handle, int n = 5);
  void monitorTags(int reader_handle);

  ~TagMonitor();
  TagMonitor();

#ifndef __PYTHON_READER_MODULE__

  void startMonitoring(int reader_handle);
  void stopMonitoring(void);
  void waitTermination(void);

private:
  std::atomic<bool> monitoring;
  std::thread monitorThread;

  std::thread make_thread(int handle) {
    return std::thread([this, handle] { monitorTags(handle); });
  }
#endif
};

#endif
