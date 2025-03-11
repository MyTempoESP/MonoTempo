
#ifndef _HID_H_
#define _HID_H_

#include <cstddef>

struct hid_device_;
typedef struct hid_device_ hid_device; 

typedef enum {
		HID_API_BUS_UNKNOWN = 0x00,
		HID_API_BUS_USB = 0x01,
		HID_API_BUS_BLUETOOTH = 0x02,
		HID_API_BUS_I2C = 0x03,
		HID_API_BUS_SPI = 0x04,
	} hid_bus_type;
	
struct hid_device_info {
		char* path;
		unsigned short vendor_id;
		unsigned short product_id;
		wchar_t* serial_number;
		unsigned short release_number;
		wchar_t* manufacturer_string;
		wchar_t* product_string;
		unsigned short usage_page;
		unsigned short usage;
		int interface_number;
		struct hid_device_info* next;
		hid_bus_type bus_type;
	};
	
#ifdef __cplusplus
extern "C" {
#endif

int hid_init(void);

int hid_exit(void);

struct hid_device_info * hid_enumerate(unsigned short vendor_id, unsigned short product_id);

void hid_free_enumeration(struct hid_device_info* devs);

hid_device* hid_open(unsigned short vendor_id, unsigned short product_id, const wchar_t* serial_number);

hid_device* hid_open_path(const char* path);

int hid_write(hid_device* dev, const unsigned char* data, size_t length);

int hid_read_timeout(hid_device* dev, unsigned char* data, size_t length, int milliseconds);

int hid_read(hid_device* dev, unsigned char* data, size_t length);

int hid_set_nonblocking(hid_device* dev, int nonblock);

int hid_send_feature_report(hid_device* dev, const unsigned char* data, size_t length);

int hid_get_feature_report(hid_device* dev, unsigned char* data, size_t length);

int hid_get_input_report(hid_device* dev, unsigned char* data, size_t length);

void hid_close(hid_device* dev);

int hid_get_manufacturer_string(hid_device* dev, wchar_t* string, size_t maxlen);

int hid_get_product_string(hid_device* dev, wchar_t* string, size_t maxlen);

int hid_get_serial_number_string(hid_device* dev, wchar_t* string, size_t maxlen);

struct hid_device_info * hid_get_device_info(hid_device* dev);

int hid_get_indexed_string(hid_device* dev, int string_index, wchar_t* string, size_t maxlen);

int hid_get_report_descriptor(hid_device* dev, unsigned char* buf, size_t buf_size);

const wchar_t* hid_error(hid_device* dev);

const struct hid_api_version* hid_version(void);

const char* hid_version_str(void);

#ifdef __cplusplus
}
#endif

#endif

