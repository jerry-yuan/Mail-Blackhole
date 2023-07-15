package web

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/pat"
	"github.com/jerry-yuan/mail-blackhole/config"
	"github.com/sirupsen/logrus"
)

var APIHost string
var WebPath string

//go:embed frontendv1
var frontendv1 embed.FS

////go:embed frontendv2/dist
//var frontendv2 embed.FS

type Web struct {
	config *config.WebUIConfig
	asset  func(string) ([]byte, error)
}

func CreateWeb(cfg *config.WebUIConfig, r http.Handler) *Web {
	web := &Web{
		config: cfg,
		asset: func(s string) ([]byte, error) {
			file, err := frontendv1.Open(s)
			if err != nil {
				return nil, err
			}
			return io.ReadAll(file)
		},
	}

	pat := r.(*pat.Router)

	WebPath = cfg.WebPath

	logrus.Infof("Serving under http://%s%s/", cfg.BindAddr, WebPath)

	pat.Path(WebPath + "/images/{file:.*}").Methods("GET").HandlerFunc(web.Static("frontendv1/images/{{file}}"))
	pat.Path(WebPath + "/css/{file:.*}").Methods("GET").HandlerFunc(web.Static("frontendv1/css/{{file}}"))
	pat.Path(WebPath + "/js/{file:.*}").Methods("GET").HandlerFunc(web.Static("frontendv1/js/{{file}}"))
	pat.Path(WebPath + "/fonts/{file:.*}").Methods("GET").HandlerFunc(web.Static("frontendv1/fonts/{{file}}"))
	pat.StrictSlash(true).Path(WebPath + "/").Methods("GET").HandlerFunc(web.Index())

	return web
}

func (web Web) Static(pattern string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fp := strings.TrimSuffix(pattern, "{{file}}") + req.URL.Query().Get(":file")
		if b, err := web.asset(fp); err == nil {
			ext := filepath.Ext(fp)

			w.Header().Set("Content-Type", mime.TypeByExtension(ext))
			w.WriteHeader(200)
			w.Write(b)
			return
		}
		logrus.Printf("[UI] File not found: %s", fp)
		w.WriteHeader(404)
	}
}

func (web Web) Index() func(http.ResponseWriter, *http.Request) {
	tmpl := template.New("index.html")
	tmpl.Delims("[:", ":]")

	asset, err := web.asset("frontendv1/templates/index.html")
	if err != nil {
		logrus.Fatalf("[UI] Error loading index.html: %s", err)
	}

	tmpl, err = tmpl.Parse(string(asset))
	if err != nil {
		logrus.Fatalf("[UI] Error parsing index.html: %s", err)
	}

	layout := template.New("layout.html")
	layout.Delims("[:", ":]")

	asset, err = web.asset("frontendv1/templates/layout.html")
	if err != nil {
		logrus.Fatalf("[UI] Error loading layout.html: %s", err)
	}

	layout, err = layout.Parse(string(asset))
	if err != nil {
		logrus.Fatalf("[UI] Error parsing layout.html: %s", err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		data := map[string]interface{}{
			"config":  web.config,
			"Page":    "Browse",
			"APIHost": APIHost,
		}

		b := new(bytes.Buffer)
		err := tmpl.Execute(b, data)

		if err != nil {
			logrus.Printf("[UI] Error executing template: %s", err)
			w.WriteHeader(500)
			return
		}

		data["Content"] = template.HTML(b.String())

		b = new(bytes.Buffer)
		err = layout.Execute(b, data)

		if err != nil {
			logrus.Printf("[UI] Error executing template: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write(b.Bytes())
	}
}
