package sknet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/pp"

	"gopkg.in/op/go-logging.v1"
)

var (
	logger    = logging.MustGetLogger("exchange.net")
	QueueSize = 1000
)

// HandlerFunc important element for implementing the middleware function.
type HandlerFunc func(c *Context)

// Engine is the core of the net package.
type Engine struct {
	handlerFunc   map[string]HandlerFunc
	handlers      []HandlerFunc
	groupHandlers map[string]*Group
	connPool      chan net.Conn
}

// New create an engine.
func New(quit chan bool) *Engine {
	e := &Engine{
		handlerFunc:   make(map[string]HandlerFunc),
		groupHandlers: make(map[string]*Group),
		connPool:      make(chan net.Conn, QueueSize),
	}

	for i := 0; i < QueueSize; i++ {
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

	root := ps[0]
	for p := range engine.groupHandlers {
		if strings.HasPrefix(p, root) {
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
func (engine *Engine) Run(port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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

// Recovery is middleware for catching panic.
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Critical("%s", r)
				debug.PrintStack()
			}
		}()
		c.Next()
	}
}

// Logger middleware
func Logger() HandlerFunc {
	return func(c *Context) {
		logger.Debug("request path:%s", c.Request.GetPath())
		c.Next()
	}
}

func MakeRequest(path string, data interface{}) (*Request, error) {
	r := &Request{}
	r.Path = pp.PtrString(path)
	if data == nil {
		return r, nil
	}
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	r.Data = d[:]
	return r, nil
}

func Write(w io.Writer, v interface{}) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}

	buf := make([]byte, 4+len(d))
	binary.BigEndian.PutUint32(buf[:], uint32(len(d)))
	copy(buf[4:], d)
	if err := binary.Write(w, binary.BigEndian, buf); err != nil {
		return err
	}
	return nil
}

// WriteRsp send request to server, then read response and return.
func WriteRsp(rw io.ReadWriter, v interface{}) (*Response, error) {
	d, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4+len(d))
	binary.BigEndian.PutUint32(buf[:], uint32(len(d)))
	copy(buf[4:], d)
	if err := binary.Write(rw, binary.BigEndian, buf); err != nil {
		return nil, err
	}

	rsp := Response{}
	if err := Read(rw, &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func Read(r io.Reader, v interface{}) error {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return err
	}
	d := make([]byte, len)
	if err := binary.Read(r, binary.BigEndian, &d); err != nil {
		return err
	}
	switch r := v.(type) {
	case *Request:
		if err := json.Unmarshal(d, v); err != nil {
			logger.Error("%s", err)
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
