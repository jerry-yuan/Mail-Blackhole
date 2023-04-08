package main

import (
	"fmt"
	"os"

	gohttp "net/http"

	"github.com/gorilla/pat"
	"github.com/jerry-yuan/mail-blackhole/api"
	"github.com/jerry-yuan/mail-blackhole/config"
	"github.com/jerry-yuan/mail-blackhole/data"
	"github.com/jerry-yuan/mail-blackhole/http"
	"github.com/jerry-yuan/mail-blackhole/smtp"
	"github.com/jerry-yuan/mail-blackhole/storage"
	"github.com/jerry-yuan/mail-blackhole/web"
	"github.com/sirupsen/logrus"
)

var conf *config.Config

var messageChan chan *data.Message = make(chan *data.Message, 10)
var exitCh chan struct{} = make(chan struct{})
var version string

func main() {

	// load configurations
	conf = config.Configure()

	if conf.PrintVersion {
		fmt.Println("Mail Blackhole version: " + version)
		os.Exit(0)
	}

	// initialize log

	level, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logrus.Warnf("Unrecognized log level %s, fallback to use info", conf.Log.Level)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{})

	// initialize storage
	msgStorage, err := storage.Create(conf.Storage)
	if err != nil {
		logrus.Fatal(err)
	}

	// initialize smtp server
	go smtp.Listen(&conf.SMTP, msgStorage, messageChan, exitCh)

	// initialize http server
	if conf.WebUI.AuthFile != "" {
		http.AuthFile(conf.WebUI.AuthFile)
	}

	if conf.WebUI.BindAddr == conf.API.BindAddr {
		cb := func(r gohttp.Handler) {
			web.CreateWeb(&conf.WebUI, r.(*pat.Router))
			api.CreateAPI(&conf.API, msgStorage, messageChan, r.(*pat.Router))
		}
		go http.Listen(conf.API.BindAddr, exitCh, cb)
	} else {
		cb1 := func(r gohttp.Handler) {
			api.CreateAPI(&conf.API, msgStorage, messageChan, r.(*pat.Router))
		}
		cb2 := func(r gohttp.Handler) {
			web.CreateWeb(&conf.WebUI, r.(*pat.Router))
		}
		go http.Listen(conf.API.BindAddr, exitCh, cb1)
		go http.Listen(conf.WebUI.BindAddr, exitCh, cb2)
	}

	<-exitCh
	logrus.Infof("Received exit signal")
}

/*

Add some random content to the end of this file, hopefully tricking GitHub
into recognising this as a Go repo instead of Makefile.

A gopher, ASCII art style - borrowed from
https://gist.github.com/belbomemo/b5e7dad10fa567a5fe8a

          ,_---~~~~~----._
   _,,_,*^____      _____``*g*\"*,
  / __/ /'     ^.  /      \ ^@q   f
 [  @f | @))    |  | @))   l  0 _/
  \`/   \~____ / __ \_____/    \
   |           _l__l_           I
   }          [______]           I
   ]            | | |            |
   ]             ~ ~             |
   |                            |
    |                           |

*/
