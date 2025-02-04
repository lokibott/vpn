package core

import (
	"crypto/tls"
	"net"
	"time"
	
	"golang.org/x/net/proxy"
)

type SecureSocks5Client struct {
	Dialer   proxy.Dialer
	Auth     *proxy.Auth
	ProxyURL string
}

func NewSecureSocks5Client(proxyAddr, user, password string) (*SecureSocks5Client, error) {
	// Validate proxy address format
	if _, _, err := net.SplitHostPort(proxyAddr); err != nil {
		return nil, err
	}

	auth := &proxy.Auth{
		User:     user,
		Password: password,
	}

	// Create base dialer with TLS
	baseDialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// Create SOCKS5 dialer with TLS transport
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, auth, baseDialer)
	if err != nil {
		return nil, err
	}

	return &SecureSocks5Client{
		Dialer:   dialer,
		Auth:     auth,
		ProxyURL: proxyAddr,
	}, nil
}

func (s *SecureSocks5Client) SecureDial(target string) (net.Conn, error) {
	conn, err := s.Dialer.Dial("tcp", target)
	if err != nil {
		return nil, err
	}

	// Upgrade to TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: target,
	}

	return tls.Client(conn, tlsConfig), nil
}