
#include "common.h"

// SEC is a macro that expands to create an ELF section which bpf loaders parse.
// we want our function to be executed whenever syscall execve (program execution) is called
SEC("tracepoint/syscalls/sys_enter_execve")
int bpf_prog(void *ctx) {
  char msg[] = "Hello world";
  // bpf_printk is a bpf helper function which writes strings to /sys/kernel/debug/tracing/trace_pipe (good for debugging purposes)
  bpf_printk("%s", msg);
  // bpf programs need to return an int
  return 0;
}

char LICENSE[] SEC("license") = "GPL";