package handlers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"hubproxy/config"
	"hubproxy/utils"

	"github.com/google/go-containerregistry/pkg/name"
)

type registryHostGuard struct {
	base http.RoundTripper
}

func (g registryHostGuard) RoundTrip(request *http.Request) (*http.Response, error) {
	if !allowedRegistryAddress(request.URL.Host) {
		return nil, fmt.Errorf("registry request host is not configured")
	}
	return g.base.RoundTrip(request)
}

func validateImageInfoRegistry(ref name.Reference) error {
	if !allowedRegistryAddress(ref.Context().RegistryStr()) {
		return fmt.Errorf("registry host is not enabled")
	}
	return nil
}

func imageInfoTransport() (http.RoundTripper, error) {
	client := utils.GetGlobalHTTPClient()
	if client == nil {
		return nil, fmt.Errorf("registry client is unavailable")
	}
	global := client.Transport
	transport, ok := global.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("registry transport is unavailable")
	}
	guarded := transport.Clone()
	// The metadata endpoint is small; direct connections let us validate DNS
	// results before dialing. Other mirror traffic keeps its configured proxy.
	guarded.Proxy = nil
	guarded.DialContext = safeRegistryDialContext
	return registryHostGuard{base: guarded}, nil
}

func safeRegistryDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if !allowedRegistryAddress(address) {
		return nil, fmt.Errorf("registry dial host is not configured")
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("invalid registry address")
	}
	allowPrivate := registryAllowsPrivateNetwork(address)
	ip := net.ParseIP(host)
	if ip != nil {
		if unsafeRegistryIP(ip, allowPrivate) {
			return nil, fmt.Errorf("registry address resolves to a restricted network")
		}
		return (&net.Dialer{}).DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
	}
	resolved, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("resolve registry host: %w", err)
	}
	if len(resolved) == 0 {
		return nil, fmt.Errorf("registry host has no address")
	}
	for _, candidate := range resolved {
		if unsafeRegistryIP(candidate.IP, allowPrivate) {
			return nil, fmt.Errorf("registry host resolves to a restricted network")
		}
	}
	return (&net.Dialer{}).DialContext(ctx, network, net.JoinHostPort(resolved[0].IP.String(), port))
}

func unsafeRegistryIP(ip net.IP, allowPrivate bool) bool {
	if ip == nil {
		return true
	}
	ip = ip.To16()
	if ip == nil {
		return true
	}
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}
	if ip.IsUnspecified() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return true
	}
	if ip.IsPrivate() {
		return !allowPrivate
	}
	if ip.To4() != nil && ip.Mask(net.CIDRMask(10, 32)).Equal(net.ParseIP("100.64.0.0").To4()) {
		return true
	}
	return !ip.IsGlobalUnicast()
}

func allowedRegistryAddress(raw string) bool {
	address := normalizeRegistryAddress(raw)
	switch address {
	case "docker.io", "index.docker.io", "registry-1.docker.io", "auth.docker.io":
		return true
	}
	cfg := config.GetConfig()
	for registry, mapping := range cfg.Registries {
		if !mapping.Enabled {
			continue
		}
		if address == normalizeRegistryAddress(registry) || address == normalizeRegistryAddress(mapping.Upstream) || address == normalizeRegistryAddress(mapping.AuthHost) {
			return true
		}
	}
	return false
}

func registryAllowsPrivateNetwork(raw string) bool {
	address := normalizeRegistryAddress(raw)
	for registry, mapping := range config.GetConfig().Registries {
		if !mapping.Enabled || !mapping.AllowPrivateNetwork {
			continue
		}
		if address == normalizeRegistryAddress(registry) || address == normalizeRegistryAddress(mapping.Upstream) {
			return true
		}
	}
	return false
}

func normalizeRegistryAddress(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return ""
	}
	if parsed, err := url.Parse(raw); err == nil && parsed.Host != "" {
		raw = parsed.Host
	} else if slash := strings.IndexByte(raw, '/'); slash >= 0 {
		raw = raw[:slash]
	}
	return strings.TrimSuffix(raw, ".")
}
