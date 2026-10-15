package handlers

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
)

func TestValidateImageInfoRegistry(t *testing.T) {
	for _, host := range []string{"ghcr.io", "registry.k8s.io", "docker.io"} {
		ref, err := name.ParseReference(host + "/team/image:latest")
		if err != nil {
			t.Fatal(err)
		}
		if err := validateImageInfoRegistry(ref); err != nil {
			t.Errorf("host %q rejected: %v", host, err)
		}
	}
	ref, err := name.ParseReference("attacker.example/team/image:latest")
	if err != nil {
		t.Fatal(err)
	}
	if err := validateImageInfoRegistry(ref); err == nil {
		t.Fatal("unconfigured registry host accepted")
	}
}

func TestRegistryHostGuardRejectsRedirectToUnknownHost(t *testing.T) {
	base := roundTripperFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("base transport called for an unknown host")
		return nil, nil
	})
	request := httptest.NewRequest(http.MethodGet, "https://169.254.169.254/latest/meta-data", nil)
	if _, err := (registryHostGuard{base: base}).RoundTrip(request); err == nil {
		t.Fatal("unknown redirect host was accepted")
	}
}

func TestUnsafeRegistryIP(t *testing.T) {
	tests := []struct {
		ip           string
		allowPrivate bool
		unsafe       bool
	}{
		{ip: "127.0.0.1", allowPrivate: true, unsafe: true},
		{ip: "10.0.0.2", unsafe: true},
		{ip: "10.0.0.2", allowPrivate: true},
		{ip: "169.254.169.254", allowPrivate: true, unsafe: true},
		{ip: "8.8.8.8", unsafe: false},
	}
	for _, tt := range tests {
		if got := unsafeRegistryIP(net.ParseIP(tt.ip), tt.allowPrivate); got != tt.unsafe {
			t.Errorf("unsafeRegistryIP(%s, %t) = %t, want %t", tt.ip, tt.allowPrivate, got, tt.unsafe)
		}
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
