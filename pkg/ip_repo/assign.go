package ip_repo

import (
	"fmt"
)

func (r IPRepo) AssignIPAddress(addressType string) (*IPAddress, error) {
	if addressType != IP_ADDRESS_TYPE_LOCAL && addressType != IP_ADDRESS_TYPE_GLOBAL {
		return nil, fmt.Errorf("invalid IP Address Type: %s", addressType)
	}

	return nil, nil
}
