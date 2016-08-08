package net

import (
	"fmt"
	"net"

	"gopkg.in/op/go-logging.v1"
)

var logger = logging.MustGetLogger("exchange.net")

type HandlerFunc func(w ResponseWriter, r *Request) error

type Engine struct {
	handlerFunc   map[string]HandlerFunc
	beforeHandler []HandlerFunc
	afterHandler  []HandlerFunc
	workPool      chan worker
}

func New() *Engine {
	e := &Engine{
		handlerFunc: make(map[string]HandlerFunc),
		workPool:    make(chan worker, 1000),
	}

	for i := 0; i < 1000; i++ {
		w := &NetWorker{
			Pool: e.workPool,
			ID:   i,
		}
		e.workPool <- w
	}
	return e
}

// add middleware
func (engine *Engine) Before(handler HandlerFunc) {
	engine.beforeHandler = append(engine.beforeHandler, handler)
}

func (engine *Engine) After(handler HandlerFunc) {
	engine.afterHandler = append(engine.afterHandler, handler)
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
		w := <-engine.workPool
		w.Work(c, engine)
	}
}
