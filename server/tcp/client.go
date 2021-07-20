package tcp

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	defaultMaxRetry       = 3
	defaultConnectTimeout = 3 * time.Second

	clientTcpConnectionPool   = make(map[string][]*TcpConnection)
	clientTcpConnectionPoolMu sync.Mutex
)

type ClientOption func(o *Client)
type Client struct {
	addr     string
	retry    int8
	headSize uint32
	timeout  time.Duration
	conn     *TcpConnection
}

func Addr(addr string) ClientOption {
	return func(s *Client) {
		s.addr = addr
	}
}

func Retry(retry int8) ClientOption {
	return func(s *Client) {
		s.retry = retry
	}
}

func Timeout(timeout time.Duration) ClientOption {
	return func(s *Client) {
		s.timeout = timeout
	}
}

func ClientHeadSize(headSize uint32) ClientOption {
	return func(s *Client) {
		s.headSize = headSize
	}
}

func NewClient(opts ...ClientOption) (client *Client, err error) {
	client = &Client{
		addr:     defaultAddr,
		retry:    int8(defaultMaxRetry),
		headSize: defaultHeadSize,
		timeout:  defaultConnectTimeout,
	}
	for _, o := range opts {
		o(client)
	}

	clientTcpConnectionPoolMu.Lock()
	conn, ok := clientTcpConnectionPool[client.addr]
	if ok {
		connNum := len(conn)
		if connNum > 0 {
			client.conn = conn[0]
			copy(conn, conn[1:])
			conn = conn[:connNum-1]
			clientTcpConnectionPool[client.addr] = conn
			clientTcpConnectionPoolMu.Unlock()
		}
	} else {
		clientTcpConnectionPoolMu.Unlock()
		var c net.Conn
		c, err = net.DialTimeout("tcp", client.addr, client.timeout)
		if err != nil {
			err = fmt.Errorf("dial tcp addr(%s) error: %s", client.addr, err.Error())
			return
		}
		tcpConn := newTcpConnection(c, client.headSize)
		client.conn = tcpConn
	}
	return
}

func (c *Client) Send(data []byte) (res []byte, err error) {
	echan := make(chan error)
	go func() {
		err = c.conn.Send(data)
		echan <- err
	}()

	select {
	case err = <-echan:
		if err != nil {
			fmt.Println("echan err: ", err.Error())
		}
	case <-time.After(5 * time.Second):
		echan <- fmt.Errorf("echan time out")
	}
	if err == nil {
		rchan := make(chan error)
		go func() {
			data, err := c.conn.Receive()
			res = make([]byte, 0, len(data))
			res = append(res, data...)
			rchan <- err
		}()
		select {
		case err := <-rchan:
			if err != nil {
				fmt.Println("rchan err: ", err.Error())
			}
		case <-time.After(5 * time.Second):
			rchan <- fmt.Errorf("rchan time out")
		}
	}
	if err == nil && c.conn != nil {
		clientTcpConnectionPoolMu.Lock()
		clientTcpConnectionPool[c.addr] = append(clientTcpConnectionPool[c.addr], c.conn)
		clientTcpConnectionPoolMu.Unlock()
	}
	return
}
