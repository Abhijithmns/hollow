## Namespaces

A namespace wraps A namespace wraps a global system resource in an abstraction that makes it appear to the processes within the namespace that they have 
their own isolated instance of the global resource

for the ones mentioned in the runtime-spec

- `pid` -> `CLONE_NEWPID`
- `network` -> `CLONE_NEWNET`
- `mount` -> `CLONE_NEWNS`
- `uts` -> `CLONE_NEWUTS` (isolate namespaces)
- `ipc` -> `CLONE_NEWIPC`
- `user` -> `CLONE_NEWUSER`
- `cgroup` -> `CLONE_NEWCGROUP`
- `time` -> `CLONE_NEWTIME`

 In Linux, namespaces are manipulated using three system calls:
 `clone(2)`: Spawns a new process in newly created namespaces.
 `unshare(2)`: Moves the current thread into new namespaces.
 `setns(2)`: Joins the current thread to existing namespaces.

 Also Linux Namespaces are thread scoped
 Because Go is a multi-threaded runtime, calling  setns(2)  or  unshare(2)  directly within a running Go process is not supported/reliable, since those system calls apply only to the
  calling thread. The project solves this using a C dynamic linker constructor ( __attribute__((constructor)) ) in a c file via cgo. The package is blank-imported in
  main.go so it executes before the Go runtime initializes.

So at this exact moment:
  1. The process is single-threaded (only the main thread exists).
  2. It is fully legal to call  setns  to join any namespace (including the mount namespace).
  3. We can safely call  chroot  and  chdir  to pivot the filesystem.
