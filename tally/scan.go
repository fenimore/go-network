package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}

	addrs, err := ifaces[2].Addrs()
	addr := addrs[0]

	ip, ipnet, err := net.ParseCIDR(addr.String())
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		fmt.Println(ip)
		fmt.Println(ipnet.Contains(ip))
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
