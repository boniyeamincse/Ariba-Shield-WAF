package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"sync"
)

// netSender is a syslog transport over UDP, TCP or TLS (P2.29). TCP and TLS
// connections are redialed automatically on write failure.
type netSender struct {
	network string
	addr    string
	tlsCfg  *tls.Config

	mu   sync.Mutex
	conn net.Conn
}

func newNetSender(transport, host, tlsCA string, insecure bool) (*netSender, error) {
	s := &netSender{
		network: transport,
		addr:    host,
	}
	if transport == "tls" {
		cfg := &tls.Config{MinVersion: tls.VersionTLS12}
		if insecure {
			cfg.InsecureSkipVerify = true
		}
		if tlsCA != "" {
			pem, err := os.ReadFile(tlsCA)
			if err != nil {
				return nil, fmt.Errorf("read CA bundle %q: %w", tlsCA, err)
			}
			pool, err := x509.SystemCertPool()
			if err != nil || pool == nil {
				pool = x509.NewCertPool()
			}
			if !pool.AppendCertsFromPEM(pem) {
				return nil, fmt.Errorf("no valid certificates in %q", tlsCA)
			}
			cfg.RootCAs = pool
		}
		s.tlsCfg = cfg
	}
	if err := s.connect(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *netSender) connect() error {
	if s.network == "tls" {
		conn, err := tls.Dial("tcp", s.addr, s.tlsCfg)
		if err != nil {
			return fmt.Errorf("dial tls %s: %w", s.addr, err)
		}
		s.conn = conn
		return nil
	}
	conn, err := net.Dial(s.network, s.addr)
	if err != nil {
		return fmt.Errorf("dial %s %s: %w", s.network, s.addr, err)
	}
	s.conn = conn
	return nil
}

func (s *netSender) send(msg []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		if err := s.connect(); err != nil {
			return err
		}
	}
	if _, err := s.conn.Write(msg); err != nil {
		if s.network == "udp" {
			return err
		}
		// TCP/TLS: the peer may have closed the connection; redial once.
		_ = s.conn.Close()
		s.conn = nil
		if err := s.connect(); err != nil {
			return err
		}
		_, werr := s.conn.Write(msg)
		return werr
	}
	return nil
}

func (s *netSender) close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		err := s.conn.Close()
		s.conn = nil
		return err
	}
	return nil
}
