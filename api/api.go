package api

import (
	gohttp "net/http"

	"github.com/gorilla/pat"
	"github.com/jerry-yuan/mail-blackhole/config"
	"github.com/jerry-yuan/mail-blackhole/data"
	"github.com/jerry-yuan/mail-blackhole/storage"
)

func CreateAPI(conf *config.APIConfig, msgStorage storage.Storage, msgCh chan *data.Message, r gohttp.Handler) {
	apiv1 := createAPIv1(conf, msgStorage, r.(*pat.Router))
	apiv2 := createAPIv2(conf, msgStorage, r.(*pat.Router))

	go func() {
		for {
			select {
			case msg := <-msgCh:
				apiv1.messageChan <- msg
				apiv2.messageChan <- msg
			}
		}
	}()
}
