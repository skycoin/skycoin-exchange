package sknet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	logging "github.com/op/go-logging"
)

var (
	logger               = logging.MustGetLogger("exchange.net")
	queueSize            = 1000
	version       uint32 = 1
	maxReqPkgSize uint32 = 32 * 1024 // set max request package size: 32kb
)

// HandlerFunc important element for implementing the middleware function.
type HandlerFunc func(c *Context) error

// Engine is the core of the net package.
type Engine struct {
	handlers      []HandlerFunc
	handlerFunc   map[string]HandlerFunc
	groupHandlers map[string]*Group
	connPool      chan net.Conn
}

// New create an engine.
func New(seckey string, quit chan bool) *Engine {
	e := &Engine{
		handlerFunc:   make(map[string]HandlerFunc),
		groupHandlers: make(map[string]*Group),
		connPool:      make(chan net.Conn, queueSize),
	}

	e.Use(Authorize(seckey))

	for i := 0; i < queueSize; i++ {
		w := &Worker{
			ID:   i,
			Enge: e,
		}
		w.Start(quit)
	}

	return e
}

// Use add middleware.
func (engine *Engine) Use(handler HandlerFunc) {
	engine.handlers = append(engine.handlers, handler)
}

// Register add request handlers.
func (engine *Engine) Register(path string, handler HandlerFunc) {
	if _, ok := engine.handlerFunc[path]; ok {
		panic(fmt.Sprintf("duplicate router %s", path))
	}
	engine.handlerFunc[path] = handler
}

// Group create request handler group, and bind middleware to this group.
func (engine *Engine) Group(path string, handlers ...HandlerFunc) *Group {
	// check if the group path conflict.
	ps := strings.Split(path, "/")
	if len(ps) == 0 {
		panic("empty path")
	}

	root := ps[1]
	for p := range engine.groupHandlers {
		if strings.HasPrefix(p[1:], root) {
			panic(fmt.Sprintf("conflict group path name:%s with %s", path, p))
		}
	}

	gp := &Group{
		Path:        path,
		preHandlers: handlers,
		regHandlers: make(map[string]HandlerFunc),
	}

	engine.groupHandlers[path] = gp
	return gp
}

// Run start the engine.
func (engine *Engine) Run(ip string, port int) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		logger.Debug("new connection:%s", c.RemoteAddr())
		engine.connPool <- c
	}
}

// Logger middleware
func Logger() HandlerFunc {
	return func(c *Context) error {
		logger.Debug("request path:%s", c.Request.GetPath())
		return c.Next()
	}
}

func Write(w io.Writer, v interface{}) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var data = []interface{}{
		version,        // protocol version
		uint32(len(d)), // payload len
		d,              // payload
	}

	for _, dt := range data {
		if err := binary.Write(w, binary.BigEndian, dt); err != nil {
			return err
		}
	}
	return nil
}

// Read read data from reader and unmarshal to specific struct.
// |  4 bytes | 4 bytes | .........
// |  version |   len   | payload |
func Read(r io.Reader, v interface{}) error {
	// read prefix head version
	var ver uint32
	if err := binary.Read(r, binary.BigEndian, &ver); err != nil {
		return err
	}
	if ver != version {
		return fmt.Errorf("invalid request")
	}

	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return err
	}

	if len > maxReqPkgSize {
		return fmt.Errorf("request data length > %v, check if your request is legal", maxReqPkgSize)
	}

	d := make([]byte, len)
	if err := binary.Read(r, binary.BigEndian, &d); err != nil {
		return err
	}
	switch r := v.(type) {
	case *Request:
		if err := json.Unmarshal(d, v); err != nil {
			logger.Error(err.Error())
			return err
		}
		return nil
	case *Response:
		r.Body = bytes.NewBuffer(d)
		return nil
	default:
		return errors.New("unknow read type")
	}
}
