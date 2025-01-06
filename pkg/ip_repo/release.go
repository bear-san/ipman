package ip_repo

func (r IPRepo) ReleaseIPAddress(address IPAddress) error {
	address.Using = false
	if err := r.WriteToSheet(address); err != nil {
		return err
	}

	return nil
}
