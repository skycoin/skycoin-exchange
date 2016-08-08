package net

import (
	"net"

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
	w.c = c

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

		if h, ok := engine.handlerFunc[r.GetPath()]; ok {
			context.handlers = append(engine.handlers, h)
		} else {
			logger.Error("no handler for path: %s", r.GetPath())
			res := pp.MakeErrResWithCode(pp.ErrCode_ServerError)
			context.JSON(res)
			return
		}

		context.handlers[0](&context)
	}
}
