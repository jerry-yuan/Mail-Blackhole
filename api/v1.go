package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/pat"
	"github.com/jerry-yuan/mail-blackhole/config"
	"github.com/jerry-yuan/mail-blackhole/data"
	"github.com/jerry-yuan/mail-blackhole/storage"
	"github.com/sirupsen/logrus"

	"github.com/ian-kent/goose"
)

// APIv1 implements version 1 of the MailHog API
//
// The specification has been frozen and will eventually be deprecated.
// Only bug fixes and non-breaking changes will be applied here.
//
// Any changes/additions should be added in APIv2.
type APIv1 struct {
	config      *config.APIConfig
	storage     storage.Storage
	messageChan chan *data.Message
}

// FIXME should probably move this into APIv1 struct
var stream *goose.EventStream

func createAPIv1(conf *config.APIConfig, msgStorage storage.Storage, r *pat.Router) *APIv1 {
	logrus.Infof("Creating API v1 with WebPath: " + conf.WebPath)
	apiv1 := &APIv1{
		config:      conf,
		storage:     msgStorage,
		messageChan: make(chan *data.Message),
	}

	stream = goose.NewEventStream()

	r.Path(conf.WebPath + "/api/v1/messages").Methods("GET").HandlerFunc(apiv1.messages)
	r.Path(conf.WebPath + "/api/v1/messages").Methods("DELETE").HandlerFunc(apiv1.delete_all)
	r.Path(conf.WebPath + "/api/v1/messages").Methods("OPTIONS").HandlerFunc(apiv1.defaultOptions)

	r.Path(conf.WebPath + "/api/v1/messages/{id}").Methods("GET").HandlerFunc(apiv1.message)
	r.Path(conf.WebPath + "/api/v1/messages/{id}").Methods("DELETE").HandlerFunc(apiv1.delete_one)
	r.Path(conf.WebPath + "/api/v1/messages/{id}").Methods("OPTIONS").HandlerFunc(apiv1.defaultOptions)

	r.Path(conf.WebPath + "/api/v1/messages/{id}/download").Methods("GET").HandlerFunc(apiv1.download)
	r.Path(conf.WebPath + "/api/v1/messages/{id}/download").Methods("OPTIONS").HandlerFunc(apiv1.defaultOptions)

	r.Path(conf.WebPath + "/api/v1/messages/{id}/mime/part/{part}/download").Methods("GET").HandlerFunc(apiv1.download_part)
	r.Path(conf.WebPath + "/api/v1/messages/{id}/mime/part/{part}/download").Methods("OPTIONS").HandlerFunc(apiv1.defaultOptions)

	r.Path(conf.WebPath + "/api/v1/messages/{id}/release").Methods("OPTIONS").HandlerFunc(apiv1.defaultOptions)

	r.Path(conf.WebPath + "/api/v1/events").Methods("GET").HandlerFunc(apiv1.eventstream)
	r.Path(conf.WebPath + "/api/v1/events").Methods("OPTIONS").HandlerFunc(apiv1.defaultOptions)

	go func() {
		keepaliveTicker := time.Tick(time.Minute)
		for {
			select {
			case msg := <-apiv1.messageChan:
				logrus.Tracef("Got message in APIv1 event stream")
				bytes, _ := json.MarshalIndent(msg, "", "  ")
				json := string(bytes)
				logrus.Tracef("Sending content: %s\n", json)
				apiv1.broadcast(json)
			case <-keepaliveTicker:
				apiv1.keepalive()
			}
		}
	}()

	return apiv1
}

func (apiv1 *APIv1) defaultOptions(w http.ResponseWriter, req *http.Request) {
	if len(apiv1.config.CORSOrigin) > 0 {
		w.Header().Add("Access-Control-Allow-Origin", apiv1.config.CORSOrigin)
		w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET,POST,DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	}
}

func (apiv1 *APIv1) broadcast(json string) {
	logrus.Println("[APIv1] BROADCAST /api/v1/events")
	b := []byte(json)
	stream.Notify("data", b)
}

// keepalive sends an empty keep alive message.
//
// This not only can keep connections alive, but also will detect broken
// connections. Without this it is possible for the server to become
// unresponsive due to too many open files.
func (apiv1 *APIv1) keepalive() {
	logrus.Println("[APIv1] KEEPALIVE /api/v1/events")
	stream.Notify("keepalive", []byte{})
}

