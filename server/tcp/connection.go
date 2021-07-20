package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	globalTcpConnectionId uint64

// noDeadline            = time.Time{}
)

type TcpConnection struct {
	id        uint64
	headSize  uint32
	conn      net.Conn
	readBuf   []byte
	writeBuf  []byte
	sendMutex sync.Mutex
	recvMutex sync.Mutex

	// About close
	closeFlag int32
}

func newTcpConnection(conn net.Conn, headSize uint32) *TcpConnection {
	return &TcpConnection{
		id:       atomic.AddUint64(&globalTcpConnectionId, 1),
		headSize: headSize,
		conn:     conn,
		readBuf:  make([]byte, headSize),
		writeBuf: make([]byte, headSize),
	}
}

func (c *TcpConnection) SetKeepAlive(period time.Duration) (err error) {
	if tc, ok := c.conn.(*net.TCPConn); ok {
		if err = tc.SetKeepAlive(true); err != nil {
			return
		}
		if err = tc.SetKeepAlivePeriod(period); err != nil {
			return
		}
	}
	return
}

func (c *TcpConnection) Receive() (msg []byte, err error) {
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()
	_, err = io.ReadFull(c.conn, c.readBuf[:c.headSize])
	if err != nil {
		err = fmt.Errorf("read head error: %s", err.Error())
		return
	}
	n := binary.BigEndian.Uint32(c.readBuf[:c.headSize])
	msg = make([]byte, n)
	_, err = io.ReadFull(c.conn, msg[:n])
	if err != nil {
		err = fmt.Errorf("read body error: %s", err.Error())
	}
	return
}

func (c *TcpConnection) Send(msg []byte) (err error) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	binary.BigEndian.PutUint32(c.writeBuf[:c.headSize], uint32(len(msg)))
	_, err = c.conn.Write(c.writeBuf[:c.headSize])
	if err != nil {
		return fmt.Errorf("write head error: %s", err.Error())
	}
	_, err = c.conn.Write(msg)
	return
}

func (c *TcpConnection) Close() {
	if atomic.CompareAndSwapInt32(&c.closeFlag, 0, 1) {
		c.conn.Close()
	}
}

func (c *TcpConnection) readPacket() {
	if atomic.CompareAndSwapInt32(&c.closeFlag, 0, 1) {
		c.conn.Close()
	}
}

func (c *TcpConnection) writePacket(msg []byte) {
	if atomic.CompareAndSwapInt32(&c.closeFlag, 0, 1) {
		c.conn.Close()
	}
}
