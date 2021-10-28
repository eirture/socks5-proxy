package socks5

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

func Auth(client net.Conn) (err error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return fmt.Errorf("reading header: %w", err)
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}

	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return fmt.Errorf("reading methods: %w", err)
	}

	n, err = client.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return fmt.Errorf("write resp err: %w", err)
	}

	return
}

func Connect(client net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(client, buf[:4])
	if n != 4 {
		return nil, fmt.Errorf("read header: %w", err)
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid ver/cmd")
	}

	addr := ""
	switch atyp {
	case 1:
		n, err := io.ReadFull(client, buf[:4])
		if n != 4 {
			return nil, fmt.Errorf("invalid IPv4: %w", err)
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 3:
		n, err = io.ReadFull(client, buf[:1])
		if n != 1 {
			return nil, fmt.Errorf("invalid hostname: %w", err)
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(client, buf[:addrLen])
		if n != addrLen {
			return nil, fmt.Errorf("invalid hostname: %w", err)
		}
		addr = string(buf[:addrLen])
	case 4:
		return nil, errors.New("IPv6: no supported yet")
	default:
		return nil, errors.New("invalid atyp")
	}

	n, err = io.ReadFull(client, buf[:2])
	if n != 2 {
		return nil, fmt.Errorf("read port: %w", err)
	}
	port := binary.BigEndian.Uint16(buf[:2])

	destAddrPort := fmt.Sprintf("%s:%d", addr, port)
	dest, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		return nil, fmt.Errorf("dial dst: %w", dest)
	}

	n, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dest.Close()
		return nil, fmt.Errorf("write rsp: %w", err)
	}
	return dest, nil
}
