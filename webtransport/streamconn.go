package webtransport

import (
	"net"

	"github.com/quic-go/webtransport-go"
)

var _ net.Conn = (*StreamConn)(nil)

// StreamsConn wraps webtransport.Stream into net.Conn
type StreamConn struct {
	webtransport.Stream
	Session *webtransport.Session
}

func NewAcceptedConn(s webtransport.Stream, ssn *webtransport.Session) net.Conn {
	WebtransportConnsAccepted.Add(1)
	return &StreamConn{s, ssn}
}

func NewOpenedConn(s webtransport.Stream, ssn *webtransport.Session) net.Conn {
	WebtransportConnsOpened.Add(1)
	return &StreamConn{s, ssn}
}

// Close calls CancelRead to avoid memory leak, see
// - https://github.com/quic-go/quic-go/issues/3558
// - https://pkg.go.dev/github.com/quic-go/webtransport-go#Stream
func (sc *StreamConn) Close() error {
	sc.Stream.CancelRead(CancelRead)
	WebtransportConnsClosed.Add(1)
	return sc.Stream.Close()
}

var CancelRead webtransport.StreamErrorCode = 3558
var CancelWrite webtransport.StreamErrorCode = 3559

// LocalAddr is required to impl net.Conn
func (sc *StreamConn) LocalAddr() net.Addr { return sc.Session.LocalAddr() }

// RemoteAddr is required to impl net.Conn
func (sc *StreamConn) RemoteAddr() net.Addr { return sc.Session.RemoteAddr() }

func (sc *StreamConn) CloseRead() error {
	sc.Stream.CancelRead(CancelRead)
	return nil
}

func (sc *StreamConn) CloseWrite() error {
	sc.Stream.CancelWrite(CancelWrite)
	return nil
}