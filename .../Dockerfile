FROM registry.redhat.io/rhel8/llvm-toolset
USER 0
#RUN dnf install clang llvm -y
RUN dnf install go -y
RUN go get github.com/dropbox/goebpf
COPY ./ .
RUN clang -I. -O2 -target bpf -c ebpf_prog/xdp_fw.c  -o ebpf_prog/xdp_fw.elf
RUN go build main.go
CMD ['sh']
#main -iface tap0 -drop 8.8.8.8

#yum install bcc bcc-tools kernel-devel-`uname -r`