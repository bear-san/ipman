package ip_repo

import (
	"fmt"
)

func (r IPRepo) AssignIPAddress(addressType string, description string) (*IPAddress, error) {
	if addressType != IP_ADDRESS_TYPE_LOCAL && addressType != IP_ADDRESS_TYPE_GLOBAL {
		return nil, fmt.Errorf("invalid IP Address Type: %s", addressType)
	}

	addresses, err := r.GetAddresses()
	if err != nil {
		return nil, err
	}

	var nominatedAddress *IPAddress
	for _, address := range addresses {
		if !address.AutoAssignEnabled || address.Using || address.AddressType != addressType {
			continue
		}

		nominatedAddress = &address
		break
	}

	if nominatedAddress == nil {
		return nil, fmt.Errorf("%s IP address is out of stock", addressType)
	}

	nominatedAddress.Using = true
	nominatedAddress.Description = description
	if err := r.WriteToSheet(*nominatedAddress); err != nil {
		return nil, err
	}

	return nominatedAddress, nil
}
