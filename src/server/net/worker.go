package net

import (
	"net"
	"strings"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

type Worker struct {
	ID   int
	Enge *Engine
}

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

func (engine *Engine) findGroupHandlers(path string) (handlers []HandlerFunc, find bool) {
	for p, gp := range engine.groupHandlers {
		if strings.Contains(path, p) {
			h, ok := gp.handlers[path]
			if !ok {
				return
			}
			handlers = append(gp.midHandlers, h)
			find = true
			break
		}
	}
	return
}
