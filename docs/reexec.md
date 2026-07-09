/proc/self/exe -> a sym link to the executable the process is running

- this process will eventually become a container

- The primary reason for doing this is that a container runtime performs two completely different roles during container creations
    The original runtime process parses the OCI configuration, validates arguments, creates namespaces, configures cgroups, prepares networking, and manages the overall life cycle
    At the same time, there must also be a process that actually enters the container environment and eventually becomes the user's application
