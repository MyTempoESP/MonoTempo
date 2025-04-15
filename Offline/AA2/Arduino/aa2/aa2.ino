/**
 * @file aa2.ino
 * @brief This file contains the implementation of a system for managing and displaying
 *        RFID reader data, network status, and system information on an LCD screen.
 *        It also handles user input through buttons and serial communication.
 *
 * @details
 * The system is designed to:
 * - Display various system statuses and data on an LCD screen.
 * - Handle button inputs to navigate between screens and trigger actions.
 * - Parse and process data received via serial communication.
 * - Lock and unlock the screen for specific operations.
 * - Provide confirmation prompts for critical actions.
 *
 * @author Rodrigo Monteiro Junior
 * @date 2025-03-04
 *
 * @dependencies
 * - SafeString.h: For safe string handling.
 * - SafeStringReader.h: For reading and parsing serial input.
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
 *   - PCData: A struct to hold system data such as tag counts, statuses, and version info.
 * - Global Variables:
 *   - g_system_data: Holds the current system data.
 *   - g_current_screen: Tracks the currently displayed screen.
 *   - g_locked: Indicates whether the screen is locked.
 *   - g_screen_waiting_confirmation: Indicates if the system is waiting for user confirmation.
 * - Functions:
 *   - check_clicked(): Checks if a button is clicked and returns the button ID.
 *   - parse_data(): Parses serial input to update system data.
 *   - parse_time(): Parses and converts a timestamp to a human-readable date and time.
 *   - check_sum(): Validates the checksum of a received serial message.
 *   - screen_build(): Builds the virtual screen content based on the current screen.
 *   - screen_draw(): Draws the virtual screen content on the LCD.
 *   - screen_lock(): Locks the screen for specific operations.
 *   - screen_unlock(): Unlocks the screen after operations.
 *   - screen_next(): Navigates to the next screen.
 *   - screen_confirm(): Prepares the system for a confirmation action.
 *   - screen_wait_confirm(): Handles the confirmation waiting state.
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
 * - Updated `parse_data()` to include parsing of additional fields like timestamps.
 * - Added `handle_serial()` to process serial input and update system state.
 * - Added `handle_buttons()` to handle button inputs for navigation and actions.
 * - Updated `screen_build()` to include new screens and data fields.
 * - Added `screen_wait_confirm()` to handle confirmation waiting logic.
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

int check_clicked()
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
typedef struct PCData
{
	int64_t tags;
	int unique_tags;
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

#ifdef READER_HAS_MULTIPLE_ANTENNAS
int64_t g_antenna0 = 0;
int64_t g_antenna1 = 0;
int64_t g_antenna2 = 0;
int64_t g_antenna3 = 0;
#endif

PCData g_system_data;

int g_major_version = 1;

// create a SafeString reader to read the struct data
createSafeStringReader(serial_reader, 80, '\n', true);
createBufferedOutput(serial_writer, 12, BLOCK_IF_FULL, true);

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

	for (size_t i = 1; i < idx_star; i++)
	{
		sum ^= msg[i];
	}

	return (sum == 0);
}

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

	if (stime > 2147483647)
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

bool parse_data(SafeString &msg)
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

#ifdef READER_HAS_MULTIPLE_ANTENNAS
	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!field.toInt64_t(g_antenna0))
		return false;

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!field.toInt64_t(g_antenna1))
		return false;

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!field.toInt64_t(g_antenna2))
		return false;

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!field.toInt64_t(g_antenna3))
		return false;
#endif

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!field.toInt64_t(g_system_data.tags))
		return false;

	idx = msg.stoken(field, idx, delims, returnEmptyFields);

	if (!field.toInt(g_system_data.unique_tags))
		return false;

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
#ifdef READER_HAS_MULTIPLE_ANTENNAS
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
#define NAV_SCREENS_COUNT 11

#define OFFMSG_SCREEN 11
#define CONFRM_SCREEN 12
#define WAITNG_SCREEN 13
#define WAITON_SCREEN 14
#define SCREENS_COUNT 15
#else
#define INFORM_SCREEN 0
#define NETWRK_SCREEN 1
#define NETCFG_SCREEN 2
#define USBCFG_SCREEN 3
#define DATTME_SCREEN 4
#define SYSTEM_SCREEN 5
#define UPLOAD_SCREEN 6
#define BACKUP_SCREEN 7
#define DELETE_SCREEN 8
#define SHTDWN_SCREEN 9
#define NAV_SCREENS_COUNT 10

#define OFFMSG_SCREEN 10
#define CONFRM_SCREEN 11
#define WAITNG_SCREEN 12
#define WAITON_SCREEN 13
#define SCREENS_COUNT 14
#endif

unsigned int g_current_screen = 0;
unsigned int g_confirm_target = 0; // target screen for events that need confirmation
unsigned int g_eta = 30;	   // countdown for shutdown message

int g_unlocks;
bool g_locked;
bool g_screen_waiting_confirmation;
int32_t g_screen_waiting_timestamp;

const char desc[SCREENS_COUNT][VIRT_SCR_COLS] = {
    "START:Reset tela   ",
#ifdef READER_HAS_MULTIPLE_ANTENNAS
    "                   ",
#endif
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
};

void screen_build()
{
	unsigned int l0 = 0, l1 = 0, l2 = 0;

	switch (g_current_screen)
	{
	case INFORM_SCREEN:
		l1 = virt_scr_sprintf(0, 1, "Regist.: %" PRId32, g_system_data.tags);
		l2 = virt_scr_sprintf(0, 2, "Atletas: %"
					    "d",
				      g_system_data.unique_tags);
		break;
#ifdef READER_HAS_MULTIPLE_ANTENNAS
	case ANTNNA_SCREEN:

		int a[4] = {g_antenna0, g_antenna1, g_antenna2, g_antenna3};
		int postfix[4] = {' ', ' ', ' ', ' '};

		for (int i = 0; i < 4; i++)
		{
			if (a[i] > 1000)
			{
				a[i] /= 1000;
				postfix[i] = 'K';
			}
			else if (a[i] > 1000000)
			{
				a[i] /= 1000000;
				postfix[i] = 'M';
			}
		}

		l1 = virt_scr_sprintf(0, 1, "A0: %ld%c"
					    " A1: %ld%c",
				      a[0], postfix[0], a[1], postfix[1]);
		l2 = virt_scr_sprintf(0, 2, "A3: %ld%c"
					    " A4: %ld%c",
				      a[2], postfix[2], a[3], postfix[3]);
		break;
#endif
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

void screen_poweroff_countdown()
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

void screen_unlock()
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

void screen_lock()
{
	if (g_current_screen == SHTDWN_SCREEN)
	{
		screen_poweroff_countdown();
	}

	g_current_screen = WAITNG_SCREEN;
	g_locked = true;

	serial_reader.flushInput();
}

void screen_lock(int screen)
{
	g_current_screen = screen;
	g_locked = true;

	serial_reader.flushInput();
}

void screen_next()
{
	g_current_screen = (g_current_screen + 1) % NAV_SCREENS_COUNT;
}

void screen_confirm()
{
	g_screen_waiting_timestamp = millis();
	g_screen_waiting_confirmation = true;
	g_confirm_target = g_current_screen;
	g_current_screen = CONFRM_SCREEN;
}

void screen_wait_confirm()
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

void screen_draw()
{
	screen_build();

	for (int i = 0; i < VIRT_SCR_ROWS; i++)
	{
		lcd.setCursor(0, i);
		for (char *c = g_virt_scr[i], i = 0; *c != '\0' && i < VIRT_SCR_COLS; c++, i++)
			lcd.write(*c);
	}
}

void screen_init()
{
	lcd.init();	 // Initialize the LCD
	lcd.backlight(); // Turn on the backlight
	memset(g_virt_scr, '\0', sizeof(g_virt_scr));
}

#define START_DELIMITER 0x3C
#define END_DELIMITER 0x3E
void event_send()
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

	screen_lock();
}

void handle_serial()
{
	// serial_reader.skipToDelimiter();
	if (!serial_reader.read())
		return;

	serial_reader.trim();

	if (!check_sum(serial_reader))
		return;

	if (!serial_reader.startsWith("$MYTMP;"))
		return;

	if (parse_data(serial_reader))
		screen_unlock();
}

void handle_buttons()
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

void setup()
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

void loop()
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
