- so throughout the life cycle of a container , the runtime will execute hooks according to the phase of the lifecycle.
- hooks are just commands that are run in either the runtime or container namespace

## these are some of them mentioned in the OCI runtime-spec
1. prestart
2. createRuntime 
3. createContainer
4. startContainer
5. poststart
6. poststop

- the `prestart` and `createRuntime` hooks are basically the same
- `prestart` hook is deprecated in favor of all these new `createRuntime`, `createContainer` and `startConatainer`
- but the OCI runtime test suite, Docker and others still use it so ill use it too0

## Some important things to be noted

- When a hook is executed, the runtime is expected to pass a JSON ( !_! ) representation of the current state of the container to stdin
- When a hook fails to execute, it MUST return a error, with the exeption of the `poststop` hook which only logs a warning
