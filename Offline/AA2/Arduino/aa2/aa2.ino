//#include <EnableInterrupt.h>
#include <LiquidCrystal_I2C.h>
//#include <Wire.h>
//#include <HardwareSerial.h>
#include <nanoFORTH.h>
#include <string.h>

#define LABEL_COUNT 23

const char* labels[] = {
  "PORTAL   My",
  "ATLETAS  ",
  "REGIST.  ",
  "COMUNICANDO ",
  "LEITOR ",
  "LTE/4G: ",
  "WIFI: ",
  "IP: ",
  "LOCAL: ",
  "PROVA: ",
  "PING: ",
  "HORA: ",
  "USB: ",
  "AGUARDE...",
  "ERRO, TENTAR",
  "  NOVAMENTE", // 15

  "RFID  -  ",
  "SERIE:   ",
  "SIST.    ", // 18

  "PRESSIONE",
  "PARA CONFIRMAR", // 20

  "OFFLINE",
  "DATA: " // 22
};
const int labels_len[LABEL_COUNT] = {
  11,9,9,12,7,8,6,4,7,7,6,6,5,10,12,11,9,9,9,9,14,7,6
};

#define VALUE_COUNT 9

const char* values[] = {
  "WEB",
  "CONECTAD",
  "DESLIGAD",
  "AUTOMATIC",
  "OK",
  "X",
  "  ",
  "A",
  ": "
};

const char code[] PROGMEM =          ///< define preload Forth code here

// Button.fth
  "VAR bac\n"
  "VAR bst\n"
  "VAR ba2\n"
  "VAR bs2\n"
  ": btn 6 IN 0 = ;\n"
  ": bt2 7 IN 0 = ;\n"
  ": b1 btn DUP bst @ NOT AND IF 1 bac ! THN bst ! ;\n"
  ": b2 bt2 DUP bs2 @ NOT AND IF 1 ba2 ! THN bs2 ! ;\n"
  ": chb b1 b2 ;\n"
  ": ba@ bac @ . ;\n"
  ": b2@ ba2 @ . ;\n"

// Screen.fth
  ": lbl  5   API ;\n"
  ": fwd  2   API ;\n"
  ": lit  API fwd ;\n"
  ": fnm  1   lit ;\n"
  ": fni  1   API ;\n" // Multi-Column
  ": num  4   lit ;\n"
  ": nui  4   API ;\n" // Multi-Column
  ": val  6   lit ;\n"
  ": ip   7   lit ;\n"
  ": ms   3   lit ;\n"
  ": hms  256 ip  ;\n"
  ": usb  12  lbl ;\n"
  ": tim  11  lbl ;\n"
  ": hex  16  fnm ;\n"
  
  // Text Decorations
  ": a    7 6 API ;\n" // Multi-Column
  ": spc  6 6 API ;\n" // Multi-Column
  ": sep  8 6 API ;\n" // Multi-Column

  // Antenna Data
  ": atn " // ( N Mag N Mag N Mag N Mag -- )
    "a 1 nui sep fni spc a 2 nui sep fnm "
    "a 3 nui sep fni spc a 4 nui sep fnm "
  ";\n"

  "10 0 TMI chb 1 TME\n"
;

#define VIRT_SCR_COLS 20
#define VIRT_SCR_ROWS 4

uint8_t g_x, g_y;
char g_virt_scr[VIRT_SCR_ROWS][VIRT_SCR_COLS];

#define virt_scr_sprintf(fmt, ...) \
  snprintf(g_virt_scr[g_y] + g_x, ((VIRT_SCR_COLS + 1) - g_x), fmt, __VA_ARGS__);

LiquidCrystal_I2C lcd(0x27, VIRT_SCR_COLS, VIRT_SCR_ROWS);

void
setup()
{
  lcd.init();      // Initialize the LCD
  lcd.backlight(); // Turn on the backlight
  
  memset(g_virt_scr, '\0', sizeof(g_virt_scr));

  Serial.begin(115200);
  while(!Serial);

  n4_setup(code);

  n4_api(0, draw);
  n4_api(1, print_forthNumber);
  n4_api(2, forth_line_feed);
  n4_api(3, forth_millis);

  n4_api(4, forth_number);
  n4_api(5, forth_label);
  n4_api(6, forth_value);
  n4_api(7, forth_ip);

  pinMode(7, INPUT_PULLUP);
  pinMode(6, INPUT_PULLUP);
}

void
forth_millis()
{
  int v;

  if ((v = n4_pop()) < 1000) {
    g_x += virt_scr_sprintf("%dms", v);
    return;
  }

  v /= 1000;
  g_x += virt_scr_sprintf("%ds", v);
}

void
forth_value()
{
  int v;

  if ((v = n4_pop()) > VALUE_COUNT || v < 0) return;

  g_x += virt_scr_sprintf("%s", values[v]);
}

void
print_forthNumber()
{
  int mag, v;
  char postfix;

  mag = n4_pop();
  v = n4_pop();

  if (mag == 16) { // (special case) hex
	  g_x += virt_scr_sprintf("%04x", v);
	  return;
  }

  postfix = (mag == 0) ?
      ' ' :
      (mag >= 3 && mag < 6 ? 'K' : 'M');

  // 'X'  if Magnitude = 0, 'XK' if 6 > Magnitude >= 3
  // 'XM' if Magnitude >= 6

  g_x += virt_scr_sprintf("%d%c", v, postfix);
}

void
forth_ip()
{
  int f = n4_pop();

  if (f >= 0xDA7E) {
    g_x += virt_scr_sprintf("%02d/%02d/%04d", n4_pop(), n4_pop(), n4_pop());
  } else if (f >= 256) {
    g_x += virt_scr_sprintf("%02d:%02d:%02d", n4_pop(), n4_pop(), n4_pop());
  } else {
    g_x += virt_scr_sprintf( "%d.%d.%d.%d", n4_pop(), n4_pop(), n4_pop(), f);
  }
}

void
forth_number()
{
  g_x += virt_scr_sprintf("%d", n4_pop());
}

void
forth_label()
{
  int v;

  if ((v = n4_pop()) >= LABEL_COUNT || v < 0) return;

  g_x = labels_len[v];

  memcpy(g_virt_scr[g_y], labels[v], labels_len[v]);
}

void
forth_line_feed()
{
  for (; g_x < VIRT_SCR_COLS - 1; g_x++)
    g_virt_scr[g_y][g_x] = ' ';

  g_virt_scr[g_y][g_x] = '\0';

  g_x = 0;

  g_y++;

  if (g_y >= (VIRT_SCR_ROWS - 1))
    g_y = VIRT_SCR_ROWS - 1;
}

void
draw()
{
  // resetting
  g_y = 0;
  g_x = 0;

  for (int i = 0; i < VIRT_SCR_ROWS; i++){

    lcd.setCursor(0, i);

    for (char* c = g_virt_scr[i], i = 0; *c != '\0' && i < VIRT_SCR_COLS; c++, i++)
      lcd.write(*c);
  }
}

void
loop()
{
  n4_run();
}
