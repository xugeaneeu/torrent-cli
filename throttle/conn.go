package throttle

import "net"

// NewConn возвращает обёртку net.Conn, которая при каждом Write
// ждёт `Take(len(p))`.
func NewConn(inner net.Conn) net.Conn {
	return &throttledConn{Conn: inner}
}

type throttledConn struct {
	net.Conn
}

func (t *throttledConn) Write(p []byte) (int, error) {
	// оттормозить согласно baket
	Take(len(p))
	return t.Conn.Write(p)
}
