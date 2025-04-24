/**
 * @file aa2.ino
 * @brief This file implements a system for managing and displaying RFID reader data,
 * network status, and system information on an LCD screen.
 * It also handles user input through buttons and serial communication.
 *
 * @details
 * The system is designed to:
 * - Display various system statuses and data on an LCD screen.
 * - Handle button inputs to navigate between screens and trigger actions.
 * - Parse and process data received via serial communication.
 * - Lock and unlock the screen for specific operations.
 * - Provide confirmation prompts for critical actions.
 * - Manage power-off countdown and sleep mode.
 *
 * @author Rodrigo Monteiro Junior
 * @date 2025-03-04
 * @version 2.0
 *
 * @copyright (c) 2025 Rodrigo Monteiro Junior
 *
 * @dependencies
 * - SafeString.h: For safe string handling.
 * - SafeStringReader.h: For reading and parsing serial input.
 * - BufferedOutput.h: For buffered serial output.
 * - avr/sleep.h: For sleep mode functionality.
 * - LiquidCrystal_I2C.h: For controlling the LCD screen.
 * - Standard C++ libraries: inttypes.h, time.h.
 *
 * @hardware
 * - 20x4 LCD screen with I2C interface.
 * - Two buttons connected to digital pins (BUTTON_VANCE and BUTTON_START).
 * - RFID reader for tag data.
 *
 * @sections
 * - Data Structures:
 *   - PCTagData: A struct to hold RFID tag data such as counts and antenna statuses.
 *   - PCData: A struct to hold system data such as tag counts, statuses, version info, and timestamps.
 * - Global Variables:
 *   - g_system_data: Holds the current system data.
 *   - g_current_screen: Tracks the currently displayed screen.
 *   - g_locked: Indicates whether the screen is locked.
 *   - g_screen_waiting_confirmation: Indicates if the system is waiting for user confirmation.
 *   - g_does_antenna_reports: Indicates if antenna reports are enabled.
 * - Functions:
 *   - check_clicked(): Checks if a button is clicked and returns the button ID.
 *   - parse_pc_data(): Parses serial input to update system data.
 *   - parse_time(): Parses and converts a timestamp to a human-readable date and time.
 *   - check_sum(): Validates the checksum of a received serial message.
 *   - screen_build(): Builds the virtual screen content based on the current screen.
 *   - screen_draw(): Draws the virtual screen content on the LCD.
 *   - screen_lock(): Locks the screen for specific operations.
 *   - screen_unlock(): Unlocks the screen after operations.
 *   - screen_next(): Navigates to the next screen.
 *   - screen_confirm(): Prepares the system for a confirmation action.
 *   - screen_wait_confirm(): Handles the confirmation waiting state.
 *   - screen_poweroff_countdown(): Handles the power-off countdown and puts the system into sleep mode.
 *   - event_send(): Sends the current screen's event via serial communication.
 *   - handle_serial(): Handles serial input and updates the system state.
 *   - handle_buttons(): Handles button inputs for navigation and actions.
 *   - setup(): Initializes the system, LCD, and serial communication.
 *   - loop(): Main program loop to handle tasks and update the screen.
 *
 * @screens
 * - Informational screens:
 *   - Display tag counts, network statuses, system version, and more.
 * - Action screens:
 *   - Allow actions like uploading data, resetting, or shutting down.
 * - Confirmation screens:
 *   - Prompt the user for confirmation before critical actions.
 * - Power-off screen:
 *   - Displays a countdown before shutting down the system.
 *
 * @constants
 * - BUTTON_VANCE, BUTTON_START: Digital pins for buttons.
 * - START_DELIMITER, END_DELIMITER: Delimiters for serial communication.
 * - NAV_SCREENS_COUNT, SCREENS_COUNT: Number of navigable and total screens.
 * - VIRT_SCR_COLS, VIRT_SCR_ROWS: Dimensions of the virtual screen buffer.
 *
 * @usage
 * - Use the buttons to navigate between screens and trigger actions.
 * - Monitor the LCD for system statuses and prompts.
 * - Send data via serial communication in the expected format for updates.
 *
 * @changes
 * - Added `parse_time()` function to handle timestamp parsing and conversion.
 * - Added `check_sum()` function to validate the checksum of serial messages.
 * - Updated `parse_pc_data()` to include parsing of additional fields like timestamps and antenna data.
 * - Added `handle_serial()` to process serial input and update system state.
 * - Added `handle_buttons()` to handle button inputs for navigation and actions.
 * - Updated `screen_build()` to include new screens and data fields.
 * - Added `screen_wait_confirm()` to handle confirmation waiting logic.
 * - Added `screen_poweroff_countdown()` to manage power-off countdown and sleep mode.
 */

