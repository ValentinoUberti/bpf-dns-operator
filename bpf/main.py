#!/usr/bin/python

from bcc import BPF
import socket
import os

device = "enp7s0u2u1u2"

bpf = BPF(src_file="filter.c")
fn = bpf.load_func("ip_filter", BPF.XDP)
BPF.attach_xdp(device, fn, 0)

try:
    bpf.trace_print()
except KeyboardInterrupt:
    pass

bpf.remove_xdp(device, 0)
