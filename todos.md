# hollow - an oci container runtime in go

a simple oci-compliant container runtime written in go to learn how containers work :)

### runtime
- [x] build the cli interface
- [x] parse an oci bundle (`config.json`)
- [x] create and persist container state
- [x] load existing containers
- [x] implement the container lifecycle (`create`, `start`, `state`, `kill`, `delete`)
- [x] execute oci lifecycle hooks *(partial)*

### process isolation
- [x] implement linux namespaces
- [x] configure the root filesystem
- [x] configure mounts
- [x] set the hostname and domain name

### resource management
- [x] implement cgroups v2 (memory, pids limits — verified via dmesg fork-rejection under real limits)

### security
- [ ] configure uid/gid mappings
- [ ] configure linux capabilities
- [ ] configure kernel parameters (`sysctl`)
- [ ] configure resource limits (`rlimits`)

### runtime features
- [ ] wire up a console/tty
- [ ] write the container pid file
- [ ] adjust the oom score
- [ ] handle pid 1 / zombie reaping (init process for multi-process containers)
- [ ] implement the oci `features` api

### if everything works
- [ ] pass the oci runtime specification test suite
- [ ] use hollow with docker/containerd
- [ ] run doom
