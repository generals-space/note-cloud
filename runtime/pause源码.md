参考文章

1. [pause.c](https://github.com/kubernetes/kubernetes/blob/master/build/pause/linux/pause.c)

```c++
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#define STRINGIFY(x) #x
#define VERSION_STRING(x) STRINGIFY(x)

#ifndef VERSION
#define VERSION HEAD
#endif

static void sigdown(int signo) {
    psignal(signo, "Shutting down, got signal");
    exit(0);
}

static void sigreap(int signo) {
    while (waitpid(-1, NULL, WNOHANG) > 0)
        ;
}

int main(int argc, char **argv) {
    int i;
    for (i = 1; i < argc; ++i) {
        if (!strcasecmp(argv[i], "-v")) {
            printf("pause.c %s\n", VERSION_STRING(VERSION));
            return 0;
        }
    }

    if (getpid() != 1) {
        /* Not an error because pause sees use outside of infra containers. */
        fprintf(stderr, "Warning: pause should be the first process\n");
    }

    if (sigaction(SIGINT, &(struct sigaction){.sa_handler = sigdown}, NULL) < 0) {
        return 1;
    }
    if (sigaction(SIGTERM, &(struct sigaction){.sa_handler = sigdown}, NULL) < 0) {
        return 2;
    }
    if (sigaction(SIGCHLD, &(struct sigaction){.sa_handler = sigreap, .sa_flags = SA_NOCLDSTOP}, NULL) < 0) {
        return 3;
    }

    for (;;) {
        pause();
    }

    fprintf(stderr, "Error: infinite loop terminated\n");
    return 42;
}
```
