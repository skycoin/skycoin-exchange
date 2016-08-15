package sknet

import (
	"fmt"
	"path/filepath"
)

// Group maintains a group of handlers of specific prefix path,
type Group struct {
	Path        string                 // group path.
	preHandlers []HandlerFunc          // handlers must be executed before registed handlers.
	regHandlers map[string]HandlerFunc // registed handlers.
}

// Register register path and handlers into this group.
// the real path will be prefixed with group path.
func (gp *Group) Register(path string, handler HandlerFunc) {
	if _, ok := gp.regHandlers[path]; ok {
		panic(fmt.Sprintf("dubplicate path:%s in group:%s", path, gp.Path))
	}
	gp.regHandlers[filepath.Join(gp.Path, path)] = handler
}
