// Copyright (c) 2019 Dropbox, Inc.
// Full license can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dropbox/goebpf"
	"golang.org/x/sys/unix"
)

type ipAddressList []string

// Set the RLIMIT_MEMLOCK limits to infinity <--- REQUIRED!
func setRLimitMemlock() error {

	var rLimit syscall.Rlimit
	rLimit.Max = unix.RLIM_INFINITY
	rLimit.Cur = unix.RLIM_INFINITY

	err := syscall.Setrlimit(unix.RLIMIT_MEMLOCK, &rLimit)
	return err
}

const (
	// Size of structure used to pass metadata
	metadataSize = 12
)

// In sync with xdp_dump.c  "struct perf_event_item"
type perfEventItem struct {
	SrcIp, DstIp uint32
}

func main() {

	var elf = flag.String("elf", "ebpf_prog/xdp_fw_events.elf", "clang/llvm compiled binary file")
	var ipList ipAddressList

	if setRLimitMemlock() != nil {
		fatalError("Error setting RLIMIT_MEMLOCK: check permissions")
	}

	iface, exists := os.LookupEnv("BPF_IFACE")
	if !exists {
		fatalError("BPF_IFACE env not set")
	}

	flag.Var(&ipList, "drop", "IPv4 CIDR to DROP traffic from, repeatable")
	flag.Parse()

	if len(ipList) == 0 {
		fatalError("at least one IPv4 address to DROP required (-drop)")
	}

	// Create eBPF system
	bpf := goebpf.NewDefaultEbpfSystem()
	// Load .ELF files compiled by clang/llvm
	err := bpf.LoadElf(*elf)
	if err != nil {
		fatalError("LoadElf() failed: %v", err)
	}
	printBpfInfo(bpf)

	// Get eBPF maps
	matches := bpf.GetMapByName("matches")
	if matches == nil {
		fatalError("eBPF map 'matches' not found")
	}
	blacklist := bpf.GetMapByName("blacklist")
	if blacklist == nil {
		fatalError("eBPF map 'blacklist' not found")
	}

	// Find special "PERF_EVENT" eBPF map
	perfmap := bpf.GetMapByName("perfmap")
	if perfmap == nil {
		fatalError("eBPF map 'perfmap' not found")
	}

	// Get XDP program. Name simply matches function from xdp_fw.c:
	//      int firewall(struct xdp_md *ctx) {
	xdp := bpf.GetProgramByName("firewall")
	if xdp == nil {
		fatalError("Program 'firewall' not found.")
	}

	// Populate eBPF map with IPv4 addresses to block
	fmt.Println("Blacklisting IPv4 addresses...")
	for index, ip := range ipList {
		fmt.Printf("\t%s\n", ip)
		err := blacklist.Insert(goebpf.CreateLPMtrieKey(ip), index)
		if err != nil {
			fatalError("Unable to Insert into eBPF map: %v", err)
		}
	}
	fmt.Println()

	// Load XDP program into kernel
	err = xdp.Load()
	if err != nil {
		fatalError("xdp.Load(): %v", err)
	}

	// Attach to interface
	err = xdp.Attach(iface)
	if err != nil {
		fatalError("xdp.Attach(): %v", err)
	}
	defer xdp.Detach()

	// Add CTRL+C handler
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)

	// Start listening to Perf Events
	perf, _ := goebpf.NewPerfEvents(perfmap)
	perfEvents, err := perf.StartForAllProcessesAndCPUs(4096)
	if err != nil {
		fatalError("perf.StartForAllProcessesAndCPUs(): %v", err)
	}

	fmt.Println("XDP program successfully loaded and attached. Counters refreshed every second.")
	fmt.Println("Press CTRL+C to stop.")
	fmt.Println()

	/*
		// Print stat every second / exit on CTRL+C
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("IP                 DROPs")
				for i := 0; i < len(ipList); i++ {
					value, err := matches.LookupInt(i)
					if err != nil {
						fatalError("LookupInt failed: %v", err)
					}
					fmt.Printf("%18s    %d\n", ipList[i], value)
				}
				fmt.Println()
			case <-ctrlC:
				fmt.Println("\nDetaching program and exit")
				return
			}
		}*/

	go func() {
		var event perfEventItem
		for {
			if eventData, ok := <-perfEvents; ok {
				reader := bytes.NewReader(eventData)
				binary.Read(reader, binary.LittleEndian, &event)
				fmt.Printf("TCP: %v -> %v BLOCKED ",
					intToIPv4(event.SrcIp),
					intToIPv4(event.DstIp),
				)
				if len(eventData)-metadataSize > 0 {
					// event contains packet sample as well
					fmt.Println(hex.Dump(eventData[metadataSize:]))
				}

				for i := 0; i < len(ipList); i++ {
					value, err := matches.LookupInt(i)
					if err != nil {
						fatalError("LookupInt failed: %v", err)
					}

					// Remove last three chars "/xx"
					if ipList[i][:len(ipList[i])-3] == intToIPv4(event.SrcIp).String() {
						if value <= 1 {
							fmt.Printf("%d time\n", value)
						} else {
							fmt.Printf("%d times\n", value)

						}

					}

					//fmt.Printf("%18s    %d\n", ipList[i], value)
				}
			} else {
				// Update channel closed
				break
			}
		}
	}()

	// Wait until Ctrl+C pressed
	<-ctrlC

	// Stop perf events and print summary
	perf.Stop()
	fmt.Println("\nSummary:")
	fmt.Printf("\t%d Event(s) Received\n", perf.EventsReceived)
	fmt.Printf("\t%d Event(s) lost (e.g. small buffer, delays in processing)\n", perf.EventsLost)
	fmt.Println("\nDetaching program and exit...")

}

func fatalError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func printBpfInfo(bpf goebpf.System) {
	fmt.Println("Maps:")
	for _, item := range bpf.GetMaps() {
		fmt.Printf("\t%s: %v, Fd %v\n", item.GetName(), item.GetType(), item.GetFd())
	}
	fmt.Println("\nPrograms:")
	for _, prog := range bpf.GetPrograms() {
		fmt.Printf("\t%s: %v, size %d, license \"%s\"\n",
			prog.GetName(), prog.GetType(), prog.GetSize(), prog.GetLicense(),
		)

	}
	fmt.Println()
}

// Implements flag.Value
func (i *ipAddressList) String() string {
	return fmt.Sprintf("%+v", *i)
}

// Implements flag.Value
func (i *ipAddressList) Set(value string) error {
	if len(*i) == 16 {
		return errors.New("Up to 16 IPv4 addresses supported")
	}
	// Validate that value is correct IPv4 address
	if !strings.Contains(value, "/") {
		value += "/32"
	}
	if strings.Contains(value, ":") {
		return fmt.Errorf("%s is not an IPv4 address", value)
	}
	_, _, err := net.ParseCIDR(value)
	if err != nil {
		return err
	}
	// Valid, add to the list
	*i = append(*i, value)
	return nil
}

func intToIPv4(ip uint32) net.IP {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, ip)
	return net.IP(res)
}

func ntohs(value uint16) uint16 {
	return ((value & 0xff) << 8) | (value >> 8)
}
