# Build Instructions

## Prerequisites

### 1. **Arduino IDE**: Install the Arduino IDE from [Arduino's official website](https://www.arduino.cc/en/software)

### 2. **Libraries**: Ensure the following libraries are installed in the Arduino IDE

- `SafeString`
- `LiquidCrystal_I2C`

### 3. **Hardware**

- Arduino Mega 2560 board.
- 20x4 LCD screen with I2C interface.
- Two buttons connected to digital pins (`BUTTON_VANCE` and `BUTTON_START`).
- RFID reader for tag data.

## Steps to Build and Upload

1. Open the Arduino IDE.
2. Load the `aa2.ino` sketch file.
3. Select the correct board and port:

    - Board: `Arduino Mega 2560`
    - Port: `COM6` (or the port your Arduino is connected to).

4. Install the required libraries if not already installed:

    - Go to **Sketch > Include Library > Manage Libraries**.
    - Search for and install `SafeString` and `LiquidCrystal_I2C`.

5. Compile the sketch to ensure there are no errors.
6. Upload the sketch to the Arduino Mega 2560.

---

## Program Description

This program is designed to manage and display system information on an LCD screen, handle user inputs via buttons, and process data received through serial communication. It is primarily used for monitoring and controlling an RFID-based system.

### Key Features

- **LCD Display**: Displays system statuses, tag counts, network information, and other data.
- **Button Navigation**: Two buttons (`BUTTON_VANCE` and `BUTTON_START`) allow users to navigate between screens and perform actions.
- **Serial Communication**: Parses and processes data received via serial input.
- **Screen Locking**: Locks the screen for specific operations and provides confirmation prompts for critical actions.
- **System Actions**: Includes functionalities like uploading data, resetting, shutting down, and managing backups.
- **Power-Off Countdown**: Implements a countdown before shutting down the system.

### Screens

1. **Informational Screens**:

    - Display tag counts, network statuses, and system version.

2. **Action Screens**:

    - Allow actions like uploading data, resetting, or shutting down.

3. **Confirmation Screens**:

    - Prompt the user for confirmation before executing critical actions.

4. **Power-Off Screen**:

    - Displays a countdown before shutting down the system.

---

## Serial Communication Input Format

The program now supports two types of serial messages: **Antenna Data (A)** and **PCData (P)**. The format is as follows:

```perl
$MYTMP;<tags>;<unique_tags>;<type>;<data>;<timestamp>*<checksum>
```

### Field Descriptions

- **`<tags>`**: Total number of tags read (integer).
- **`<unique_tags>`**: Number of unique tags read (integer).
- **`<type>`**: Message type (`A` for Antenna Data, `P` for PCData).
- **`<data>`**: Data fields specific to the message type (see below).
- **`<timestamp>`**: The UNIX timestamp for updating time.
- **`<checksum>`**: XOR checksum of the message (excluding the `$`).

#### Antenna's `<data>` Fields (`A`)

- **`<antenna1>`**: Count of tags read by antenna 1 (integer).
- **`<antenna2>`**: Count of tags read by antenna 2 (integer).
- **`<antenna3>`**: Count of tags read by antenna 3 (integer).
- **`<antenna4>`**: Count of tags read by antenna 4 (integer).

#### PCData's `<data>` Fields (`P`)

- **`<comm_status>`**: Communication status (`1` for true, `0` for false).
- **`<rfid_status>`**: RFID reader status (`1` for true, `0` for false).
- **`<usb_status>`**: USB connection status (`1` for true, `0` for false).
- **`<sys_version>`**: System version number (integer).
- **`<num_serie>`**: Serial number of the system (integer).
- **`<backups>`**: Number of backups stored (integer).
- **`<permanent_unique_tags>`**: Number of permanent unique tags (integer).

### Examples

#### Antenna Data Example

```perl
$MYTMP;12345;678;A;100;200;300;400;1744846380*<checksum>
```

- **`<tags>`**: 12345
- **`<unique_tags>`**: 678
- **`<type>`**: `A` (Antenna Data)
- **`<antenna1>`**: 100
- **`<antenna2>`**: 200
- **`<antenna3>`**: 300
- **`<antenna4>`**: 400
- **`<timestamp>`**: 1744846380 (Wed, 16 Apr 2025 23:33:00 GMT)

#### PCData Example

```perl
$MYTMP;12345;678;P;1;1;0;42;123;5;100;1744846380*<checksum>
```

- **`<tags>`**: 12345
- **`<unique_tags>`**: 678
- **`<type>`**: `P` (PCData)
- **`<comm_status>`**: 1 (true)
- **`<rfid_status>`**: 1 (true)
- **`<usb_status>`**: 0 (false)
- **`<sys_version>`**: 42
- **`<num_serie>`**: 123
- **`<backups>`**: 5
- **`<permanent_unique_tags>`**: 100
- **`<timestamp>`**: 1744846380 (Wed, 16 Apr 2025 23:33:00 GMT)

---

## Usage Instructions

1. **Navigation**:

    - Use `BUTTON_VANCE` to navigate between screens.
    - Use `BUTTON_START` to trigger actions on the current screen.

2. **Serial Communication**:

     - Send data in the specified format to update system information.

3. **Monitor LCD**:

     - Observe the LCD for system statuses, prompts, and confirmation requests.

4. **Power-Off Countdown**:

    - When shutting down, the system will display a countdown. Wait for the countdown to complete before powering off the device.

---

## Recent Changes

- **Added Power-Off Countdown**: The system now includes a countdown before shutting down.
- **Improved Serial Parsing**: Enhanced parsing of serial messages to include additional fields like timestamps and antenna data.
- **Screen Locking**: Added functionality to lock and unlock the screen for specific operations.
- **Confirmation Prompts**: Added confirmation screens for critical actions like data deletion and shutdown.
- **New Screens**: Added screens for displaying antenna data and USB status
