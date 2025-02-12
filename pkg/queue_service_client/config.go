package queue_service_client

import (
	"fmt"
	"os"
)

type Config struct {
	Host string
	Port string
}

func (c *Config) EnvParse(prefix string) error {
	const errorMsg = "from queue client config error: %w"

	var exists bool

	nameAddressEnv := prefix + "HOST"
	c.Host, exists = os.LookupEnv(nameAddressEnv)
	if !exists {
		return fmt.Errorf(errorMsg, fmt.Errorf("env variable %s must be set", nameAddressEnv))
	}

	namePortEnv := prefix + "PORT"
	c.Port, exists = os.LookupEnv(namePortEnv)
	if !exists {
		return fmt.Errorf(errorMsg, fmt.Errorf("env variable %s must be set", namePortEnv))
	}

	return nil
}
