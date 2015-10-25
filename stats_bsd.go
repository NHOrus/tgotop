// +build freebsd darwin

package main

func (m *memData) Update() error {
	return nil
}

func getifnum() (int, error) {
	return 0, nil
}

func (nd *netData) Setup() error {
	return nil
}

func (nd *netData) Update() error {
	return nil
}
