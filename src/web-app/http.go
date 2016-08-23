package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/op/go-logging.v1"
)

var (
	logger = logging.MustGetLogger("webapp.gui")
)

const (
	resourceDir = "build/"
	devDir      = "dev/"
	indexPage   = "index.html"
)

// Begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, mux *http.ServeMux) error {
	logger.Info("Starting web interface on http://%s", host)
	logger.Warning("HTTPS not in use!")
	logger.Info("Web resources directory: %s", staticDir)

	appLoc, err := determineResourcePath(staticDir)
	if err != nil {
		return err
	}
	logger.Debug("static dir:%s", appLoc)

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	// mux := http.NewServeMux()
	mux.HandleFunc("/", newIndexHandler(appLoc))

	fileInfos, _ := ioutil.ReadDir(appLoc)
	for _, fileInfo := range fileInfos {
		route := fmt.Sprintf("/%s", fileInfo.Name())
		if fileInfo.IsDir() {
			route = route + "/"
		}
		mux.Handle(route, http.FileServer(http.Dir(appLoc)))
	}

	// Runs http.Serve() in a goroutine
	serve(listener, mux)
	return nil
}

func serve(listener net.Listener, mux *http.ServeMux) {
	go func() {
		if err := http.Serve(listener, mux); err != nil {
			log.Panic(err)
		}
	}()
}

func determineResourcePath(staticDir string) (string, error) {
	//check "dev" directory first
	appLoc := filepath.Join(staticDir, devDir)
	if !strings.HasPrefix(appLoc, "/") {
		// Prepend the binary's directory path if appLoc is relative
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return "", err
		}

		appLoc = filepath.Join(dir, appLoc)
	}

	if _, err := os.Stat(appLoc); os.IsNotExist(err) {
		//check dist directory
		appLoc = filepath.Join(staticDir, resourceDir)
		if !strings.HasPrefix(appLoc, "/") {
			// Prepend the binary's directory path if appLoc is relative
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				return "", err
			}

			appLoc = filepath.Join(dir, appLoc)
		}

		if _, err := os.Stat(appLoc); os.IsNotExist(err) {
			return "", err
		}
	}

	return appLoc, nil
}

// Returns a http.HandlerFunc for index.html, where index.html is in appLoc
func newIndexHandler(appLoc string) http.HandlerFunc {
	// Serves the main page
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("url:", r.URL)
		page := filepath.Join(appLoc, indexPage)
		logger.Debug("Serving index page: %s", page)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, page)
		} else {
			Error404(w)
		}
	}
}
