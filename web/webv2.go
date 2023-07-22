package web

import (
	"embed"
	"github.com/gorilla/pat"
	"github.com/jerry-yuan/mail-blackhole/config"
	"github.com/sirupsen/logrus"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

//go:embed frontendv2/build
var frontendv2 embed.FS

func CreateWebV2(cfg *config.WebUIConfig, r http.Handler) {
	/// web := &Web{
	/// 	config: cfg,
	/// 	asset: func(s string) ([]byte, error) {
	/// 		file, err := frontendv2.Open(s)
	/// 		if err != nil {
	/// 			return nil, err
	/// 		}
	/// 		return io.ReadAll(file)
	/// 	},
	/// }

	pat := r.(*pat.Router)

	WebPath = cfg.WebPath

	logrus.Infof("Serving under http://%s%s/", cfg.BindAddr, WebPath)

	pat.PathPrefix("/").Methods("GET").HandlerFunc(serveAssets)
}

func serveAssets(writer http.ResponseWriter, request *http.Request) {
	path := "frontendv2/build" + request.URL.Path

	if file, err := frontendv2.Open(path); err == nil {
		bytes, err := io.ReadAll(file)
		if err == nil {
			ext := filepath.Ext(path)
			writer.Header().Set("Content-Type", mime.TypeByExtension(ext))
			writer.WriteHeader(200)
			writer.Write(bytes)
			return
		}
	}
	// file not found
	path = "frontendv2/build/index.html"
	if file, err := frontendv2.Open(path); err == nil {
		bytes, err := io.ReadAll(file)
		if err == nil {
			writer.WriteHeader(200)
			writer.Write(bytes)
			return
		}
	}
	// index is missing
	writer.WriteHeader(404)
	writer.Write([]byte("Index is missing."))
}
