# hollow - an oci container runtime in go

a simple oci-compliant container runtime written in go for learning how containers work under the hood.

## roadmap

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
- [ ] set the hostname and domain name

### security
- [ ] configure uid/gid mappings
- [ ] configure linux capabilities
- [ ] configure kernel parameters (`sysctl`)
- [ ] configure resource limits (`rlimits`)

### resource management
- [ ] implement cgroups

### runtime features
- [ ] wire up a console/tty
- [ ] write the container pid file
- [ ] adjust the oom score
- [ ] implement the oci `features` api

### if everything works 
- [ ] pass the oci runtime specification test suite
- [ ] use hollow with docker/containerd
- [ ] run doom
