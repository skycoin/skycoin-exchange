package net

import (
	"net"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

// Worker tcp handler.
type Worker struct {
	ID   int
	Enge *Engine
}

// Start start the worker.
func (wk *Worker) Start(quit chan bool) {
	go func() {
		for {
			select {
			case c := <-wk.Enge.connPool:
				process(wk.ID, c, wk.Enge)
			case <-quit:
				return
			}
		}
	}()
}

// process handle the incoming connection, will read request from conn, setup the middle,
// and dispatch the request.
func process(id int, c net.Conn, engine *Engine) {
	logger.Debug("[%d] working", id)
	r := &Request{}
	w := &NetResponse{c: c}

	defer func() {
		c.Close()
		logger.Debug("[%d] worker done", id)
	}()

	var err error
	for {
		r.Reset()
		if err = r.Read(c); err != nil {
			return
		}
		context := Context{
			Request: r,
			Resp:    w,
			Data:    make(map[string]interface{}),
		}

		// check if the path belongs to group.
		hds, find := engine.findGroupHandlers(r.GetPath())
		if find {
			context.handlers = append(engine.handlers, hds...)
		} else {
			if h, ok := engine.handlerFunc[r.GetPath()]; ok {
				context.handlers = append(engine.handlers, h)
			} else {
				logger.Error("no handler for path: %s", r.GetPath())
				res := pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				context.JSON(res)
				return
			}
		}

		context.handlers[0](&context)
	}
}

// findGroupHandlers find group of specific path.
func (engine *Engine) findGroupHandlers(path string) (handlers []HandlerFunc, find bool) {
	for p, gp := range engine.groupHandlers {
		if strings.Contains(path, p) {
			h, ok := gp.regHandlers[path]
			if !ok {
				return
			}
			handlers = append(gp.preHandlers, h)
			find = true
			break
		}
	}
	return
}
