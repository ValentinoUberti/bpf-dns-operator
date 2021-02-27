#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
int ip_filter(struct xdp_md *ctx)
{

    unsigned int ip = 134744072; //8.8.8.8

    //bpf_trace_printk("got a packet\n");
    void *data = (void *)(long)ctx->data;         // (8)
    void *data_end = (void *)(long)ctx->data_end; // (8)
    struct ethhdr *eth = data;

    // check packet size
    if ((void *)eth + sizeof(*eth) > data_end)
    {
        bpf_trace_printk("First size check");
        return XDP_PASS;
    }

    //check if the packet is an IP packet
    if (ntohs(eth->h_proto) != ETH_P_IP)
    {
        return XDP_PASS;
    }

    struct iphdr *iph = data + sizeof(struct ethhdr);
    if ((void *)(iph + 1) > data_end)
    {
        bpf_trace_printk("Error final check");
        return XDP_PASS;
    }

    if ((void *)iph + sizeof(*iph) <= data_end)
    {

        unsigned int ip_src = iph->saddr;
        //bpf_trace_printk("Destination ip address is %u\n", ip_src);

        // drop the packet if the ip source address is equal to ip
        if (ip_src == ip)
        {
            bpf_trace_printk("Dropping destination ip");
            return XDP_DROP;
        }
    }

    return XDP_PASS;
}