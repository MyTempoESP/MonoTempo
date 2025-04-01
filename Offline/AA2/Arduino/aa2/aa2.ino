#include <LiquidCrystal_I2C.h>
#include <string.h>

#define VIRT_SCR_COLS 20
#define VIRT_SCR_ROWS 4

uint8_t g_x, g_y;
char g_virt_scr[VIRT_SCR_ROWS][VIRT_SCR_COLS + 1];

#define virt_scr_sprintf(x, y, fmt, ...) \
  snprintf(g_virt_scr[y] + x, (VIRT_SCR_COLS - x), fmt, __VA_ARGS__);

LiquidCrystal_I2C lcd(0x27, VIRT_SCR_COLS, VIRT_SCR_ROWS);

void
setup()
{
  lcd.init();      // Initialize the LCD
  lcd.backlight(); // Turn on the backlight
  
  memset(g_virt_scr, '\0', sizeof(g_virt_scr));

  Serial.begin(115200);
  while(!Serial);

  pinMode(7, INPUT_PULLUP);
  pinMode(6, INPUT_PULLUP);
}

void
draw()
{

  for (int i = 0; i < VIRT_SCR_ROWS; i++){

    lcd.setCursor(0, i);

    for (char* c = g_virt_scr[i], i = 0; *c != '\0' && i < VIRT_SCR_COLS; c++, i++)
      lcd.write(*c);
  }
}

void
loop()
{
  handleButtons();
  handleSerial();
  draw();
}
