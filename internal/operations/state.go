package operations

import (
	"encoding/json"
	"fmt"

	"github.com/Abhijithmns/hollow/internal/container"
)

// state <container-id>

type StateOpts struct {
	ID string
}

func State(opts *StateOpts) (string, error) {
	cont , err := container.Load(opts.ID)
	if err != nil {
		return "", fmt.Errorf("load container: %w", err)
	}
	
	state, err := json.Marshal(cont.State)
	if err != nil {
		return "", fmt.Errorf("marshall container state: %w", err)
	}

	return string(state), nil
}
