package net

import (
	"fmt"
	"net"

	"gopkg.in/op/go-logging.v1"
)

var (
	logger    = logging.MustGetLogger("exchange.net")
	QueueSize = 1000
)

type HandlerFunc func(c *Context)

type Engine struct {
	handlerFunc map[string]HandlerFunc
	handlers    []HandlerFunc
	connPool    chan net.Conn
}

func New(quit chan bool) *Engine {
	e := &Engine{
		handlerFunc: make(map[string]HandlerFunc),
		connPool:    make(chan net.Conn, QueueSize),
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

// add middleware
func (engine *Engine) Use(handler HandlerFunc) {
	engine.handlers = append(engine.handlers, handler)
}

func (engine *Engine) Register(path string, handler HandlerFunc) {
	if _, ok := engine.handlerFunc[path]; ok {
		panic(fmt.Sprintf("duplicate router %s", path))
	}
	engine.handlerFunc[path] = handler
}

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
