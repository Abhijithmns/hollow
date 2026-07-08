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

## How do we execute hooks?

- Create a new Context to use for the execution of the hook.
- If the hook has a Timeout specified, we update the Context with a timeout.
- Lookup the path to the specified executable.
- Create a new Cmd with the context, executable and path to it.
- Apply the Args and Env to the command.
- Run the command which executes the hook.

- when the runtime execs a hook, it is launching another process and waiting for it to finish before Continuing with the container creation.
- if it doesnt wait for a hook to finish, it could start the container before the essential tasks were compleate
