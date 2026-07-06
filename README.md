# Hollow - a container runtime in Go

## Todos

- [x] Building the CLI interface for a container runtime
- [x] Reading a bundle config and saving a containers state
- [x] Loading a container, getting its state, and deleting it
- [ ] Initialising a container and starting the user process
- [ ] Executing container runtime lifecycle hooks
- [ ] Sending signals to a running container using ‘kill’
- [ ] Setting up the OCI Runtime Spec test suite
- [ ] Isolating the container process using namespaces
- [ ] Managing container resources using cgroups & rlimits
- [ ] Set up the root filesystem (rootfs) of a container
- [ ] Modifying runtime kernel parameters of the container
- [ ] Mounting masked and readonly paths of the container
- [ ] Propagating the rootfs mount of the container
- [ ] Setting the hostname and domain name of the container
- [ ] Setting capabilities and restricting privileges of the container
- [ ] Applying scheduling policies and I/O priority to the container
- [ ] Applying UID/GID and additional GIDs to the container process
- [ ] Wiring up a console to the container
- [ ] Writing the container process PID out to file
- [ ] Adjusting the OOM score of the container process
- [ ] Implementing the ‘features’ API for the container runtime
- [ ] Use hollow with Docker and other tools
- [ ] Run DOOM