#include <SafeString.h>
#include <SafeStringReader.h>
#include <BufferedOutput.h>

#include <avr/sleep.h>

#include <time.h>

#define __STDC_FORMAT_MACROS
#include <inttypes.h>

#define BUTTON_VANCE 6
#define BUTTON_START 7

bool g_clicked_vance = false;
bool g_clicked_start = false;

int check_clicked(void)
{
	if (digitalRead(BUTTON_VANCE) == 0)
	{
		if (!g_clicked_vance)
		{
			g_clicked_vance = true;
			return BUTTON_VANCE;
		}
	}
	else
		g_clicked_vance = false;

	if (digitalRead(BUTTON_START) == 0)
	{
		if (!g_clicked_start)
		{
			g_clicked_start = true;
			return BUTTON_START;
		}
	}
	else
		g_clicked_start = false;

	return 0;
}

/*
| **Data**             | **Type** | **Description**                                            | **Target screen** | **Format**                                                                                                                       |
|----------------------|----------|------------------------------------------------------------|-------------------|----------------------------------------------------------------------------------------------------------------------------------|
| Tags                 | Int32    | Number of tags read by the RFID reader.                    | 1                 | %dK (K only included after 10K tags) |
| Unique-Tags          | Int32    | Number of Unique tags read by the RFID reader.             | 1                 | %d                                                                                                                               |
| Communication status | Bool     | Status of the connection between the PC and mytempo.esp.br | 2                 | SIM(True) / NÃO(False)                                                                                                           |
| WI-FI status         | Bool     | Status of the general internet connection of the device.   | 3                 | OK(True) / X(false)                                                                                                              |
| 4G status            | Bool     | Status of the LTE/4G connection of the device.             | 3                 | OK(True) / X(false)                                                                                                              |
| Reader status        | Bool     | Status of the RFID reader.                                 | 4                 | OK(True) / X(false)                                                                                                              |
| System version       | Int32    | Version number of the system.                              | 5                 | %d                                                                                                                               |
| Backup count         | Int32    | Number of backups currently stored.                        | 6                 | %d                                                                                                                               |
| Envio count          | Int32    | Number of envios currently stored.                         | 7                 | %d                                                                                                                               |
*/

typedef struct PCTagData
{
	int32_t tags;
	int unique_tags;
	int32_t antenna1;
	int32_t antenna2;
	int32_t antenna3;
	int32_t antenna4;
} PCTagData;

typedef struct PCData
{
	PCTagData tag_data;

	bool comm_status;
	bool rfid_status;
	bool usb_status;
	int sys_version;
	int num_serie;
	int backups;
	int permanent_unique_tags;

	// DateTime
	int year;
	int month;
	int day;
	int hour;
	int minute;
	int second;
} PCData;

typedef struct LOGData
{
	int uploadcount;
	int databases;
	int avgproctime;
	int errcount;
	bool status;
} LOGData;

LOGData g_system_logs;
PCData g_system_data;

// whether or not the antenna reports are enabled
bool g_does_antenna_reports = false;

// whether or not the log reports are enabled
bool g_does_log_reports = false;

// create a SafeString reader to read the struct data
createSafeStringReader(serial_reader, 120, '\n', true);

