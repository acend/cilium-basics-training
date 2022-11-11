package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	//handle Ctrl+c
	stopper := make(chan os.Signal, 1)

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Here we load our bpf code into the kernel, these functions are in the
	// .go file created by bpf2go
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	//SEC("tracepoint/syscalls/sys_enter_execve")
	kp, err := link.Tracepoint("syscalls", "sys_enter_execve", objs.BpfProg, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
	}
	defer kp.Close()

	// code below is only to show tracepipe data in stdout (not relevant for eBPF)
	const tracePipeFile = "/sys/kernel/debug/tracing/trace_pipe"

	f, err := os.Open(tracePipeFile)
	if err != nil {
		log.Fatalf("opening trace_pipe: %s", err)
	}
	reader := bufio.NewReader(f)

	for {
		select {
		case <-stopper:
			break
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					continue
				}
				log.Fatalf("error read trace_pipe: %s", err)
			}
			fmt.Printf("%+v\n", line)
		}

	}

	log.Println("Received signal, exiting program..")
}
