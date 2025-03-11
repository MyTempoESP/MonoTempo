#include "reader/reader.hpp"

#include <unistd.h>

Reader::Reader() {
	handle = 0;
}

Reader::~Reader() {
  WARN("DEBUG: Reader destructor called");

  if (handle != 0) {
    WARN("Reader open, closing...");

#ifndef __PYTHON_READER_MODULE__
    tm.stopMonitoring();
#endif

    stopInventory();
    close();
  }
}

#ifdef __PYTHON_READER_MODULE__
void Reader::read_n(int n) {
	INFOF("Reading %d tags", n);
	tm.monitorNTags(handle, n);
}
#else
void Reader::createMonitor(void) {
  tm.startMonitoring(handle);
}
void Reader::joinMonitor(void) {
  INFO("Creating and joining thread");
  // whatever
}
#endif

STATUS Reader::close(void) {
  ASSERT(CloseDevice(handle), "couldn't close device (it will close anyway)");
  return OK;
}

STATUS Reader::openAndFetchParams(DevicePara *param) {
  handle = 0;
  std::string ip = "192.168.1.200";

  char *cip = ip.data();

#ifndef PRODUCTION
  WARNF("DEBUG: reader ip = %s", cip);
#endif

  ASSERT(OpenNetConnection(&handle, cip, 2022, 3000),
         "Couldn't open net connection");

  WARNF( "Reader handle = %d", handle);

  INFO("Requesting device parameters from reader");

  ASSERT(GetDevicePara(handle, param), "GetDevicePara failed");

  return OK;
}

void Reader::SetFrequency() { // TODO: incomplete
  int region = 0;

  unsigned long result = 0x00;

  DevicePara param;

  GetDevicePara(handle, &param);
  param.REGION = region - 1;

  unsigned short inr;
  unsigned short dec;
  unsigned short step;

  inr = param.STRATFREI[0]; // INT PART
  dec = param.STRATFRED[0]; // DECIMAL PART

  // it's better to see 'char' as an int8 in this specific situation
  param.STRATFREI[0] = *(((char *)&inr) + 1);
  param.STRATFREI[1] = *((char *)&inr);
  param.STRATFRED[0] = *(((char *)&dec) + 1);
  param.STRATFRED[1] = *((char *)&dec);

  switch (region) {
  case 2: // USA
    step = 500;
    param.CN = 50;
    break;
  case 3: // Korea
    step = 200;
    param.CN = 15;
    break;
  case 4: // Europe
    step = 200;
    param.CN = 15;
    break;
  case 5: // Japan
    step = 200;
    param.CN = 8;
    break;
  case 6: // Malaysia
    step = 500;
    param.CN = 7;
    break;
  case 7: // Europe3
    step = 600;
    param.CN = 4;
    break;
  case 8: // China_1
    step = 250;
    param.CN = 20;
    break;
  case 9: // China_2
    step = 250;
    param.CN = 20;
    break;
  default:
    break;
  }

  param.STEPFRE[0] = *(((char *)&step) + 1);
  param.STEPFRE[1] = *((char *)&step);

  result = SetDevicePara(handle, param);
}

STATUS Reader::startInventory(void) {
  ASSERT(InventoryContinue(handle, 0, 0), "couldn't start inventory");
  return OK;
}

STATUS Reader::stopInventory(void) {
  ASSERT(InventoryStop(handle, 10000),
         "couldn't stop taking inventory (it will stop eventually)");
  return OK;
}
