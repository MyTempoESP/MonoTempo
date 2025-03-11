#include "config/config.h"
#include "reader/reader.hpp"

#include "chafon_reader.hpp"

int main(void) {
  Reader reader;

  INFO("Opening connection to Reader");

  DevicePara param;

  DEFER(reader.openAndFetchParams(&param), quit);

  INFO("Reader opened");
  INFO("Starting Inventory");

  DEFER(reader.startInventory(), quit);

  INFO("Starting monitor...");
  reader.createMonitor();

  INFO("Joining monitor...");
  reader.joinMonitor();

  return 0;

quit:
  return 1;
}