// create a BufferedOutput to write screen data out
createBufferedOutput(serial_writer, 12, BLOCK_IF_FULL, true);

// Perform basic XOR checksum8 to input, excluding the starting '$' sign.
bool check_sum(SafeString &msg)
{
	int idx_star = msg.indexOf('*');

	cSF(check_sum_hex, 2);

	msg.substring(check_sum_hex, idx_star + 1);

	long sum = 0;

	if (!check_sum_hex.hexToLong(sum))
	{
		return false;
	}

	// skip the first character '$'
	for (size_t i = 1; i < idx_star; i++)
	{
		sum ^= msg[i];
	}

	return (sum == 0);
}

// Parses a UNIX timestamp and correctly offsets it to arduino's standards
// (arduino starts at 2000-01-01 rather than 1970)
bool parse_time(SafeString &timeField)
{
	int64_t stime = 0;

	if (!timeField.toInt64_t(stime))
	{
		return false;
	}

	if (stime < 0)
	{
		return false;
	}

	if (stime > UINT32_MAX)
	{
		return false;
	}

	if (stime < UNIX_OFFSET) // 2000-01-01 00:00:00
	{
		return false;
	}

	stime -= UNIX_OFFSET;

	time_t time = static_cast<time_t>(stime);

	struct tm *tm_ptr = gmtime(&time);
	if (!tm_ptr)
		return false;

	g_system_data.year = tm_ptr->tm_year + 1900;
	g_system_data.month = tm_ptr->tm_mon + 1;
	g_system_data.day = tm_ptr->tm_mday;
	g_system_data.hour = tm_ptr->tm_hour;
	g_system_data.minute = tm_ptr->tm_min;
	g_system_data.second = tm_ptr->tm_sec;

	return true;
}

// XXX: gets an int64_t and casts it to a int32_t
// the result is stored in &dest only if the conversion succeeds.
// @param field The field to convert
// @return true if the conversion was successful, false otherwise
bool getInt32Field(SafeString &field, int32_t &dest)
{
	int64_t value = 0;

	if (!field.toInt64_t(value))
		return false;

	if (value <= INT32_MIN || value >= INT32_MAX)
		return false;

	dest = static_cast<int32_t>(value);

	return true;
}

