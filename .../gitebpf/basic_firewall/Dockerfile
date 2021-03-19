FROM registry.redhat.io/rhel8/llvm-toolset AS builder
USER 0
#RUN dnf install clang llvm -y
RUN dnf install go -y
#RUN go get github.com/dropbox/goebpf
COPY ./ .
RUN clang -I. -O2 -target bpf -c ebpf_prog/xdp_fw_events.c  -o ebpf_prog/xdp_fw_events.elf
RUN go build main.go
#RUN dnf install iproute iputils -y


FROM registry.access.redhat.com/ubi8/ubi-minimal
WORKDIR /opt/app-root/src/
COPY --from=builder /opt/app-root/src/main ./firewall
COPY --from=builder /opt/app-root/src/ebpf_prog/xdp_fw_events.elf ./ebpf_prog/xdp_fw_events.elf
CMD [ "/bin/bash" ]