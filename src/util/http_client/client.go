package http_client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

var blockedCallbackCIDRs = mustParseCallbackCIDRs([]string{
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
	"2001:db8::/32",
})

// ClientFactory is overridden in tests to stub outbound HTTP calls.
var ClientFactory = resty.New

// GetHttpClient 获取请求客户端
func GetHttpClient(proxys ...string) *resty.Client {
	client := ClientFactory()
	// 如果有代理
	if len(proxys) > 0 {
		proxy := proxys[0]
		client.SetProxy(proxy)
	}
	client.SetTimeout(time.Second * 10)
	return client
}

// GetSafeHttpClient returns an outbound client for merchant callbacks. It does
// not use a proxy and re-checks resolved IP addresses at dial/redirect time to
// prevent notify_url from reaching private networks or metadata endpoints.
func GetSafeHttpClient() *resty.Client {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetTransport(&http.Transport{
		Proxy:               nil,
		DialContext:         safeCallbackDialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	})
	client.SetRedirectPolicy(
		resty.FlexibleRedirectPolicy(10),
		resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return ValidateOutboundURL(req.URL.String())
		}),
	)
	return client
}

func ValidateOutboundURL(rawURL string) error {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return fmt.Errorf("invalid callback url: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return errors.New("callback url must include scheme and host")
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
	default:
		return errors.New("callback url scheme must be http or https")
	}
	if u.User != nil {
		return errors.New("callback url must not include user info")
	}
	host := u.Hostname()
	if host == "" {
		return errors.New("callback url host is empty")
	}
	if callbackPrivateNetworksAllowed() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return validateCallbackHost(ctx, host)
}

func safeCallbackDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("invalid callback address %q: %w", address, err)
	}

	ips, err := resolveCallbackIPs(ctx, host)
	if err != nil {
		return nil, err
	}
	if !callbackPrivateNetworksAllowed() {
		for _, ip := range ips {
			if isBlockedCallbackIP(ip) {
				return nil, fmt.Errorf("callback host resolves to blocked address %s", ip.String())
			}
		}
	}
	ip, err := selectDialIP(ips, network)
	if err != nil {
		return nil, err
	}
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	return dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
}

func validateCallbackHost(ctx context.Context, host string) error {
	ips, err := resolveCallbackIPs(ctx, host)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		if isBlockedCallbackIP(ip) {
			return fmt.Errorf("callback host resolves to blocked address %s", ip.String())
		}
	}
	return nil
}

func resolveCallbackIPs(ctx context.Context, host string) ([]net.IP, error) {
	if ip := net.ParseIP(host); ip != nil {
		return []net.IP{ip}, nil
	}
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("resolve callback host %q: %w", host, err)
	}
	ips := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		if addr.IP != nil {
			ips = append(ips, addr.IP)
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("resolve callback host %q: no addresses", host)
	}
	return ips, nil
}

func selectDialIP(ips []net.IP, network string) (net.IP, error) {
	for _, ip := range ips {
		if strings.HasSuffix(network, "4") && ip.To4() == nil {
			continue
		}
		if strings.HasSuffix(network, "6") && ip.To4() != nil {
			continue
		}
		return ip, nil
	}
	return nil, fmt.Errorf("no callback address compatible with network %s", network)
}

func callbackPrivateNetworksAllowed() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("EZPAY_ALLOW_PRIVATE_CALLBACKS")), "true")
}

func isBlockedCallbackIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() {
		return true
	}
	for _, block := range blockedCallbackCIDRs {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func mustParseCallbackCIDRs(raw []string) []*net.IPNet {
	nets := make([]*net.IPNet, 0, len(raw))
	for _, cidr := range raw {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		nets = append(nets, ipNet)
	}
	return nets
}
