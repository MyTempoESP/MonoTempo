#include <LiquidCrystal_I2C.h>
#include <string.h>

#define LABEL_COUNT 33

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
  "DATA: ", // 22

  "PRESSIONE CONFIRMA",
  "PARA FAZER UPLOAD",
  "DE ATLETAS", // 25

  "DOS BACKUPS", // 26

  "UPLOAD EM ANDAMENTO", // 27

  "<START: RESET TELA>",
  "<START: RESET WIFI>",
  "<START: RESET 4G>",
  "<START: BACKUP USB>",
  "<START: APAGA TUDO>" // 32
};
const int labels_len[LABEL_COUNT] = {
  11,9,9,12,7,8,6,4,7,7,6,6,5,10,12,11,9,9,9,9,14,7,6,18,17,10,11,19,19,19,17,19,19
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

#define VIRT_SCR_COLS 20
#define VIRT_SCR_ROWS 4

uint8_t g_x, g_y;
char g_virt_scr[VIRT_SCR_ROWS][VIRT_SCR_COLS];

#define virt_scr_sprintf(fmt, ...)																			\
  snprintf(g_virt_scr[g_y] + g_x, ((VIRT_SCR_COLS + 1) - g_x), fmt, __VA_ARGS__);

LiquidCrystal_I2C lcd(0x27, VIRT_SCR_COLS, VIRT_SCR_ROWS);

#define SERIAL_BUFSIZE 256
char serial_buffer[SERIAL_BUFSIZE];

bool new_data = false;

int
serial_recv_delimited()
{
	char start_delimiter = 0x3C;
	char end_delimiter = 0x3E;

	char rb;
	int n;

	new_data = false;

	for (rb = 0; Serial.available() > 0; rb = Serial.read()) if (rb == start_delimiter) goto read_until_delimiter; return;

 read_until_delimiter:
	for (rb = Serial.read(), n = 0, new_data = true						\
				 ; Serial.available() > 0														\
				 && rb != end_delimiter															\
				 && n < SERIAL_BUFSIZE - 1													\
				 ; n++, rb = Serial.read()) serial_buffer[n] = rb;

	serial_buffer[n] = '\0';

	return n;
}

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
	char* c;
	
	int count = serial_recv_delimited();

	if (!new_data) return;
	
	char* start = serial_buffer;

	for (int i = 0; i < VIRT_SCR_ROWS; i++) {
		for (c = start; *c != '\0' && *c != '\n'; c++);
		int len = c - start;
		len = (len<(VIRT_SCR_COLS-1)?len:(VIRT_SCR_COLS-1));
		memcpy(g_virt_scr[i], start, len);
		g_virt_scr[i][len] = '\0';
		if (*c == '\0') return;
		
		start = ++c;
	}
}
