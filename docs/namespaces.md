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

