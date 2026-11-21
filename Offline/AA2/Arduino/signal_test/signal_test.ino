/**
 * Signal bar LCD test sketch.
 *
 * Displays 4 antenna signal indicators on a 20x4 I2C LCD.
 * Send single digits 0-3 via Serial (115200) to cycle antenna levels.
 * Each digit press cycles the next antenna (A1, A2, A3, A4, A1, ...).
 *
 * Layout:
 *   Row 0: "Signal Bar Test"
 *   Row 1: A1: [][][]  A2: [][][]
 *   Row 2: A3: [][][]  A4: [][][]
 *   Row 3: "Send 0-3 to set"
 */

#include <LiquidCrystal_I2C.h>

LiquidCrystal_I2C lcd(0x27, 20, 4);

byte signal_outline[] = {
	B11111,
	B10001,
	B10001,
	B10001,
	B10001,
	B10001,
	B10001,
	B11111
};

byte signal_filled[] = {
	B11111,
	B11111,
	B11111,
	B11111,
	B11111,
	B11111,
	B11111,
	B11111
};

#define CHAR_OUTLINE 1
#define CHAR_FILLED 2

int levels[4] = {0, 1, 2, 3}; // start showing all levels
int current_antenna = 0;

void draw_signal_bars(int col, int row, int level)
{
	lcd.setCursor(col, row);
	for (int i = 0; i < 3; i++)
	{
		lcd.write((i < level) ? CHAR_FILLED : CHAR_OUTLINE);
	}
}

void draw_screen()
{
	lcd.setCursor(0, 0);
	lcd.print("Signal Bar Test     ");

	lcd.setCursor(0, 1);
	lcd.print("A1:    A2:   ");
	lcd.print("       ");
	draw_signal_bars(3, 1, levels[0]);
	draw_signal_bars(10, 1, levels[1]);

	lcd.setCursor(0, 2);
	lcd.print("A3:    A4:   ");
	lcd.print("       ");
	draw_signal_bars(3, 2, levels[2]);
	draw_signal_bars(10, 2, levels[3]);

	lcd.setCursor(0, 3);
	lcd.print("Send 0-3 to set lvl ");
}

void setup()
{
	lcd.init();
	lcd.backlight();

	lcd.createChar(CHAR_OUTLINE, signal_outline);
	lcd.createChar(CHAR_FILLED, signal_filled);

	Serial.begin(115200);
	while (!Serial)
		;

	draw_screen();

	Serial.println("Signal Bar Test Ready");
	Serial.println("Send 0, 1, 2, or 3 to set level for next antenna");
	Serial.print("Current: A1=");
	Serial.print(levels[0]);
	Serial.print(" A2=");
	Serial.print(levels[1]);
	Serial.print(" A3=");
	Serial.print(levels[2]);
	Serial.print(" A4=");
	Serial.println(levels[3]);
}

void loop()
{
	if (Serial.available())
	{
		char c = Serial.read();

		if (c >= '0' && c <= '3')
		{
			int level = c - '0';
			levels[current_antenna] = level;

			Serial.print("A");
			Serial.print(current_antenna + 1);
			Serial.print(" -> ");
			Serial.println(level);

			current_antenna = (current_antenna + 1) % 4;

			Serial.print("Next: A");
			Serial.println(current_antenna + 1);

			draw_screen();
		}
	}
}
