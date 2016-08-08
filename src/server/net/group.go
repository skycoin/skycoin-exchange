package net

import (
	"fmt"
	"path/filepath"
)

type Group struct {
	Path        string
	midHandlers []HandlerFunc
	handlers    map[string]HandlerFunc
}

func (gp *Group) Register(path string, handler HandlerFunc) {
	if _, ok := gp.handlers[path]; ok {
		panic(fmt.Sprintf("dubplicate path:%s in group:%s", path, gp.Path))
	}
	gp.handlers[filepath.Join(gp.Path, path)] = handler
}
