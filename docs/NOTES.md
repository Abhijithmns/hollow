## Initialising the container (hollow create <container-id>)
- the runtime creates a container process and applies the required configuration. it then sets up an IPC channel over a unix domain socket so that it can contribute to communicate 
  with the container once the container process is detached

## Starting the container 
- runtime sends a message over the IPC channel to the container process instructing the container to start. The container process applies the remaining configuration, then execs the user process defined in the config

- A container process shouldn't start executing the user's program (/bin/sh, /bin/ash) until the runtime has finished configuring everything. So the runtime and the child process communicate over a Unix domain socket (UDS).

## How do you actually initialise and start a container? (this is what runc does)

1. The user (either directly or via some higher-level runtime tooling, such as Docker) issues the hollow create <container-id> command to the runtime.
2. The runtime does its part of the configuration of the container based on the config.json spec.
3. The runtime creates a Unix domain socket to use for communication between the runtime and the container process.
4. The runtime listens (asynchronously) on the socket for any messages coming from the container.
5. The runtime ‘reexecs’, applying the configuration from step 2 to the new process.
6. The runtime releases the container process, and continues listening on the socket, waiting to receive a ready message from the container.
7. The container process does its part of the configuration based on the config.json spec.
8. The container process sends a ready message to the socket to indicate it’s completed configuring and is ready to receive further commands.
9. The container process listens (synchronously) on the socket for further instructions.
10. The runtime receives the ready message from the socket.
11. The runtime exits, leaving the container process running in the ‘background’.

