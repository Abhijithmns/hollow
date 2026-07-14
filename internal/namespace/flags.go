package namespace

import (
	"fmt"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func CloneFlags(spec *specs.Spec) (uintptr, error) {
	// no namespaces in config.json
	if spec.Linux == nil || spec == nil {
		return 0, fmt.Errorf("no namespaces mentioned in 'config.json' ")
	}
	var flags uintptr
	for _,ns := range spec.Linux.Namespaces {
		if ns.Path == "" {
			// do something
			continue
		}
		flag, ok := NamespaceFlags[ns.Type]
		if !ok {
			// unknown namespace 
			// the OCI spec defines all the valid types, so just ignore
			return 0, fmt.Errorf("unsupported namespace type: %s", ns.Type)
		}
		flags |= flag
	}
	return flags, nil
}
