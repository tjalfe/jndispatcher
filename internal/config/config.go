package config

import (
	"fmt"

	"github.com/tjalfe/jndispatcher/internal/types"
)

func ReadConfig() (types.Config, error) {
	config, err := readConfig()
	if err != nil {
		return types.Config{}, fmt.Errorf("error reading config: %w", err)
	}
	return config, nil
}
