#include <linux/bpf.h>
#include <linux/if_link.h>
#include <assert.h>
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <arpa/inet.h>
#include <netinet/ether.h>
#include <unistd.h>
#include <time.h>

// 24628 should be 10.8.0.3
void main()
{
    char ip_param[] = "8.8.8.8";
    struct sockaddr_in sa_param;
    inet_pton(AF_INET, ip_param, &(sa_param.sin_addr));
    __u32 ip = sa_param.sin_addr.s_addr;
    printf("the ip to filter is %s/%u\n", ip_param, ip);
}
