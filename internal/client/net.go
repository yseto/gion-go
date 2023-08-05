package client

import (
	"errors"
	"fmt"
	"net"
	"syscall"
)

var Error = errors.New("host is blocked")

var (
	v4List []*net.IPNet
	v6List []*net.IPNet
)

func init() {
	// cidr lists refer paranoidhttp

	var v4cidrs = []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"127.0.0.0/8",
		"169.254.0.0/16",
	}
	for _, i := range v4cidrs {
		v4List = append(v4List, mustParseCIDR(i))
	}

	var v6cidrs = []string{
		"0000::/3",
		"4000::/2",
		"8000::/1",
		"2001::/32",
		"2001:10::/28",
		"2001:20::/28",
		"2001:db8::/32",
		"2002::/16",
		"::1/128",
	}
	for _, i := range v6cidrs {
		v6List = append(v6List, mustParseCIDR(i))
	}
}

func mustParseCIDR(cidr string) *net.IPNet {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic("can not parse cidr")
	}

	return ipNet
}

func control(network, address string, c syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("cannot parse address %q: %w", address, err)
	}

	addr := net.ParseIP(host)
	if addr == nil {
		return fmt.Errorf("cannot parse host %q", host)
	}

	if addr.To4() != nil {
		for _, n := range v4List {
			if n.Contains(addr) {
				return Error
			}
		}
	} else if addr.To16() != nil {
		for _, n := range v6List {
			if n.Contains(addr) {
				return Error
			}
		}
	} else {
		return Error
	}

	return nil
}
