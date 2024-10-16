#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>

#include "nmconnect.h"

#define RED "\033[31;1m"
#define YLW "\033[33;1m"
#define RST "\033[0m"

/* dbg */
void _Noreturn fatal(const char *func, int line) {
  fprintf(stderr, RED "Fatal error at '%s', line %d\n" RST, func, line);
  exit(1);
}

#define FATAL() fatal(__FUNCTION__, __LINE__)

/*
 * parse space separator
 * and remove trailing '\n'
 * */
int parse(char *data) {

  char *c;

  int count = -1;

  for (c = data; *c != '\0' && *c != '\n'; c++) {
    if (*c == ' ') {
      *c = '\0';

      count = (int)(c - data);
    }
  }

  /* the for loop stops in either '\0' or '\n'
   * so i set it to '\0' anyway. This may be redundant for
   * the '\0' case, but it removes the newline in case
   * it's there. */
  *c = '\0';

  return count;
}

#define DATA_BUF_SIZE 64

#define COMMAND(...)                                                           \
  do {                                                                         \
    int pid, status;                                                           \
    pid = fork();                                                              \
    if (0 > pid)                                                               \
      FATAL();                                                                 \
    if (pid == 0)                                                              \
      execlp(__VA_ARGS__);                                                     \
    if (0 > wait(&status))                                                     \
      FATAL();                                                                 \
    if (WIFEXITED(status))                                                     \
      fprintf(stderr, "Info: got status %d\n", WEXITSTATUS(status));           \
  } while (0)

void connect(const char *network, const char *password) {

  /*
   * sed -i
   * 	"/^ssid=/ s/=.*\$/=$NETWORK/"
   * 	/etc/NetworkManager/system-connections/Wifi.nmconnection
   */

/* TODO: defs in nmconnect.h */
#define NETWORK_SSID_LENGTH 30 /* '%s' is counted by LEN so +2 chars */
#define NETWORK_PSK_LENGTH 30

#define STRING(v, s)                                                           \
  const char *v = s;                                                           \
  const int v##_sz = sizeof s;
#define LEN(v) v##_sz

  STRING(network_fmt, "/^ssid=/ s/=.*$/=%" QUOTE(NETWORK_SSID_LENGTH) "s/;");
  STRING(password_fmt, "/^psk=/ s/=.*$/=%" QUOTE(NETWORK_PSK_LENGTH) "s/");

  /* Creating the sed command */
  char command[LEN(network_fmt) + LEN(password_fmt) + NETWORK_SSID_LENGTH +
               NETWORK_PSK_LENGTH + DATA_BUF_SIZE];

  int count = 0;

  count = sprintf(command, network_fmt, network);
  sprintf(command + count, password_fmt, password);

  COMMAND("sed", "sed", "-i", command, NETWORK_CONFIG_FILE, (char *)NULL);
  // COMMAND("nmcli", "nmcli", "reload", (char *)NULL);
}

int main(void) {
  int fd = -1; /* file descriptor for the FIFO communication file */

  char *network, *password;
  network = NULL;
  password = NULL;

  int count = -1; /* length of the network name (will be calculated later) */

  char data[DATA_BUF_SIZE] = {0}; /* REMINDER: this is Network AND Password! */

  for (;;) {

    if (0 > (fd = open(NETWORK_FIFO, O_RDONLY)))
      FATAL();

    read(fd, data, sizeof(data));

    if (0 > (count = parse(data))) {
      fprintf(stderr, YLW "Warning: couldn't parse '%s'\n" RST, data);
      continue;
    }

    network = data;
    password = network + count + 1; /* skipping the \0 */

    if (network == NULL || password == NULL)
      FATAL();

    /* dbg */
    fprintf(stderr,
            YLW "Debg: network  -> %p\n"
                "|     password -> %p\n" RST,
            network, password);

    fprintf(stderr, "Info: connect( %s, %s )\n", network, password);

    connect(network, password);

    close(fd);
  }
}
