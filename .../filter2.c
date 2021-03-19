#define RATIO 10
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <linux/tcp.h>
#include <linux/inet.h>


int ip_filter(struct __sk_buff *skb)
{
    void *data = (void *)(uintptr_t)skb->data;
    void *data_end = (void *)(uintptr_t)skb->data_end;
    struct ethhdr *eth = data;
    struct iphdr *iph = (struct iphdr *)(eth + 1);
    struct tcphdr *tcphdr = (struct tcphdr *)(iph + 1);

    /* sanity check needed by the eBPF verifier */
    if ((void *)(tcphdr + 1) > data_end)
        return TC_ACT_OK;

   
    return TC_ACT_OK;
}