/*
## Serial Communication Input Format

The program now supports two types of serial messages: Antenna Data (A) and PCData (P). The format is as follows:

```perl
$MYTMP;<tags>;<unique_tags>;<type>;<data>;<timestamp>*<checksum>
```

### Field Descriptions

- `<tags>`: Total number of tags read (integer).
- `<unique_tags>`: Number of unique tags read (integer).
- `<type>`: Message type (`A` for Antenna Data, `P` for PCData).
- `<data>`: Data fields specific to the message type (see below).
- `<timestamp>`: The UNIX timestamp for updating time.
- `<checksum>`: XOR checksum of the message (excluding the `$`).

#### Antenna's `<data>` Fields (`A`)

- `<antenna1>`: Count of tags read by antenna 1 (integer).
- `<antenna2>`: Count of tags read by antenna 2 (integer).
- `<antenna3>`: Count of tags read by antenna 3 (integer).
- `<antenna4>`: Count of tags read by antenna 4 (integer).

#### PCData's `<data>` Fields (`P`)

- `<comm_status>`: Communication status (`1` for true, `0` for false).
- `<rfid_status>`: RFID reader status (`1` for true, `0` for false).
- `<usb_status>`: USB connection status (`1` for true, `0` for false).
- `<sys_version>`: System version number (integer).
- `<num_serie>`: Serial number of the system (integer).
- `<backups>`: Number of backups stored (integer).
- `<permanent_unique_tags>`: Number of permanent unique tags (integer).

### Examples

#### Antenna Data Example

```perl
$MYTMP;12345;678;A;100;200;300;400;1744846380*<checksum>
```

- `<tags>`: 12345
- `<unique_tags>`: 678
- `<type>`: `A` (Antenna Data)
- `<antenna1>`: 100
- `<antenna2>`: 200
- `<antenna3>`: 300
- `<antenna4>`: 400
- `<timestamp>`: 1744846380 (Wed, 16 Apr 2025 23:33:00 GMT)

#### PCData Example

```perl
$MYTMP;12345;678;P;1;1;0;42;123;5;100;1744846380*<checksum>
```

- `<tags>`: 12345
- `<unique_tags>`: 678
- `<type>`: `P` (PCData)
- `<comm_status>`: 1 (true)
- `<rfid_status>`: 1 (true)
- `<usb_status>`: 0 (false)
- `<sys_version>`: 42
- `<num_serie>`: 123
- `<backups>`: 5
- `<permanent_unique_tags>`100
- `<timestamp>`: 1744846380 (Wed, 16 Apr 2025 23:33:00 GMT)
*/
bool parse_pc_data(SafeString &msg)
{
	cSF(field, 11);

	char delims[] = ";*";
	bool returnEmptyFields = true;

	int idx = 0;

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (field != "$MYTMP")
	{
		return false;
	}

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!getInt32Field(field, g_system_data.tag_data.tags))
		return false;

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!field.toInt(g_system_data.tag_data.unique_tags))
		return false;

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (field.equals("T"))
	{
		// do nothing, just short-circuit
	} // do antenna update
	else if (field.equals("A"))
	{
		// enable antenna reports if one is actually received
		g_does_antenna_reports = true;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!getInt32Field(field, g_system_data.tag_data.antenna1))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!getInt32Field(field, g_system_data.tag_data.antenna2))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!getInt32Field(field, g_system_data.tag_data.antenna3))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!getInt32Field(field, g_system_data.tag_data.antenna4))
			return false;
	} // do PCData update
	else if (field.equals("P"))
	{
		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		g_system_data.comm_status = field.equals("1");

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		g_system_data.rfid_status = field.equals("1");

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		g_system_data.usb_status = field.equals("1");

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!field.toInt(g_system_data.sys_version))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!field.toInt(g_system_data.num_serie))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!field.toInt(g_system_data.backups))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!field.toInt(g_system_data.permanent_unique_tags))
			return false;
	} // do LOGData update
	else if (field.equals("L"))
	{
		g_does_log_reports = true;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!field.toInt(g_system_logs.uploadcount))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!field.toInt(g_system_logs.databases))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!field.toInt(g_system_logs.avgproctime))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		if (!field.toInt(g_system_logs.errcount))
			return false;

		idx = msg.stoken(field, idx, delims, returnEmptyFields);

		g_system_logs.status = field.equals("1");
	}

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!parse_time(field))
		return false;

	return true;
}

/* SCREEN_H */
#include <LiquidCrystal_I2C.h>
#include <string.h>

#define VIRT_SCR_COLS 20
#define VIRT_SCR_ROWS 4

char g_virt_scr[VIRT_SCR_ROWS][VIRT_SCR_COLS + 1];

LiquidCrystal_I2C lcd(0x27, VIRT_SCR_COLS, VIRT_SCR_ROWS);

const char fill_pattern[20] = "                   ";

#define virt_scr_sprintf(x, y, fmt, ...) \
	snprintf(g_virt_scr[y] + x, (VIRT_SCR_COLS - x), fmt, __VA_ARGS__);
#define virt_scr_fill_from(n, y) \
	(n < 20 && snprintf(g_virt_scr[y] + n, (VIRT_SCR_COLS - n), fill_pattern + n));

