package websocketkoffee

import (
	"fmt"
	"koffee/pkg/pool"
	"math/rand"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mailru/easygo/netpoll"
)

// ConnectionData For a new connection or closed connection
type ConnectionData struct {
	Address    string
	ID         uint64
	Headers    map[string][]byte
	Connection net.Conn
}

// Message struct
type Message struct {
	Conn            net.Conn
	Data            []byte
	ErrorConnection error
	ID              uint64
}

// Connection struct
type Connection interface {
	OnClose(func(*ConnectionData))
	OnOpen(func(*ConnectionData))
	OnMessage(func(*Message) error)
}

type connection struct {
	conn      net.Conn
	onClose   func(*ConnectionData)
	onOpen    func(*ConnectionData)
	onMessage func(*Message) error
}

// OnClose
func (c *connection) OnClose(onclose func(*ConnectionData)) {
	c.onClose = onclose
}

// OnOpen
func (c *connection) OnOpen(onopen func(*ConnectionData)) {
	c.onOpen = onopen
}

// OnMessage
func (c *connection) OnMessage(onmessage func(*Message) error) {
	c.onMessage = onmessage
}
func (c *connection) Conn() net.Conn {
	return c.conn
}

func uint64g() uint64 {
	return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}

type uniqueRand struct {
	generated map[uint64]bool
}

func (u *uniqueRand) Int() uint64 {
	i := uint64g()
	return i
}

func iserr(e error) bool {
	if e != nil {
		return true
	}
	return false
}

// Write writes a msg to the client
func Write(c net.Conn, buff []byte) error {
	return wsutil.WriteServerText(c, buff)
}

var poller, _ = netpoll.New(nil)

// !FIXME Too much indirection, I don't think it's cool passing all of these when we had them as values. ATM we are gonna use the anonymous fn
func handleConnection(c connection, ln net.Listener, pool *pool.Pool, u *ws.Upgrader, unique *uniqueRand) func() {
	return func() {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
		_, err = u.Upgrade(conn)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
		c.conn = conn
		connectionData := &ConnectionData{Address: conn.RemoteAddr().String(), ID: unique.Int()}
		c.onOpen(connectionData)
		pool.Schedule(handler(&c, conn, connectionData, pool))
	}
}

// initialize handler for messages from client
func handler(c *connection, conn net.Conn, connection *ConnectionData, pool *pool.Pool) func() {
	return func() {
		desc := netpoll.Must(netpoll.HandleRead(conn))
		poller.Start(desc, func(ev netpoll.Event) {
			// pool.Schedule(func() {
			// Prepare function for closing the connection
			close := func() {
				fmt.Println("Closed messaging")
				conn.Close()
				poller.Stop(desc)
				c.onClose(connection)
			}

			// When ReadHup or Hup received, this means that client has
			// closed at least write end of the connection or connections
			// itself. So we want to stop receive events about such conn.
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				close()
				return
			}
			b, code, e := wsutil.ReadClientData(conn)
			// Error means close connection
			if e != nil {
				close()
				return
			}
			// Continuation means close connection
			if code == ws.OpContinuation {
				close()
				return
			}
			// Close means close connection
			if code == ws.OpClose {
				close()
				return
			}
			// If error reading message close connection
			if len(b) > 0 &&
				iserr(c.onMessage(&Message{Data: b, ErrorConnection: e, ID: connection.ID, Conn: conn})) {
				close()
				return
			}
		})
		// })
	}
}

// Configuration struct for the WS Server
type Configuration struct {
	// Port e. :3333
	Port string
	// Timeout that the server will fire if the worker pool is full (very unlikely) default: 2seconds
	Timeout time.Duration
	// Optional: Function that will fire when there is a header parse on a connection.
	OnHeader func(key, value []byte)
	// If you want the headers of the connection or message stored. Default: false
	StoreHeaders bool
	// How many goroutines do you want for the server MAX Default: 128
	PoolSize int
	// Queue size for the pool thread Default: 32
	QueueSize int
	// Starting threads Default: 64
	Spawn int
	// Max Headers
	MaxHeaders int
}

func (conf *Configuration) defaults() {

	if conf.Port == "" {
		conf.Port = ":3333"
	}

	if conf.Timeout == 0 {
		conf.Timeout = time.Second * 2
	}

	if conf.PoolSize == 0 {
		conf.PoolSize = 128
	}

	if conf.QueueSize == 0 {
		conf.QueueSize = 32
	}

	if conf.Spawn == 0 {
		conf.Spawn = 64
	}

	if conf.MaxHeaders == 0 {
		conf.MaxHeaders = 128
	}

}

// WSKoffeeHandle a handler for the koffee WebSockets
func WSKoffeeHandle(conf Configuration, connectionHandler func(Connection)) error {
	conf.defaults()
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost%s", conf.Port))
	if iserr(err) {
		return err
	}
	desc := netpoll.Must(netpoll.HandleListener(ln, netpoll.EventRead|netpoll.EventOneShot))
	pool := pool.NewPool(conf.PoolSize, conf.QueueSize, conf.Spawn)
	c := connection{onClose: func(*ConnectionData) {}, onOpen: func(*ConnectionData) {}, onMessage: func(*Message) error { return nil }}
	// Initialize onClose, onOpen and onMessage by the user.
	connectionHandler(&c)
	unique := uniqueRand{}
	// Main listener. We get rid of the for loop that affects performance a lot.
	poller.Start(desc, func(e netpoll.Event) {
		errorScheduling := pool.ScheduleTimeout(time.Millisecond, func() {
			conn, err := ln.Accept()
			connectionData := &ConnectionData{Address: conn.RemoteAddr().String(), ID: unique.Int(), Headers: nil, Connection: conn}
			if err != nil {
				fmt.Println(err)
				conn.Close()
				return
			}
			u := ws.Upgrader{}
			if conf.StoreHeaders {
				connectionData.Headers = make(map[string][]byte, conf.MaxHeaders)
			}
			u.OnHeader = func(key, value []byte) (err error) {
				if conf.OnHeader != nil {
					conf.OnHeader(key, value)
				}
				if conf.StoreHeaders {
					if len(connectionData.Headers)+1 > conf.MaxHeaders {
						conn.Close()
						return
					}
					connectionData.Headers[string(key)] = value
				}
				return
			}
			_, err = u.Upgrade(conn)
			if err != nil {
				fmt.Println(err)
				conn.Close()
				return
			}
			c.conn = conn
			c.onOpen(connectionData)
			// add to the pool of threads the handler
			pool.Schedule(handler(&c, conn, connectionData, pool))
		})

		if iserr(errorScheduling) {
			poller.Stop(desc)
			time.Sleep(conf.Timeout)
		}
		poller.Resume(desc)
	})
	return nil
}
