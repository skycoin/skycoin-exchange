package net

import (
	"net"

	"github.com/skycoin/skycoin-exchange/src/pp"
)

type worker interface {
	Work(c net.Conn, engine *Engine)
}

type NetWorker struct {
	Pool chan worker
	ID   int
}

func (nw *NetWorker) Work(c net.Conn, engine *Engine) {
	logger.Debug("[%d] working", nw.ID)
	r := &Request{}
	w := &NetResponse{
		c: c,
	}
	// set back the worker to pool
	defer func() {
		c.Close()
		nw.Pool <- nw
		logger.Debug("[%d] worker done", nw.ID)
	}()

	var err error
	for {
		r.Reset()
		if err = r.Read(c); err != nil {
			return
		}
		for _, h := range engine.beforeHandler {
			if err = h(w, r); err != nil {
				logger.Error("%s", err)
				res := pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				w.SendJSON(res)
				return
			}
		}

		h, ok := engine.handlerFunc[r.GetPath()]
		if !ok {
			logger.Error("no handler for router: %s", r.GetPath())
			res := pp.MakeErrResWithCode(pp.ErrCode_ServerError)
			w.SendJSON(res)
			return
		}

		h(w, r)

		for _, h := range engine.afterHandler {
			if err = h(w, r); err != nil {
				logger.Error("%s", err)
				res := pp.MakeErrResWithCode(pp.ErrCode_ServerError)
				w.SendJSON(res)
				return
			}
		}
	}
}
