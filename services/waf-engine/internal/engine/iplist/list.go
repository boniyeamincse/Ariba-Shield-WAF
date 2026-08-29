package iplist

import (
	"net/netip"
	"sync"
)

// List is a thread-safe allow/block IP list (master plan §6.7 IP/CIDR lists).
type List struct {
	mu      sync.RWMutex
	allowed []netip.Prefix
	blocked []netip.Prefix
}

// New creates an empty list.
func New() *List {
	return &List{}
}

// SetAllowed replaces the allowed prefix set atomically.
func (l *List) SetAllowed(prefixes []string) error {
	parsed, err := parsePrefixes(prefixes)
	if err != nil {
		return err
	}
	l.mu.Lock()
	l.allowed = parsed
	l.mu.Unlock()
	return nil
}

// SetBlocked replaces the blocked prefix set atomically.
func (l *List) SetBlocked(prefixes []string) error {
	parsed, err := parsePrefixes(prefixes)
	if err != nil {
		return err
	}
	l.mu.Lock()
	l.blocked = parsed
	l.mu.Unlock()
	return nil
}

// IsAllowed returns true if addr is in the allow list.
func (l *List) IsAllowed(addr string) bool {
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return false
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, p := range l.allowed {
		if p.Contains(ip) {
			return true
		}
	}
	return false
}

// IsBlocked returns true if addr is in the block list (block takes precedence).
func (l *List) IsBlocked(addr string) bool {
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return false
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, p := range l.blocked {
		if p.Contains(ip) {
			return true
		}
	}
	return false
}

func parsePrefixes(prefixes []string) ([]netip.Prefix, error) {
	out := make([]netip.Prefix, 0, len(prefixes))
	for _, s := range prefixes {
		if s == "" {
			continue
		}
		p, err := netip.ParsePrefix(s)
		if err != nil {
			// Try a bare address -> /32 or /128.
			ip, ipErr := netip.ParseAddr(s)
			if ipErr != nil {
				return nil, err
			}
			p = netip.PrefixFrom(ip, ip.BitLen())
		}
		out = append(out, p.Masked())
	}
	return out, nil
}