/*
| **Screen**              | **Content**                                      | **Description**                                                       | **Action (optional)**                                                                       |
|-------------------------|--------------------------------------------------|-----------------------------------------------------------------------|---------------------------------------------------------------------------------------------|
| Info                    | - Regist.: (Tags)<br>- Atletas: (Unique-Tags)    | Displays the<br>tag count, and<br>unique tag count.                   | START: Reset visual tag data<br>information, does not<br>actually touch anything<br>stored. |
| Network                 | - Comunicando: (Communication status)            | Displays the communication<br>with mytempo.esp.br                     |                                                                                             |
| Network Mgmt            | - Wi-Fi: (WI-FI status)<br>- LTE/4G: (4G status) | Displays basic PC network<br>connectivity info.                       | START: Issue a reconnection<br>of both wifi and 4g networks.                                |
| System                  | - Version: (System version)                      | Displays the system version,<br>i.e. the current update.              | START: Fetch and Install the<br>latest version from github.                                 |
| Upload                  | - Regist.: (Tags)<br>- Pendentes: (Envio count)  | Displays the current tag<br>count + the number of pending<br>uploads. | START: Upload all tag data<br>currently stored + pending<br>tag data.                       |
| Upload (Backup)         | - Backups: (Backup count)                        | Displays the number of backups.                                       | START: Upload all backups.                                                                  |
| #15 (Erase data)        |                         -                        |                                   -                                   |                                              -                                              |
| #15 (Shutdown)          |                         -                        |                                   -                                   |                                              -                                              |

| #15 (Shutdown [Helper]) |                         -                        |                                   -                                   |                                              -                                              |
| Confirmation screen     | - Pressione START para confirmar...              | Waits for user confirmation before an action                          | START: Confirmation                                                                         |
*/
#define INFORM_SCREEN 0
#define ANTNNA_SCREEN 1
#define NETWRK_SCREEN 2
#define NETCFG_SCREEN 3
#define USBCFG_SCREEN 4
#define DATTME_SCREEN 5
#define SYSTEM_SCREEN 6
#define UPLOAD_SCREEN 7
#define BACKUP_SCREEN 8
#define DELETE_SCREEN 9
#define SHTDWN_SCREEN 10
#define LOGRPT_SCREEN 11
#define NAV_SCREENS_COUNT 12

#define OFFMSG_SCREEN 13
#define CONFRM_SCREEN 14
#define WAITNG_SCREEN 15
#define WAITON_SCREEN 16
#define SCREENS_COUNT 17

unsigned int g_current_screen = 0;
unsigned int g_confirm_target = 0; // target screen for events that need confirmation
unsigned int g_eta = 30;	   // countdown for shutdown message

int g_unlocks;
bool g_locked;
bool g_screen_waiting_confirmation;
int32_t g_screen_waiting_timestamp;

const char desc[SCREENS_COUNT][VIRT_SCR_COLS] = {
    "START:Reset tela   ",
    "                   ",
    "                   ",
    "START:Reconectar   ",
    "START:Salvar no USB",
    "                   ",
    "START:Atualizar    ",
    "START:Upload Regist",
    "START:Upload Backup",
    "START:Apagar tudo  ",
    "START:Desligar     ",
    "                   ",
    "START:Confirma     ",
    "                   ",
    "                   ",
    "                   "};

