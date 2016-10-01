// Most of these functions will only work on my machine, or similar
// machines.
package main

import (
	"fmt"
	"net"
)

func main() {
	// find subnet
	a, _ := LocalAddress()
	fmt.Println(a)
	addrs, _ := ListAddresses()
	fmt.Println(addrs)

	ifaces, _ := net.Interfaces()
	fmt.Println(ifaces)
}

// LocalAddress returns local addr.
func LocalAddress() (net.Addr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	addrs, err := ifaces[2].Addrs()
	addr := addrs[0]
	return addr, nil
}

// ListAddress returns slice of address.
func ListAddresses() ([]net.Addr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	addrs, err := ifaces[2].Addrs()
	return addrs, nil
}
