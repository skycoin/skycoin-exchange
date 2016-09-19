package gui

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"

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

// LaunchWebInterface begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, rt *httprouter.Router) error {
	logger.Info("Starting web interface on http://%s", host)
	logger.Warning("HTTPS not in use!")
	// logger.Info("Web resources directory: %s", staticDir)

	appLoc, err := determineResourcePath(staticDir)
	if err != nil {
		return err
	}
	logger.Debug("static dir:%s", appLoc)

	rt.NotFound = http.FileServer(http.Dir(appLoc))
	go func() {
		log.Panic(http.ListenAndServe(host, rt))
	}()
	return nil
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
		//check build directory
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
