package api

import (
	"context"
	"encoding/json"
	"net"
	"time"
)

func udpJSON(ctx context.Context, address string, payload any) (map[string]any, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	dialer := net.Dialer{Timeout: time.Second}
	conn, err := dialer.DialContext(ctx, "udp", address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := conn.Write(body); err != nil {
		return nil, err
	}

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	var response map[string]any
	if err := json.Unmarshal(buf[:n], &response); err != nil {
		return nil, err
	}
	return response, nil
}
