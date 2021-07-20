package tcp

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mike955/zrpc/log"
	"github.com/sirupsen/logrus"
)

const (
	defaultAppName         = "tcpServer"
	defaultAddr            = ":9090"
	defaultHeadSize        = 4
	defaultKeepAlivePeriod = 10 * time.Second
)

var (
	defaultMaxConnection = runtime.NumCPU()
)

type ServerOption func(o *Server)
type TcpHandler func(msg []byte) []byte

type Server struct {
	app            string
	addr           string
	headSize       uint32
	maxConnectoion uint64
	stopFlag       int32
	stopWait       sync.WaitGroup

	connections     map[uint64]*TcpConnection
	connectionMutex sync.Mutex
	listener        net.Listener
	handler         TcpHandler
	Logger          *log.Entry
}

func App(app string) ServerOption {
	return func(s *Server) {
		s.app = app
	}
}

func Address(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func MaxConnection(maxConnectoion uint64) ServerOption {
	return func(s *Server) {
		s.maxConnectoion = maxConnectoion
	}
}

func Logger(log *log.Logger) ServerOption {
	return func(s *Server) {
		s.Logger = log.WithFields(map[string]interface{}{"app": s.app})
	}
}

func HeadSize(size uint32) ServerOption {
	return func(s *Server) {
		s.headSize = size
	}
}

func Handler(handler TcpHandler) ServerOption {
	return func(s *Server) {
		s.handler = handler
	}
}

func NewServer(app string, opts ...ServerOption) (srv *Server, err error) {
	srv = &Server{
		app:            defaultAppName,
		addr:           defaultAddr,
		headSize:       defaultHeadSize,
		maxConnectoion: uint64(defaultMaxConnection),
		connections:    map[uint64]*TcpConnection{},
		Logger:         defaultLogger().WithFields(logrus.Fields{"app": app}),
	}
	for _, o := range opts {
		o(srv)
	}
	listener, err := net.Listen("tcp", srv.addr)
	if err != nil {
		panic(err)
	}
	srv.listener = listener
	return
}

func (s *Server) Serve() (err error) {
	s.Logger.Infof("tcp server start at %s", s.addr)
	for {
		c, err := s.listener.Accept()
		if err != nil {
			continue
		}
		connection, err := s.newConnection(c)
		if err != nil {
			continue
		}
		go s.serveConnection(connection)
	}
}

func (s *Server) Stop() (err error) {
	if atomic.CompareAndSwapInt32(&s.stopFlag, 0, 1) {
		s.listener.Close()
		s.closeConnections()
		s.stopWait.Wait()
	} else {
		err = fmt.Errorf("close server error")
	}
	return
}

func (s *Server) newConnection(conn net.Conn) (c *TcpConnection, err error) {
	c = newTcpConnection(conn, s.headSize)
	if err = c.SetKeepAlive(defaultKeepAlivePeriod); err != nil {
		return
	}
	s.addConnection(c)
	return
}

func (s *Server) serveConnection(c *TcpConnection) {
	for {
		req, err := c.Receive()
		if err != nil {
			s.delConnection(c)
			return
		}
		res := s.handler(req)
		err = c.Send(res)
		if err != nil {
			c.Close()
			s.delConnection(c)
			s.Logger.Info()
		}
	}
}

func (s *Server) addConnection(c *TcpConnection) {
	s.connectionMutex.Lock()
	defer s.connectionMutex.Unlock()
	s.connections[c.id] = c
	s.stopWait.Add(1)
}

func (s *Server) delConnection(c *TcpConnection) {
	s.connectionMutex.Lock()
	defer s.connectionMutex.Unlock()
	delete(s.connections, c.id)
	s.stopWait.Done()
}

func (s *Server) closeConnections() {
	for _, connection := range s.connections {
		connection.Close()
	}
}

func defaultLogger() (logger *log.Logger) {
	logger = log.NewLogger()
	return
}
