package ip_repo

import "fmt"

func (r IPRepo) ReleaseIPAddress(addressID string) error {
	addresses, err := r.GetAddresses()
	if err != nil {
		return err
	}

	var address *IPAddress
	for _, addr := range addresses {
		if addr.ID == addressID {
			address = &addr
			break
		}
	}
	if address == nil {
		return fmt.Errorf("address with ID %s not found", addressID)
	}

	address.Using = false

	if err := r.WriteToSheet(*address); err != nil {
		return err
	}

	return nil
}
