# Install goebpf package
$ go get github.com/dropbox/goebpf
$ make clang -I../../.. -O2 -target bpf -c ebpf_prog/xdp.c  -o ebpf_prog/xdp.elf
go build -v -o main


clang -I./ -O2 -target bpf -c ebpf_prog/xdp_fw.c  -o ebpf_prog/xdp_fw.elf
/home/vale/Projects/Operator-Sdk/bpf-dns-operator/bpf


###
git clone https://github.com/dropbox/goebpf.git
https://stackoverflow.com/questions/60322147/xdp-program-ip-link-error-prog-section-rejected-operation-not-permitted

CONFIG_BPF_SYSCALL?
cat /boot/config-5.10.19-200.fc33.x86_64 | grep CONFIG_BPF_SYSCALL

https://www.linuxtechi.com/set-ulimit-file-descriptors-limit-linux-servers/

