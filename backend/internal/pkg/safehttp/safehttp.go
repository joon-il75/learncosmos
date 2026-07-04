package safehttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidURL      = errors.New("safehttp: invalid url")
	ErrBlockedHost     = errors.New("safehttp: blocked host")
	ErrUnsupported     = errors.New("safehttp: unsupported scheme")
	ErrResolveHost     = errors.New("safehttp: resolve host")
	ErrTooManyRedirect = errors.New("safehttp: too many redirects")
)

var blockedCIDRs = mustParsePrefixes([]string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"192.168.0.0/16",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"::/128",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
	"ff00::/8",
})

func mustParsePrefixes(values []string) []netip.Prefix {
	prefixes := make([]netip.Prefix, 0, len(values))
	for _, value := range values {
		prefix, err := netip.ParsePrefix(value)
		if err != nil {
			panic(err)
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes
}

func ValidateHTTPURL(ctx context.Context, rawURL string) (*url.URL, error) {
	parsed, err := parseHTTPURL(rawURL)
	if err != nil {
		return nil, err
	}
	if _, err := resolveAllowedAddrs(ctx, parsed.Hostname()); err != nil {
		return nil, err
	}
	return parsed, nil
}

func parseHTTPURL(rawURL string) (*url.URL, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return nil, ErrInvalidURL
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed == nil || parsed.Host == "" || parsed.Hostname() == "" {
		return nil, ErrInvalidURL
	}
	if parsed.User != nil {
		return nil, ErrInvalidURL
	}
	scheme := strings.ToLower(strings.TrimSpace(parsed.Scheme))
	if scheme != "http" && scheme != "https" {
		return nil, ErrUnsupported
	}
	return parsed, nil
}

func resolveAllowedAddrs(ctx context.Context, host string) ([]netip.Addr, error) {
	host = strings.Trim(strings.TrimSpace(host), "[]")
	if host == "" {
		return nil, ErrInvalidURL
	}
	lowerHost := strings.ToLower(host)
	if lowerHost == "localhost" || strings.HasSuffix(lowerHost, ".localhost") {
		return nil, ErrBlockedHost
	}
	if addr, err := netip.ParseAddr(host); err == nil {
		addr = addr.Unmap()
		if isBlockedAddr(addr) {
			return nil, ErrBlockedHost
		}
		return []netip.Addr{addr}, nil
	}
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil || len(ips) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrResolveHost, host)
	}
	addrs := make([]netip.Addr, 0, len(ips))
	for _, ip := range ips {
		addr, ok := netip.AddrFromSlice(ip.IP)
		if !ok {
			return nil, ErrBlockedHost
		}
		addr = addr.Unmap()
		if isBlockedAddr(addr) {
			return nil, ErrBlockedHost
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func isBlockedAddr(addr netip.Addr) bool {
	addr = addr.Unmap()
	if !addr.IsValid() {
		return true
	}
	for _, prefix := range blockedCIDRs {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

type ClientConfig struct {
	Timeout      time.Duration
	MaxRedirects int
}

func NewClient(cfg ClientConfig) *http.Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	maxRedirects := cfg.MaxRedirects
	if maxRedirects <= 0 {
		maxRedirects = 5
	}
	dialer := &net.Dialer{Timeout: timeout}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				return nil, err
			}
			addrs, err := resolveAllowedAddrs(ctx, host)
			if err != nil {
				return nil, err
			}
			var lastErr error
			for _, addr := range addrs {
				conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(addr.String(), port))
				if err == nil {
					return conn, nil
				}
				lastErr = err
			}
			return nil, lastErr
		},
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return ErrTooManyRedirect
			}
			_, err := ValidateHTTPURL(req.Context(), req.URL.String())
			return err
		},
	}
}

func DefaultClient() *http.Client {
	return NewClient(ClientConfig{})
}