void screen_build(void)
{
	unsigned int l0 = 0, l1 = 0, l2 = 0;

	switch (g_current_screen)
	{
	case INFORM_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Regist.: %" PRId32, g_system_data.tag_data.tags);
		l2 = virt_scr_sprintf(0, 2, "Atletas: %"
					    "d",
				      g_system_data.tag_data.unique_tags);
		break;
	case ANTNNA_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "A1: %" PRId32 " A2: %" PRId32,
				      g_system_data.tag_data.antenna1, g_system_data.tag_data.antenna2);

		l2 = virt_scr_sprintf(0, 2, "A3: %" PRId32 " A4: %" PRId32,
				      g_system_data.tag_data.antenna3, g_system_data.tag_data.antenna4);
		break;
	case NETWRK_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Leitor     : %2s", g_system_data.rfid_status ? "OK" : "X");
		l2 = virt_scr_sprintf(0, 2, "Comunicando: %3s", g_system_data.comm_status ? "SIM" : "NAO");
		break;
	case NETCFG_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Reconectar a rede", NULL);
		l2 = virt_scr_sprintf(0, 2, "Wi-Fi ou 4G", NULL);
		break;
	case USBCFG_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "USB: %2s", g_system_data.usb_status ? "OK" : "X");
		break;
	case DATTME_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Data: %02d/%02d/%04d", g_system_data.day, g_system_data.month, g_system_data.year);
		l2 = virt_scr_sprintf(0, 2, "Hora: %02d:%02d:%02d", g_system_data.hour, g_system_data.minute, g_system_data.second);
		break;
	case SYSTEM_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Versao: %"
					    "d",
				      g_system_data.sys_version);
		break;

		/* down here are screens with no Heading, so they can use l0 */

	case UPLOAD_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Atletas: %"
					    "d",
				      g_system_data.permanent_unique_tags);
		// l2 = virt_scr_sprintf(0, 2, "Pendentes: %" "d", g_system_data.envios);
		break;
	case BACKUP_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Backups: %"
					    "d",
				      g_system_data.backups);
		break;
	case DELETE_SCREEN:
		l0 = virt_scr_sprintf(0, 0, "Apagar dados", NULL);
		l1 = virt_scr_sprintf(0, 1, "do equipamento", NULL);
		break;
	case SHTDWN_SCREEN:
		l0 = virt_scr_sprintf(0, 0, "Desligar o", NULL);
		l1 = virt_scr_sprintf(0, 1, "equipamento", NULL);
		break;
	case LOGRPT_SCREEN:
		l0 = virt_scr_sprintf(0, 0, "subiu: %d", g_system_logs.uploadcount);
		l1 = virt_scr_sprintf(0, 1, "dur.med. %dcs", g_system_logs.avgproctime);
		l2 = virt_scr_sprintf(0, 2, "erros: %d", g_system_logs.errcount);
		break;

		/* end of NAV_SCREENS */
		/* Extra screens */

	case OFFMSG_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Aguarde %d segundos", g_eta);
		l2 = virt_scr_sprintf(0, 2, "E pressione POWER", NULL);
		break;
	case CONFRM_SCREEN:
		l0 = virt_scr_sprintf(0, 0, "Pressione START", NULL);
		l1 = virt_scr_sprintf(0, 1, "para confirmar", NULL);
		break;
	case WAITNG_SCREEN:
		l2 = virt_scr_sprintf(0, 2, "Aguarde...", NULL);
		break;
	case WAITON_SCREEN:
		l2 = virt_scr_sprintf(0, 2, "Inicializando...", NULL);
		break;
	}

	if (g_current_screen < UPLOAD_SCREEN)
	{
		l0 = virt_scr_sprintf(0, 0, "PORTAL my%d", g_system_data.num_serie);
	}

	virt_scr_fill_from(l0, 0);
	virt_scr_fill_from(l1, 1);
	virt_scr_fill_from(l2, 2);

	virt_scr_sprintf(0, 3, "%s", desc[g_current_screen]);
}

void screen_poweroff_countdown(void)
{
	g_current_screen = OFFMSG_SCREEN;

	while (g_eta-- > 0)
	{
		screen_draw();
		delay(1000);
	}

	// turn off the screen
	lcd.noBacklight();
	lcd.clear();

	set_sleep_mode(SLEEP_MODE_PWR_DOWN);
	sleep_enable();
	sleep_cpu();
}

void screen_unlock(void)
{
	// we require three unlocks to unlock the screen
	// this is to prevent accidental unlocks when the screen is locked
	if (g_unlocks++ < 3)
		return;

	g_unlocks = 0; // reset the locks

	if (!g_locked)
		return;

	g_current_screen = INFORM_SCREEN;
	g_locked = false;
}

// Set screen to WAITNG_SCREEN, lock navigation.
// A successful Serial receive triggers an unlock.
void screen_lock(void)
{
	if (g_current_screen == SHTDWN_SCREEN)
	{
		screen_poweroff_countdown();
	}

	g_current_screen = WAITNG_SCREEN;
	g_locked = true;

	serial_reader.flushInput();
}

// Set screen to WAITNG_SCREEN, lock navigation.
// A successful Serial receive triggers an unlock.
// @param screen The screen to lock to.
void screen_lock(int screen)
{
	g_current_screen = screen;
	g_locked = true;

	serial_reader.flushInput();
}