func (apiv1 *APIv1) eventstream(w http.ResponseWriter, req *http.Request) {
	logrus.Println("[APIv1] GET /api/v1/events")

	//apiv1.defaultOptions(session)
	if len(apiv1.config.CORSOrigin) > 0 {
		w.Header().Add("Access-Control-Allow-Origin", apiv1.config.CORSOrigin)
		w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET,POST,DELETE")
	}

	stream.AddReceiver(w)
}

func (apiv1 *APIv1) messages(w http.ResponseWriter, req *http.Request) {
	logrus.Println("[APIv1] GET /api/v1/messages")

	apiv1.defaultOptions(w, req)

	// TODO start, limit
	switch apiv1.storage.(type) {
	case *storage.MongoDB:
		messages, _ := apiv1.storage.(*storage.MongoDB).List(0, 1000)
		bytes, _ := json.Marshal(messages)
		w.Header().Add("Content-Type", "text/json")
		w.Write(bytes)
	case *storage.InMemory:
		messages, _ := apiv1.storage.(*storage.InMemory).List(0, 1000)
		bytes, _ := json.Marshal(messages)
		w.Header().Add("Content-Type", "text/json")
		w.Write(bytes)
	default:
		w.WriteHeader(500)
	}
}

func (apiv1 *APIv1) message(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get(":id")
	logrus.Infof("[APIv1] GET /api/v1/messages/%s\n", id)

	apiv1.defaultOptions(w, req)

	message, err := apiv1.storage.Load(id)
	if err != nil {
		logrus.Printf("- Error: %s", err)
		w.WriteHeader(500)
		return
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		logrus.Printf("- Error: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.Write(bytes)
}

func (apiv1 *APIv1) download(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get(":id")
	logrus.Printf("[APIv1] GET /api/v1/messages/%s\n", id)

	apiv1.defaultOptions(w, req)

	w.Header().Set("Content-Type", "message/rfc822")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+id+".eml\"")

	switch apiv1.storage.(type) {
	case *storage.MongoDB:
		message, _ := apiv1.storage.(*storage.MongoDB).Load(id)
		for h, l := range message.Content.Headers {
			for _, v := range l {
				w.Write([]byte(h + ": " + v + "\r\n"))
			}
		}
		w.Write([]byte("\r\n" + message.Content.Body))
	case *storage.InMemory:
		message, _ := apiv1.storage.(*storage.InMemory).Load(id)
		for h, l := range message.Content.Headers {
			for _, v := range l {
				w.Write([]byte(h + ": " + v + "\r\n"))
			}
		}
		w.Write([]byte("\r\n" + message.Content.Body))
	default:
		w.WriteHeader(500)
	}
}

func (apiv1 *APIv1) download_part(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get(":id")
	part := req.URL.Query().Get(":part")
	logrus.Printf("[APIv1] GET /api/v1/messages/%s/mime/part/%s/download\n", id, part)

	// TODO extension from content-type?
	apiv1.defaultOptions(w, req)

	w.Header().Set("Content-Disposition", "attachment; filename=\""+id+"-part-"+part+"\"")

	message, _ := apiv1.storage.Load(id)
	contentTransferEncoding := ""
	pid, _ := strconv.Atoi(part)
	for h, l := range message.MIME.Parts[pid].Headers {
		for _, v := range l {
			switch strings.ToLower(h) {
			case "content-disposition":
				// Prevent duplicate "content-disposition"
				w.Header().Set(h, v)
			case "content-transfer-encoding":
				if contentTransferEncoding == "" {
					contentTransferEncoding = v
				}
				fallthrough
			default:
				w.Header().Add(h, v)
			}
		}
	}
	body := []byte(message.MIME.Parts[pid].Body)
	if strings.ToLower(contentTransferEncoding) == "base64" {
		var e error
		body, e = base64.StdEncoding.DecodeString(message.MIME.Parts[pid].Body)
		if e != nil {
			logrus.Printf("[APIv1] Decoding base64 encoded body failed: %s", e)
		}
	}
	w.Write(body)
}

func (apiv1 *APIv1) delete_all(w http.ResponseWriter, req *http.Request) {
	logrus.Println("[APIv1] POST /api/v1/messages")

	apiv1.defaultOptions(w, req)

	w.Header().Add("Content-Type", "text/json")

	err := apiv1.storage.DeleteAll()
	if err != nil {
		logrus.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (apiv1 *APIv1) delete_one(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get(":id")

	logrus.Printf("[APIv1] POST /api/v1/messages/%s/delete\n", id)

	apiv1.defaultOptions(w, req)

	w.Header().Add("Content-Type", "text/json")
	err := apiv1.storage.DeleteOne(id)
	if err != nil {
		logrus.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
