#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>

#include "nmconnect.h"

#define RED "\033[31;1m"
#define YLW "\033[33;1m"
#define RST "\033[0m"

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

#define DATA_BUF_SIZE 36

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
  // sed -i "/^ssid=/ s/=.*\$/=$NETWORK/"
  // /etc/NetworkManager/system-connections/Wifi.nmconnection

  char command[18 + DATA_BUF_SIZE];
  sprintf(command, "/^ssid=/ s/=.*$/=%s/", network);

  COMMAND("sed", "sed", "-i", command, NETWORK_CONFIG_FILE, (char *)NULL);
  COMMAND("nmcli", "nmcli", "reload", (char *)NULL);
}

int main(void) {
  int fd = -1;

  char *network, *password;
  network = NULL;
  password = NULL;

  int count = -1; /* length of the network name (will be calculated later) */

  char data[DATA_BUF_SIZE] = {0};

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
