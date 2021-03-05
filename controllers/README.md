https://duo.com/labs/tech-notes/writing-an-xdp-network-filter-with-ebpf

https://blogs.oracle.com/linux/notes-on-bpf-1

https://www.openshift.com/blog/linux-capabilities-in-openshift

https://github.com/dropbox/goebpf/tree/master/examples


https://prototype-kernel.readthedocs.io/en/latest/bpf/ebpf_maps.html



# Install clang/llvm to be able to compile C files into bpf arch
$ apt-get install clang llvm make

# Install goebpf package
$ go get github.com/dropbox/goebpf
$ make clang -I../../.. -O2 -target bpf -c ebpf_prog/xdp.c  -o ebpf_prog/xdp.elf
go build -v -o main

