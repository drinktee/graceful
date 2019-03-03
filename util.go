package graceful

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
)

const (
	// unixProtocol is the network protocol of unix socket.
	unixProtocol = "unix"
	tcpProtocol  = "tcp"
)

func CreateListenerFile(endpoint string) (net.Listener, *os.File, error) {
	protocol, addr, err := parseEndpointWithFallbackProtocol(endpoint, tcpProtocol)
	if err != nil {
		return nil, nil, err
	}
	if protocol == unixProtocol {
		// err := os.Remove(addr)
		// if err != nil && !os.IsNotExist(err) {
		// 	return nil, nil, fmt.Errorf("failed to remove socket file %q: %v", addr, err)
		// }
		unixAddr, err := net.ResolveUnixAddr(unixProtocol, addr)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to resolve: %v\n", err)

		}
		listener, err := net.ListenUnix(unixProtocol, unixAddr)
		if err != nil {
			return nil, nil, err
		}
		f, err := listener.File()
		if err != nil {
			return nil, nil, fmt.Errorf("ListenUnix failed to retreive fd for: %s (%s)", unixAddr, err)
		}
		return listener, f, nil
	}

	if protocol == tcpProtocol {
		tcpAddr, err := net.ResolveTCPAddr(tcpProtocol, addr)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to resolve: %v\n", err)
		}
		listener, err := net.ListenTCP(tcpProtocol, tcpAddr)
		if err != nil {
			return nil, nil, err
		}
		f, err := listener.File()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to retreive fd for: %s (%s)", addr, err)
		}
		return listener, f, nil
	}
	return nil, nil, fmt.Errorf("only support unix socket or tcp endpoint")
}

func parseEndpointWithFallbackProtocol(endpoint string, fallbackProtocol string) (protocol string, addr string, err error) {
	if protocol, addr, err = parseEndpoint(endpoint); err != nil && protocol == "" {
		fallbackEndpoint := fallbackProtocol + "://" + endpoint
		protocol, addr, err = parseEndpoint(fallbackEndpoint)
		if err == nil {
			log.Printf("Using %q as endpoint is deprecated, please consider using full url format %q.", endpoint, fallbackEndpoint)
		}
	}
	return
}

func parseEndpoint(endpoint string) (string, string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", "", err
	}

	switch u.Scheme {
	case "tcp":
		return "tcp", u.Host, nil

	case "unix":
		return "unix", u.Path, nil

	case "":
		return "", "", fmt.Errorf("Using %q as endpoint is deprecated, please consider using full url format", endpoint)

	default:
		return u.Scheme, "", fmt.Errorf("protocol %q not supported", u.Scheme)
	}
}
