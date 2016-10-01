// Most of these functions will only work on my machine, or similar
// machines.
package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	portscanner "github.com/anvie/port-scanner"
)

func main() {
	// find subnet
	a, _ := LocalAddress()

	addrs, _ := ListAddresses()
	fmt.Println("Addresses:\n", addrs)
	ifaces, _ := net.Interfaces()
	fmt.Println("Interfaces:\n", ifaces)
	fmt.Println("Me:\n", a.String(), ifaces[2])
	local := strings.TrimRight(a.String(), "/24")
	fmt.Println(local)
	//localhost := "localhost"
	//actual := "192.168.1.140"
	ps := portscanner.NewPortScanner(local, 10*time.Second)
	openedPorts := ps.GetOpenedPort(20, 30000)
	fmt.Println(openedPorts)
	for i := 0; i < len(openedPorts); i++ {
		port := openedPorts[i]
		fmt.Print(" ", port, " [open]")
		fmt.Println("  -->  ", ps.DescribePort(port))
	}
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
