package ipcheck

import (
	"fmt"
	"net/netip"
	"strings"
)

type IChecker interface {
	Check(ip string) (bool, error)
}

func New() *Checker {
	return &Checker{
		ips: make(map[netip.Addr]struct{}),
	}
}

func TrueChecker() IChecker {
	return defaultCk{v: true}
}

func FalseChecker() IChecker {
	return defaultCk{v: false}
}

type defaultCk struct {
	v bool
}

func (c defaultCk) Check(ip string) (bool, error) {
	return c.v, nil
}

// Checker stores IP addresses and CIDR prefixes for membership checks.
type Checker struct {
	ips      map[netip.Addr]struct{}
	prefixes []netip.Prefix
}

// Add registers an IP address or CIDR prefix. Entry must be a valid IPv4/IPv6
// address or prefix in CIDR notation.
func (c *Checker) Add(entry string) error {
	trimmed := strings.TrimSpace(entry)
	if trimmed == "" {
		return fmt.Errorf("ipcheck: empty entry")
	}

	if strings.Contains(trimmed, "/") {
		return c.addPrefix(trimmed)
	}
	return c.addAddr(trimmed)
}

func (c *Checker) AddList(lst ...string) error {
	for _, item := range lst {
		if err := c.Add(item); err != nil {
			return fmt.Errorf("add entry:%s failed, err:%w", item, err)
		}
	}
	return nil
}

func (c *Checker) addAddr(input string) error {
	addr, err := netip.ParseAddr(input)
	if err != nil {
		return fmt.Errorf("ipcheck: invalid IP %q: %w", input, err)
	}

	if c.ips == nil {
		c.ips = make(map[netip.Addr]struct{})
	}
	c.ips[addr] = struct{}{}
	return nil
}

func (c *Checker) addPrefix(input string) error {
	prefix, err := netip.ParsePrefix(input)
	if err != nil {
		return fmt.Errorf("ipcheck: invalid CIDR %q: %w", input, err)
	}
	if !prefix.IsValid() {
		return fmt.Errorf("ipcheck: invalid CIDR %q", input)
	}

	c.prefixes = append(c.prefixes, prefix.Masked())
	return nil
}

// Check reports whether the supplied IP address matches any registered IP or
// CIDR prefix.
func (c *Checker) Check(ip string) (bool, error) {
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil {
		return false, fmt.Errorf("ipcheck: invalid IP %q: %w", ip, err)
	}

	if _, ok := c.ips[addr]; ok {
		return true, nil
	}

	for _, prefix := range c.prefixes {
		if prefix.Contains(addr) {
			return true, nil
		}
	}

	return false, nil
}