void screen_next(void)
{
	g_current_screen = (g_current_screen + 1) % NAV_SCREENS_COUNT;

	// skip the antenna screen if antenna reports are disabled
	if (g_current_screen == ANTNNA_SCREEN && !g_does_antenna_reports)
	{
		g_current_screen = NETWRK_SCREEN;
	}

	// skip the log report screen if log reports are disabled
	if (g_current_screen == LOGRPT_SCREEN && !g_does_log_reports)
	{
		g_current_screen = INFORM_SCREEN;
	}
}

void screen_confirm(void)
{
	g_screen_waiting_timestamp = millis();
	g_screen_waiting_confirmation = true;
	g_confirm_target = g_current_screen;
	g_current_screen = CONFRM_SCREEN;
}

void screen_wait_confirm(void)
{
	if (check_clicked())
	{
		g_screen_waiting_confirmation = false;
		event_send();
	}

	if (millis() - g_screen_waiting_timestamp > 2000)
	{
		g_screen_waiting_confirmation = false;
		g_current_screen = g_confirm_target;
	}
}

void screen_draw(void)
{
	screen_build();

	for (int i = 0; i < VIRT_SCR_ROWS; i++)
	{
		lcd.setCursor(0, i);
		for (char *c = g_virt_scr[i], i = 0; *c != '\0' && i < VIRT_SCR_COLS; c++, i++)
			lcd.write(*c);
	}
}

void screen_init(void)
{
	lcd.init();	 // Initialize the LCD
	lcd.backlight(); // Turn on the backlight
	memset(g_virt_scr, '\0', sizeof(g_virt_scr));
}

void event_send(void)
{
	// screens that need confirmation
	if (g_current_screen == DELETE_SCREEN || g_current_screen == SHTDWN_SCREEN)
	{
		screen_confirm();
		return;
	}

	if (g_current_screen == CONFRM_SCREEN)
		g_current_screen = g_confirm_target;

	char buf[11];
	snprintf(buf, 11, "$MYTMP;%d", g_current_screen);
	serial_writer.write((uint8_t *)buf, 10);

	// filter out screens which have non-blocking actions or no action at all
	if (g_current_screen > NETWRK_SCREEN && g_current_screen != DATTME_SCREEN)
		screen_lock();
}

void handle_serial(void)
{
	// serial_reader.skipToDelimiter();
	if (!serial_reader.read())
		return;

	serial_reader.trim();

	if (!check_sum(serial_reader))
		return;

	if (!serial_reader.startsWith("$MYTMP;"))
		return;

	if (parse_pc_data(serial_reader))
		screen_unlock();
}

void handle_buttons(void)
{
	// LOCKED
	if (g_locked)
		return;

	switch (check_clicked())
	{
	case BUTTON_VANCE:
		screen_next();
		break;
	case BUTTON_START:
		event_send();
		break;
	}
}

void setup(void)
{
	screen_init();

	Serial.begin(115200);
	while (!Serial)
		;

	serial_reader.connect(Serial); // where SafeStringReader will read from
	serial_writer.connect(Serial); // where BufferedOutput will write to

	pinMode(BUTTON_START, INPUT_PULLUP);
	pinMode(BUTTON_VANCE, INPUT_PULLUP);

	screen_lock(WAITON_SCREEN);
	screen_draw();
}

unsigned long previous_millis = 0;

void loop(void)
{
	handle_serial();

	// blink without delay
	unsigned long ms = millis();

	if (ms - previous_millis >= 5)
	{
		// if the screen is locked, skip the screen tasks
		// if the screen is waiting for confirmation, skip the screen tasks
		// if the screen is waiting for confirmation, skip the screen tasks

		previous_millis = ms;

		if (g_locked)
			return;

		if (g_screen_waiting_confirmation)
		{
			screen_wait_confirm();
		}
		else
			handle_buttons();

		screen_draw();
	}
}